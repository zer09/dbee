package boltengine

import (
	"dbee/internal/boltengine/schema"
)

func (sx *SetTx) Wfloat(n string, v float32) {
	if pn := sx.writePayload(n, &schema.PayloadFloat{Value: v}, v); pn > 0 {
		if sx.indexableFloat == nil {
			sx.indexableFloat = make(map[uint64]float32)
		}

		sx.indexableFloat[pn] = v
	}
}

func (sx *SetTx) Wdouble(n string, v float64) {
	if pn := sx.writePayload(n, &schema.PayloadDouble{Value: v}, v); pn > 0 {
		if sx.indexableDouble == nil {
			sx.indexableDouble = make(map[uint64]float64)
		}

		sx.indexableDouble[pn] = v
	}
}

func (sx *SetTx) Wint(n string, v int64) {
	if pn := sx.writePayload(n, &schema.PayloadInt64{Value: v}, v); pn > 0 {
		if sx.indexableInt == nil {
			sx.indexableInt = make(map[uint64]int64)
		}

		sx.indexableInt[pn] = v
	}
}

func (sx *SetTx) Wsint(n string, v int64) {
	if pn := sx.writePayload(n, &schema.PayloadSint64{Value: v}, v); pn > 0 {
		if sx.indexableSint == nil {
			sx.indexableSint = make(map[uint64]int64)
		}

		sx.indexableSint[pn] = v
	}
}

func (sx *SetTx) Wuint(n string, v uint64) {
	if pn := sx.writePayload(n, &schema.PayloadUint64{Value: v}, v); pn > 0 {
		if sx.indexableUint == nil {
			sx.indexableUint = make(map[uint64]uint64)
		}

		sx.indexableUint[pn] = v
	}
}

func (sx *SetTx) Wbool(n string, v bool) {
	if pn := sx.writePayload(n, &schema.PayloadBool{Value: v}, v); pn > 0 {
		if sx.indexableBool == nil {
			sx.indexableBool = make(map[uint64]bool)
		}

		sx.indexableBool[pn] = v
	}
}

func (sx *SetTx) Wstring(n string, v string) {
	if pn := sx.writePayload(n, &schema.PayloadString{Value: v}, v); pn > 0 {
		if sx.indexableString == nil {
			sx.indexableString = make(map[uint64]string)
		}

		sx.indexableString[pn] = v
	}
}

func (sx *SetTx) Wbytes(n string, v []byte) {
	sx.writePayload(n, &schema.PayloadBytes{Value: v}, v)
}
