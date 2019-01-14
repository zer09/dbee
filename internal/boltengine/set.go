package boltengine

import (
	"crypto/rand"
	"dbee/internal/boltengine/schema"
	"dbee/store"
	"path/filepath"
	"reflect"

	"github.com/gogo/protobuf/proto"
	"github.com/oklog/ulid"
	bolt "go.etcd.io/bbolt"
)

type Set struct {
	*Instance
	// name of the set.
	name string
	// partitions of the set.
	partitions map[string]*partition
	// List of all indexed
	idxs *indexMap
}

func (s *Set) Name() string {
	return s.name
}

func (s *Set) Partitions() []string {
	keys := make([]string, len(s.partitions))
	i := 0
	for k := range s.partitions {
		keys[i] = k
	}

	return keys
}

func (s *Set) Partition(name string) (store.Partition, error) {
	return s.getPartition(name)
}

func (s *Set) Index(propName string) error {
	return s.newIDx(propName, s)
}

func (s *Set) ListIndexes() []string {
	keys := reflect.ValueOf(s.idxs.index).MapKeys()
	ks := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		ks[i] = keys[i].String()
	}

	return ks
}

func (s *Set) Get(id ...string) (store.SetTx, error) {
	return s.partitions[defaultPartition].Get(id...)
}

func (s *Set) getPartition(partitionName string) (*partition, error) {
	if v, ok := s.partitions[partitionName]; ok {
		return v, nil
	}

	part := &schema.Partition{
		Name:  partitionName,
		Index: ulid.MustNew(ulid.Now(), rand.Reader).String(),
		Store: ulid.MustNew(ulid.Now(), rand.Reader).String(),
	}

	err := s.meta.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(sets)
		pBuf := []byte(s.name)
		sBuf := b.Get(pBuf)

		setSchema := &schema.Set{}
		err := proto.Unmarshal(sBuf, setSchema)
		if err != nil {
			return err
		}

		for _, n := range setSchema.Partition {
			if n.Name == part.Name {
				return nil
			}
		}

		setSchema.Partition = append(setSchema.Partition, part)
		setBuf, err := proto.Marshal(setSchema)
		if err != nil {
			return err
		}

		return b.Put(pBuf, setBuf)
	})

	if err != nil {
		return nil, err
	}

	// recheck again if the partition exists
	// because during inside the update it could be this partition
	// is also created.
	if v, ok := s.partitions[partitionName]; ok {
		return v, nil
	}

	s.partitions[partitionName] = &partition{
		name:      part.Name,
		indexName: part.Index,
		storeName: part.Store,
		set:       s,
	}

	for _, p := range s.partitions {
		if p.opened {
			continue
		}

		p.opened = true
		if p.index, err = open(filepath.Join(s.dir, p.indexName)); err != nil {
			return nil, err
		}

		if p.store, err = open(filepath.Join(s.dir, p.storeName)); err != nil {
			return nil, err
		}
	}

	return s.partitions[partitionName], s.preparRootBucket()
}

func (s *Set) preparRootBucket() error {
	for _, v := range s.partitions {
		err := v.index.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(rootBucket)
			return err
		})

		if err != nil {
			return err
		}

		err = v.store.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(rootBucket)
			return err
		})

		if err != nil {
			return err
		}
	}

	return nil
}
