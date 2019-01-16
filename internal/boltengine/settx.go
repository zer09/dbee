package boltengine

import (
	"dbee/internal/boltengine/schema"
	"dbee/store"
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
		return tx.Bucket(rootBucket).Put(sx.idBuf, sx.payloadBuf)
	})

	return sx.err
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
