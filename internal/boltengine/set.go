package boltengine

import (
	"crypto/rand"
	"dbee/errors"
	"dbee/internal/boltengine/schema"
	"dbee/store"
	"path/filepath"
	"strings"

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

func (s *Set) Get(param ...string) (store.SetTx, error) {
	var err error
	sx := &SetTx{
		set:       s,
		idBuf:     make([]byte, 16),
		partition: defaultPartition,
		payload: &schema.Payload{
			Meta:   &schema.Meta{Deleted: false},
			Values: make(map[uint64][]byte),
		},
	}

	if len(param) > 0 {
		l := len(param[0])
		if l > 0 && l < 26 {
			return nil, errors.ErrNotValidULID
		}

		if sx.id, err = ulid.ParseStrict(param[0]); err != nil {
			return nil, err
		}

		if err = sx.loadPayload(); err != nil {
			return nil, err
		}
	} else {
		if sx.id, err = ulid.New(ulid.Now(), rand.Reader); err != nil {
			return nil, err
		}

		if err = sx.id.MarshalBinaryTo(sx.idBuf); err != nil {
			return nil, err
		}
	}

	if len(param) > 1 && len(strings.TrimSpace(param[1])) > 0 {
		sx.partition = strings.TrimSpace(strings.ToLower(param[1]))
		err = s.newPartitions(sx.partition)
		if err != nil {
			return nil, err
		}
	}

	return sx, nil
}

func (s *Set) newPartitions(partitionName string) error {
	if _, ok := s.partitions[partitionName]; ok {
		return nil
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
		return err
	}

	s.partitions[part.Name] = &partition{
		Name:  part.Name,
		Index: part.Index,
		Store: part.Store,
	}

	for _, p := range s.partitions {
		if p.opened {
			continue
		}

		p.opened = true
		if p.index, err = open(filepath.Join(s.dir, p.Index)); err != nil {
			return err
		}

		if p.store, err = open(filepath.Join(s.dir, p.Store)); err != nil {
			return err
		}
	}

	return s.preparRootBucket()
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
