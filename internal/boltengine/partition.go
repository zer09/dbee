package boltengine

import (
	"crypto/rand"
	"dbee/endian"
	"dbee/errors"
	"dbee/internal/boltengine/schema"
	"dbee/store"
	"strings"

	"github.com/oklog/ulid"
	bolt "go.etcd.io/bbolt"
)

type Partition struct {
	// name of the Partition.
	name string
	// indexName file location of the Partition.
	// indexName string
	// storeName file location of the Partition.
	// This will be the storage of the data.
	storeName string
	// index     *bolt.DB
	store *bolt.DB
	set   *Set
}

func (p *Partition) Name() string {
	return p.name
}

func (p *Partition) Get(id ...string) (store.SetTx, error) {
	var err error

	sx := &SetTx{
		set:       p.set,
		idBuf:     make([]byte, 16),
		partition: p,
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

func (p *Partition) getByBytes(ulidByte []byte) (store.SetTx, error) {
	var err error

	sx := &SetTx{
		set:       p.set,
		idBuf:     make([]byte, 16),
		partition: p,
		payload: &schema.Payload{
			Meta:   &schema.Meta{Deleted: false},
			Values: make(map[uint64][]byte),
		},
	}

	err = sx.id.UnmarshalBinary(ulidByte)
	if err != nil {
		return nil, err
	}

	err = sx.loadPayload()
	if err != nil {
		return nil, err
	}

	return sx, err
}

func (p *Partition) One(n string, v string) (store.SetTx, error) {
	sx, err := p.getOneByIndex(n, v)
	if err != nil {
		return nil, err
	}

	if sx != nil {
		return sx, nil
	}

	return nil, nil
}

func (p *Partition) getOneByIndex(n string, v string) (store.SetTx, error) {
	propIdx := p.set.idxs.getIndex(n)
	if propIdx < 1 {
		return nil, nil
	}

	propIdxBuf := endian.I64toB(propIdx)
	v = strings.ToLower(strings.TrimSpace(v))
	target := []byte(v)

	var stx store.SetTx
	var kFind []byte

	_ = p.store.View(func(tx *bolt.Tx) error {
		kb := tx.Bucket(indexBucket).Bucket(indexStringBucket).Bucket(propIdxBuf)
		if kb == nil {
			return nil
		}

		vb := kb.Bucket(target)
		if vb == nil {
			return nil
		}

		k, _ := vb.Cursor().First()
		if k == nil {
			return nil
		}

		kFind = k
		return nil
	})

	if len(kFind) != 16 {
		return nil, nil
	}

	stx, err := p.getByBytes(kFind)
	if err != nil {
		return nil, err
	}

	return stx, nil
}
