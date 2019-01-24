package boltengine

import (
	"bytes"
	"crypto/rand"
	"dbee/endian"
	"dbee/errors"
	"dbee/internal/boltengine/schema"
	"dbee/store"

	proto "github.com/golang/protobuf/proto"
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
		indexable: make(map[uint64]interface{}),
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

func (p *Partition) One(n string, v interface{}) (store.SetTx, error) {
	sx, err := p.getOneByIndex(n, v)
	if err != nil {
		return nil, err
	}

	if sx != nil {
		return sx, nil
	}

	return nil, nil
}

func (p *Partition) getOneByIndex(n string, v interface{}) (store.SetTx, error) {
	propIdx := p.set.idxs.getIndex(n)
	if propIdx < 1 {
		return nil, nil
	}

	var stx store.SetTx

	idBuf := endian.I64toB(propIdx)
	err := p.store.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(indexBucket)
		iValBuc := b.Bucket(idBuf)
		if iValBuc != nil {
			c := iValBuc.Cursor()
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				stxtmp, err := p.getByBytes(k)
				if err != nil {
					return err
				}

				switch value := v.(type) {
				case float32:
					if stxtmp.Rfloat(n) == value {
						stx = stxtmp
						return nil
					}
				case float64:
					if stxtmp.Rdouble(n) == value {
						stx = stxtmp
						return nil
					}
				case int64:
					if stxtmp.Rint(n) == value || stxtmp.Rsint(n) == value {
						stx = stxtmp
						return nil
					}
				case uint64:
					if stxtmp.Ruint(n) == value {
						stx = stxtmp
						return nil
					}
				case bool:
					if stxtmp.Rbool(n) == value {
						stx = stxtmp
						return nil
					}
				case []byte:
					if bytes.Equal(stxtmp.Rbytes(n), value) {
						stx = stxtmp
						return nil
					}
				case string:
					if stxtmp.Rstring(n) == value {
						stx = stxtmp
						return nil
					}
				default:
					if stxtmp.Rsint(n) == value {
						stx = stxtmp
						return nil
					}
				}
			}
		}

		isBuf := b.Get(idBuf)
		if isBuf != nil {
			is := &schema.IndexSlice{}
			err := proto.Unmarshal(isBuf, is)
			if err != nil {
				return err
			}

			for _, iv := range is.IDIndexes {
				stxtmp, err := p.getByBytes(iv)
				if err != nil {
					return err
				}

				switch value := v.(type) {
				case float32:
					if stxtmp.Rfloat(n) == value {
						stx = stxtmp
						return nil
					}
				case float64:
					if stxtmp.Rdouble(n) == value {
						stx = stxtmp
						return nil
					}
				case int64:
					if stxtmp.Rint(n) == value || stxtmp.Rsint(n) == value {
						stx = stxtmp
						return nil
					}
				case uint64:
					if stxtmp.Ruint(n) == value {
						stx = stxtmp
						return nil
					}
				case bool:
					if stxtmp.Rbool(n) == value {
						stx = stxtmp
						return nil
					}
				case []byte:
					if bytes.Equal(stxtmp.Rbytes(n), value) {
						stx = stxtmp
						return nil
					}
				case string:
					if stxtmp.Rstring(n) == value {
						stx = stxtmp
						return nil
					}
				default:
					if stxtmp.Rsint(n) == value {
						stx = stxtmp
						return nil
					}
				}
			}
		}

		return nil
	})

	return stx, err
}
