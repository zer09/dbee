package boltengine

import (
	"crypto/rand"
	"dbee/endian"
	"dbee/errors"
	"dbee/internal/boltengine/schema"
	"dbee/store"
	"path/filepath"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/oklog/ulid"

	bolt "go.etcd.io/bbolt"
)

var (
	// props is the bucket that will store the propertiest/column.
	props = []byte("props")
	// prop_names is a bucket inside prop that will store property name as value.
	prop_names = []byte("prop_names")
	// prop_indexes is a bucket inside prop that will store property as key.
	// And the value will be the index number of the property name.
	prop_indexes = []byte("prop_indexes")
	// sets is the bucket that hold the set information.
	sets = []byte("sets")
	// rootBucket of the store
	rootBucket = endian.I64toB(0)
)

// boltOpt the default option for bolt.
var boltOpt = &bolt.Options{Timeout: 1 * time.Second}

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
	props *instanceProp
	// sets available sets for ths instace.
	// set key name will be this formate setname:partition name,
	// every set will be have a default name partition named `default`.
	sets map[string]*Set
}

// instanceProp will hold the reverse index of property and property index.
type instanceProp struct {
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
			_, err := b.CreateBucketIfNotExists(prop_names)
			if err != nil {
				return err
			}

			_, err = b.CreateBucketIfNotExists(prop_indexes)
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
		props: &instanceProp{
			name:  make(map[uint64]string),
			index: make(map[string]uint64),
		},
		sets: make(map[string]*Set),
	}, nil
}

// Close the instance.
func (i *Instance) Close() error {
	var err error
	for _, v := range i.sets {
		for _, p := range v.partitions {
			if err = p.store.Close(); err != nil {
				return err
			}

			if err = p.index.Close(); err != nil {
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
			b := tx.Bucket(props).Bucket(prop_names)
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
	if v, ok := i.props.index[name]; ok {
		return v, nil
	} else {
		_ = i.meta.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(props).Bucket(prop_indexes)
			index := b.Get([]byte(name))
			if index != nil {
				v = endian.BtoI64(index)
			}

			return nil
		})

		if v < 1 {
			err := i.meta.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket(props).Bucket(prop_names)
				i, err := b.NextSequence()
				if err != nil {
					return err
				}

				id := endian.I64toB(i)
				err = b.Put(id, []byte(name))
				if err != nil {
					return err
				}

				v = i
				return tx.Bucket(props).Bucket(prop_indexes).Put([]byte(name), id)
			})

			if err != nil {
				return v, err
			}
		}

		return v, nil
	}
}

// Set will get an instance of set.
// name is the name of the set.
func (i *Instance) Set(name string) (store.Set, error) {
	if s, ok := i.sets[name]; ok {
		return s, nil
	}

	s := &Set{
		Instance:   i,
		name:       name,
		partitions: make(map[string]*partition),
	}

	err := i.meta.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(sets)
		// pBuf the key of the set or the set name/name.
		pBuf := []byte(s.name)
		// sBuf is the schema.Set proto containing the meta of the set.
		sBuf := b.Get(pBuf)

		if sBuf == nil {
			setSchema := &schema.Set{
				Partition: make([]*schema.Partition, 1),
			}

			defPart := &schema.Partition{
				Name:  defaultPartition,
				Index: ulid.MustNew(ulid.Now(), rand.Reader).String(),
				Store: ulid.MustNew(ulid.Now(), rand.Reader).String(),
			}

			setSchema.Partition[0] = defPart
			setBuf, err := proto.Marshal(setSchema)
			if err != nil {
				return err
			}

			s.partitions[defPart.Name] = &partition{
				Name:  defPart.Name,
				Index: defPart.Index,
				Store: defPart.Store,
			}

			return b.Put(pBuf, setBuf)
		} else {
			setSchema := &schema.Set{}
			err := proto.Unmarshal(sBuf, setSchema)
			if err != nil {
				return err
			}

			for _, defPart := range setSchema.Partition {
				s.partitions[defPart.Name] = &partition{
					Name:  defPart.Name,
					Index: defPart.Index,
					Store: defPart.Store,
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	i.sets[s.name] = s

	for _, v := range s.partitions {
		v.opened = true
		if v.index, err = open(filepath.Join(s.dir, v.Index)); err != nil {
			return nil, err
		}

		if v.store, err = open(filepath.Join(s.dir, v.Store)); err != nil {
			return nil, err
		}
	}

	return s, s.preparRootBucket()
}
