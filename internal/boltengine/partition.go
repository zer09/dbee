package boltengine

import (
	"crypto/rand"
	"dbee/errors"
	"dbee/internal/boltengine/schema"
	"dbee/store"

	"github.com/oklog/ulid"
	bolt "go.etcd.io/bbolt"
)

type partition struct {
	// name of the partition.
	name string
	// indexName file location of the partition.
	indexName string
	// storeName file location of the partition.
	// This will be the storage of the data.
	storeName string
	opened    bool
	index     *bolt.DB
	store     *bolt.DB
	set       *Set
}

func (p *partition) Name() string {
	return p.name
}

func (p *partition) Get(id ...string) (store.SetTx, error) {
	var err error

	sx := &SetTx{
		set:   p.set,
		idBuf: make([]byte, 16),
		partition:  p,
		payload: &schema.Payload{
			Meta:   &schema.Meta{Deleted: false},
			Values: make(map[uint64][]byte),
		},
	}

	if len(id) > 0 {
		if len(id[0]) != 26 {
			return nil, errors.ErrNotValidULID
		}

		if sx.id, err = ulid.ParseStrict(id[0]); err != nil {
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

	return sx, nil
}
