package boltengine

import (
	"dbee/endian"
	"dbee/internal/boltengine/schema"
	"dbee/store"
	"encoding/binary"
	"time"

	proto "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/oklog/ulid"
	bolt "go.etcd.io/bbolt"
)

type SetTx struct {
	id         ulid.ULID
	idBuf      []byte
	set        *Set
	partition  *Partition
	payload    *schema.Payload
	payloadBuf []byte
	err        error
}

func (sx *SetTx) ID() string {
	return sx.id.String()
}

func (sx *SetTx) CreatedOn() (t time.Time) {
	if sx.payload.Meta.CreatedOn == nil {
		return time.Time{}
	}

	t, sx.err = ptypes.Timestamp(sx.payload.Meta.CreatedOn)
	return t
}

func (sx *SetTx) LastUpdate() (t time.Time) {
	if sx.payload.Meta.CreatedOn == nil {
		return time.Time{}
	}

	t, sx.err = ptypes.Timestamp(sx.payload.Meta.LastUpdate)
	return t
}

func (sx *SetTx) Partition() store.Partition {
	return sx.partition
}

func (sx *SetTx) IsSoftDeleted() bool {
	return sx.payload.Meta.Deleted
}

// Delete the data softly.
// you still need to execut commit befor the item will be deleted.
func (sx *SetTx) Delete() {
	sx.payload.Meta.Deleted = true
}

// HardDelete the data.
// This will auto commit the delete.
func (sx *SetTx) HardDelete() error {
	if sx.err != nil {
		return sx.err
	}

	// Check if the item is alread stored.
	if sx.payload.Meta.CreatedOn == nil {
		return nil
	}

	sx.payload.Meta.LastUpdate, sx.err = ptypes.TimestampProto(time.Now().UTC())
	if sx.err != nil {
		return sx.err
	}

	if sx.payloadBuf, sx.err = proto.Marshal(sx.payload); sx.err != nil {
		return sx.err
	}

	sx.err = sx.partition.store.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(rootBucket).Delete(sx.idBuf)
	})

	return sx.err
}

func (sx *SetTx) Commit() error {
	if sx.err != nil {
		return sx.err
	}

	if sx.payload.Meta.CreatedOn == nil {
		sx.payload.Meta.CreatedOn, sx.err = ptypes.TimestampProto(time.Now().UTC())
		if sx.err != nil {
			return sx.err
		}
	}

	sx.payload.Meta.LastUpdate, sx.err = ptypes.TimestampProto(time.Now().UTC())
	if sx.err != nil {
		return sx.err
	}

	if sx.payloadBuf, sx.err = proto.Marshal(sx.payload); sx.err != nil {
		return sx.err
	}

	sx.err = sx.partition.store.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket(rootBucket).Put(sx.idBuf, sx.payloadBuf)
		if err != nil {
			return err
		}

		return sx.commitIndex(tx)
	})

	return sx.err
}

func (sx *SetTx) commitIndex(tx *bolt.Tx) error {
	var b *bolt.Bucket
	var err error
	for k, _ := range sx.payload.Values {
		if !sx.set.idxs.indexable(k) {
			continue
		}

		if b == nil {
			b = tx.Bucket(indexBucket)
		}

		idBuf := endian.I64toB(k)
		iValBuc := b.Bucket(idBuf)
		if iValBuc != nil {
			err = sx.storeIndexInBucket(iValBuc)
			if err != nil {
				return err
			}
		}

		err = sx.storeIndexInSlice(idBuf, b)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sx *SetTx) storeIndexInBucket(iValBuc *bolt.Bucket) error {
	return iValBuc.Put(sx.idBuf, emptySlice)
}

func (sx *SetTx) storeIndexInSlice(idBuf []byte, b *bolt.Bucket) error {
	var is *schema.IndexSlice
	var err error

	IDxsBuf := b.Get(idBuf)
	err = proto.Unmarshal(IDxsBuf, is)
	if err != nil {
		return err
	}

	if binary.Size(IDxsBuf) >= bucketMinSize {
		return sx.convertSliceToBucket(idBuf, is, b)
	}

	seen := make(map[string]struct{}, len(is.IDIndexes))
	j := 0
	for _, v := range is.IDIndexes {
		u := &ulid.ULID{}
		err = u.UnmarshalBinary(v)
		if err != nil {
			return err
		}

		us := u.String()
		if _, ok := seen[us]; ok {
			continue
		}

		seen[us] = emptyStruct
		is.IDIndexes[j] = v
		j++
	}

	is.IDIndexes = is.IDIndexes[:j]

	if _, ok := seen[sx.id.String()]; !ok {
		is.IDIndexes = append(is.IDIndexes, sx.idBuf)
	}

	IDxsBuf, err = proto.Marshal(is)
	if err != nil {
		return nil
	}

	return b.Put(idBuf, IDxsBuf)
}

func (sx *SetTx) convertSliceToBucket(
	idBuf []byte, is *schema.IndexSlice, b *bolt.Bucket) error {
	var err error
	err = b.Delete(idBuf)
	if err != nil {
		return err
	}

	//deduplicate
	m := make(map[string][]byte)
	for _, v := range is.IDIndexes {
		u := &ulid.ULID{}
		err = u.UnmarshalBinary(v)
		if err != nil {
			return err
		}

		m[u.String()] = v
	}

	iValBuc, err := b.CreateBucketIfNotExists(idBuf)
	if err != nil {
		return err
	}

	for _, v := range m {
		err = iValBuc.Put(v, emptySlice)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sx *SetTx) loadPayload() error {
	sx.err = sx.id.MarshalBinaryTo(sx.idBuf)
	if sx.err != nil {
		return sx.err
	}

	sx.err = sx.partition.store.View(func(tx *bolt.Tx) error {
		sx.payloadBuf = tx.Bucket(rootBucket).Get(sx.idBuf)
		if sx.payloadBuf != nil {
			return proto.Unmarshal(sx.payloadBuf, sx.payload)
		}

		return nil
	})

	return sx.err
}

func (sx *SetTx) writePayload(n string, m proto.Message) {
	if sx.err != nil {
		return
	}

	var pBuf []byte
	pBuf, sx.err = proto.Marshal(m)
	if sx.err != nil {
		return
	}

	var pn uint64
	pn, sx.err = sx.set.instance.GetPropIndex(n)
	if sx.err != nil {
		return
	}

	sx.payload.Values[pn] = pBuf
}

func (sx *SetTx) readPayload(n string, m proto.Message) {
	if sx.err != nil {
		return
	}

	var pn uint64
	pn, sx.err = sx.set.instance.GetPropIndex(n)
	if sx.err != nil {
		return
	}

	val, ok := sx.payload.Values[pn]
	if !ok {
		return
	}

	sx.err = proto.Unmarshal(val, m)
}
