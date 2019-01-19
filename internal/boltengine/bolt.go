package boltengine

import (
	"crypto/rand"
	"dbee/endian"
	"dbee/errors"
	"dbee/internal/boltengine/schema"
	"dbee/store"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/oklog/ulid"

	bolt "go.etcd.io/bbolt"
)

var (
	// props is the bucket that will store the propertiest/column.
	props = []byte("props")
	// propNames is a bucket inside prop that will store property name as value.
	propNames = []byte("propNames")
	// propIndexes is a bucket inside prop that will store property as key.
	// And the value will be the index number of the property name.
	propIndexes = []byte("propIndexes")
	// sets is the bucket that hold the set information.
	sets = []byte("sets")
	// rootBucket of the store
	rootBucket = endian.I64toB(0)
	// indexBucket will store all the index of the set.
	indexBucket = endian.I64toB(1)
	// emptySlice will be used for values on empty for keys.
	emptySlice = []byte{}
	// emptyStruct an empty struct, mostly used for keys only map.
	emptyStruct struct{}
)

// boltOpt the default option for bolt.
var boltOpt = &bolt.Options{Timeout: 1 * time.Second}
var pageSize = os.Getpagesize()
var bucketMinSize = (pageSize * 3) - 16

// metaMagicName is the magic file name of the meta information.
const metaMagicName string = "01CRZFW6MBXA18393078QPCDQ7"
const defaultPartition string = "default"

func open(path string) (*bolt.DB, error) {
	return bolt.Open(path, 0600, boltOpt)
}

// Instance of a bolt engine.
type Instance struct {
	dir   string
	meta  *bolt.DB
	props *instancePropMap
	mux   sync.Mutex
	// sets available sets for ths instace.
	// set key name will be this formate setname:Partition name,
	// every set will be have a default name Partition named `default`.
	sets map[string]*Set
}

// instancePropMap will hold the reverse index of property and property index.
type instancePropMap struct {
	// name will hold name as the value, and will be accessed using the index.
	name map[uint64]string
	// index will hold the index as value, and will be accessed using the prop.
	index map[string]uint64
}

// New instance of Instance.
func New(dir string) (*Instance, error) {
	meta, err := open(filepath.Join(dir, metaMagicName))
	if err != nil {
		return nil, err
	}

	err = meta.Update(func(tx *bolt.Tx) error {
		if b, err := tx.CreateBucketIfNotExists(props); err != nil {
			return err
		} else {
			_, err := b.CreateBucketIfNotExists(propNames)
			if err != nil {
				return err
			}

			_, err = b.CreateBucketIfNotExists(propIndexes)
			if err != nil {
				return err
			}
		}

		if _, err := tx.CreateBucketIfNotExists(sets); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &Instance{
		dir:  dir,
		meta: meta,
		props: &instancePropMap{
			name:  make(map[uint64]string),
			index: make(map[string]uint64),
		},
		sets: make(map[string]*Set),
	}, nil
}

// Close the instance.
func (i *Instance) Close() error {
	i.mux.Lock()
	defer i.mux.Unlock()
	var err error
	for _, v := range i.sets {
		for _, p := range v.partitions {
			if err = p.store.Close(); err != nil {
				return err
			}
		}
	}

	return i.meta.Close()
}

// Dir of the instance.
func (i *Instance) Dir() string {
	return i.dir
}

// GetPropName get the property name using an index number.
func (i *Instance) GetPropName(index uint64) (string, error) {
	if v, ok := i.props.name[index]; ok {
		return v, nil
	} else {
		err := i.meta.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(props).Bucket(propNames)
			name := b.Get(endian.I64toB(index))
			if name != nil {
				v = string(name)
				return nil
			}

			return errors.ErrPropNotFound
		})

		return v, err
	}
}

// GetPropIndex will get the property index using the property name.
// If the name is not indexed then it will store the index and property name.
func (i *Instance) GetPropIndex(name string) (uint64, error) {
	name = sanitizeProp(name)

	var v uint64
	if v, ok := i.props.index[name]; ok {
		return v, nil
	}

	err := i.meta.Update(func(tx *bolt.Tx) error {
		propIdxBucket := tx.Bucket(props).Bucket(propIndexes)
		index := propIdxBucket.Get([]byte(name))
		if index != nil {
			v = endian.BtoI64(index)
			return nil
		}

		propNamesBucket := tx.Bucket(props).Bucket(propNames)
		i, err := propNamesBucket.NextSequence()
		if err != nil {
			return err
		}

		id := endian.I64toB(i)
		nameByte := []byte(name)
		err = propNamesBucket.Put(id, nameByte)
		if err != nil {
			return err
		}

		err = propIdxBucket.Put(nameByte, id)
		if err != nil {
			return err
		}

		v = i
		return nil
	})

	return v, err
}

// Set will get an instance of set.
// name is the name of the set.
func (i *Instance) Set(name string) (store.Set, error) {
	i.mux.Lock()
	defer i.mux.Unlock()

	if s, ok := i.sets[name]; ok {
		return s, nil
	}

	s := &Set{
		instance:   i,
		name:       name,
		partitions: make(map[string]*Partition),
	}

	err := i.meta.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(sets)
		// kBuf the key of the set or the set name/name.
		kBuf := []byte(s.name)
		// sBuf is the schema.Set proto containing the meta of the set.
		sBuf := b.Get(kBuf)

		if sBuf == nil {
			setSchema := &schema.Set{
				Partitions: make(map[string]*schema.Partition, 1),
			}

			store, err := ulid.New(ulid.Now(), rand.Reader)
			if err != nil {
				return err
			}

			setSchema.Partitions[defaultPartition] = &schema.Partition{
				Name:  defaultPartition,
				Store: store.String(),
			}

			setBuf, err := proto.Marshal(setSchema)
			if err != nil {
				return err
			}

			err = b.Put(kBuf, setBuf)
			if err != nil {
				return err
			}

			s.partitions[defaultPartition] = &Partition{
				name:      defaultPartition,
				storeName: store.String(),
				set:       s,
			}

			return nil
		}

		setSchema := &schema.Set{}
		err := proto.Unmarshal(sBuf, setSchema)
		if err != nil {
			return err
		}

		for k, v := range setSchema.Partitions {
			s.partitions[k] = &Partition{
				name:      v.Name,
				storeName: v.Store,
				set:       s,
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	i.sets[s.name] = s

	for _, v := range s.partitions {
		err = openPartition(v, s.instance.dir)
		if err != nil {
			return nil, err
		}
	}

	err = s.initIdxs()
	if err != nil {
		return nil, err
	}

	return s, s.preparRootBuckets()
}
