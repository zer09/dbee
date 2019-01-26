package boltengine

import (
	"crypto/rand"
	"dbee/internal/boltengine/schema"
	"dbee/store"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/oklog/ulid"
	bolt "go.etcd.io/bbolt"
)

const indexBucketPrefix = "idx_"

type Set struct {
	instance *Instance
	// name of the set.
	name string
	// List of all indexed
	idxs *indexMap

	mux sync.Mutex
	// partitions of the set.
	partitions map[string]*Partition
}

func (s *Set) Name() string {
	return s.name
}

func (s *Set) ListIndexes() []string {
	return s.idxs.indexes()
}

func (s *Set) Index(propName string) error {
	setIdx := indexBucketPrefix + s.name
	idxStr := sanitizeProp(propName)
	idxBuf := []byte(idxStr)

	err := s.instance.meta.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(setIdx))
		if err != nil {
			return err
		}

		exists := b.Get(idxBuf)
		if exists != nil {
			return nil
		}

		return b.Put(idxBuf, emptySlice)
	})

	if err != nil {
		return err
	}

	propID, err := s.instance.GetPropIndex(idxStr)
	if err != nil {
		return err
	}

	s.idxs.add(propID, idxStr)
	return nil
}

// initIdxs should only be called once in set creation.
func (s *Set) initIdxs() error {
	setIdx := indexBucketPrefix + s.name

	return s.instance.meta.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(setIdx))

		if b == nil {
			s.idxs = &indexMap{
				name:  make(map[uint64]string),
				index: make(map[string]uint64),
			}

			return nil
		}

		s.idxs = &indexMap{
			name:  make(map[uint64]string, b.Stats().KeyN),
			index: make(map[string]uint64, b.Stats().KeyN),
		}

		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			idxStr := string(k)
			propID, err := s.instance.GetPropIndex(idxStr)
			if err != nil {
				return err
			}

			s.idxs.add(propID, idxStr)
		}

		return nil
	})
}

func (s *Set) Get(id ...string) (store.SetTx, error) {
	return s.partitions[defaultPartition].Get(id...)
}

func (s *Set) Partitions() []string {
	s.mux.Lock()
	defer s.mux.Unlock()

	keys := make([]string, len(s.partitions))
	i := 0
	for k := range s.partitions {
		keys[i] = k
		i++
	}

	return keys
}

func (s *Set) Partition(name string) (store.Partition, error) {
	return s.getPartition(name)
}

func (s *Set) getPartition(partitionName string) (*Partition, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if v, ok := s.partitions[partitionName]; ok {
		return v, nil
	}

	var p *schema.Partition
	store, err := ulid.New(ulid.Now(), rand.Reader)
	if err != nil {
		return nil, err
	}

	err = s.instance.meta.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(sets)
		kBuf := []byte(s.name)
		sBuf := b.Get(kBuf)

		setSchema := &schema.Set{}
		err := proto.Unmarshal(sBuf, setSchema)
		if err != nil {
			return err
		}

		if v, ok := setSchema.Partitions[partitionName]; ok {
			p = v
			return nil
		}

		p = &schema.Partition{
			Name:  partitionName,
			Store: store.String(),
		}

		setSchema.Partitions[partitionName] = p
		sBuf, err = proto.Marshal(setSchema)
		if err != nil {
			return err
		}

		return b.Put(kBuf, sBuf)
	})

	if err != nil {
		return nil, err
	}

	s.partitions[partitionName] = &Partition{
		name:      p.Name,
		storeName: p.Store,
		set:       s,
	}

	for _, v := range s.partitions {
		err = openPartition(v, s.instance.dir)
		if err != nil {
			return nil, err
		}
	}

	return s.partitions[partitionName], s.prepareRootBuckets()
}

func (s *Set) prepareRootBuckets() error {
	var err error
	for _, v := range s.partitions {
		err = v.store.Update(func(tx *bolt.Tx) error {
			var err error
			_, err = tx.CreateBucketIfNotExists(rootBucket)
			if err != nil {
				return err
			}

			b, err := tx.CreateBucketIfNotExists(indexBucket)
			if err != nil {
				return err
			}

			_, err = b.CreateBucketIfNotExists(indexStringBucket)
			if err != nil {
				return nil
			}

			_, err = b.CreateBucketIfNotExists(indexUintBucket)
			if err != nil {
				return err
			}

			_, err = b.CreateBucketIfNotExists(indexNintBucket)

			return err
		})

		if err != nil {
			return err
		}
	}

	return nil
}
