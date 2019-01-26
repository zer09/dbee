package boltengine

import (
	"dbee/internal/boltengine/schema"
)

func (sx *SetTx) Wint(n string, v int64) {
	pn := sx.writeIndexablePayload(n, &schema.PayloadSint64{Value: v})

	if pn < 1 {
		return
	}

	if v < 0 {
		if sx.indexableNint == nil {
			sx.indexableNint = make(map[uint64]uint64)
		}

		sx.indexableNint[pn] = uint64(abs(v))
		return
	}

	if sx.indexableUint == nil {
		sx.indexableUint = make(map[uint64]uint64)
	}

	sx.indexableUint[pn] = uint64(v)
}

func (sx *SetTx) Wuint(n string, v uint64) {
	pn := sx.writeIndexablePayload(n, &schema.PayloadUint64{Value: v})

	if pn < 1 {
		return
	}

	if sx.indexableUint == nil {
		sx.indexableUint = make(map[uint64]uint64)
	}

	sx.indexableUint[pn] = v
}

func (sx *SetTx) Wstring(n string, v string) {
	pn := sx.writeIndexablePayload(n, &schema.PayloadString{Value: v})

	if pn < 1 {
		return
	}

	if sx.indexableString == nil {
		sx.indexableString = make(map[uint64]string)
	}

	sx.indexableString[pn] = v
}

func (sx *SetTx) Wbool(n string, v bool) {
	sx.writePayload(n, &schema.PayloadBool{Value: v})
}

func (sx *SetTx) Wbytes(n string, v []byte) {
	sx.writePayload(n, &schema.PayloadBytes{Value: v})
}
