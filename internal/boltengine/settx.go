package boltengine

import (
	"dbee/endian"
	"dbee/internal/boltengine/schema"
	"dbee/store"
	"time"

	proto "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/oklog/ulid"
	bolt "go.etcd.io/bbolt"
)

type SetTx struct {
	id              ulid.ULID
	idBuf           []byte
	set             *Set
	partition       *Partition
	payload         *schema.Payload
	payloadBuf      []byte
	err             error
	onDisk          bool
	indexableUint   map[uint64]uint64
	indexableNint   map[uint64]uint64
	indexableString map[uint64]string
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
	if !sx.onDisk {
		return nil
	}

	sx.err = sx.partition.store.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(rootBucket).Delete(sx.idBuf)
	})

	if sx.err != nil {
		return sx.err
	}

	sx.onDisk = false
	return nil
}

func (sx *SetTx) OnDisk() bool {
	return sx.onDisk
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

	if sx.err == nil {
		sx.onDisk = true
	}

	return sx.err
}

func (sx *SetTx) commitIndex(tx *bolt.Tx) error {
	err := sx.indexString(tx)
	if err != nil {
		return err
	}

	err = sx.indexUint(tx)
	if err != nil {
		return err
	}

	err = sx.indexNint(tx)

	return err
}

func (sx *SetTx) indexString(tx *bolt.Tx) error {
	if len(sx.indexableString) < 1 {
		return nil
	}

	sb := tx.Bucket(indexBucket).Bucket(indexStringBucket)

	for k, v := range sx.indexableString {
		kb, err := sb.CreateBucketIfNotExists(endian.I64toB(k))
		if err != nil {
			return err
		}

		vb, err := kb.CreateBucketIfNotExists([]byte(v))
		if err != nil {
			return err
		}

		err = vb.Put(sx.idBuf, emptySlice)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sx *SetTx) indexUint(tx *bolt.Tx) error {
	if len(sx.indexableUint) < 1 {
		return nil
	}

	ub := tx.Bucket(indexBucket).Bucket(indexUintBucket)

	for k, v := range sx.indexableUint {
		kb, err := ub.CreateBucketIfNotExists(endian.I64toB(k))
		if err != nil {
			return err
		}

		vb, err := kb.CreateBucketIfNotExists(endian.I64toB(v))
		if err != nil {
			return err
		}

		err = vb.Put(sx.idBuf, emptySlice)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sx *SetTx) indexNint(tx *bolt.Tx) error {
	if len(sx.indexableNint) < 1 {
		return nil
	}

	nb := tx.Bucket(indexBucket).Bucket(indexNintBucket)

	for k, v := range sx.indexableNint {
		kb, err := nb.CreateBucketIfNotExists(endian.I64toB(k))
		if err != nil {
			return nil
		}

		vb, err := kb.CreateBucketIfNotExists(endian.I64toB(v))
		if err != nil {
			return nil
		}

		err = vb.Put(sx.idBuf, emptySlice)
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
			sx.onDisk = true
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

// writeIndexablePayload will write the prop payload to sx.
// returns uint64 greater than zero if the prop is indexable;
func (sx *SetTx) writeIndexablePayload(n string, m proto.Message) uint64 {

	if sx.err != nil {
		return 0
	}

	var pBuf []byte
	pBuf, sx.err = proto.Marshal(m)
	if sx.err != nil {
		return 0
	}

	var pn uint64
	pn, sx.err = sx.set.instance.GetPropIndex(n)
	if sx.err != nil {
		return 0
	}

	sx.payload.Values[pn] = pBuf

	if sx.set.idxs.indexable(pn) {
		return pn
	}

	return 0
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
