package boltengine

import (
	"dbee/internal/boltengine/schema"
)

func (sx *SetTx) Wfloat(n string, v float32) {
	sx.writePayload(n, &schema.PayloadFloat{Value: v})
}

func (sx *SetTx) Wdouble(n string, v float64) {
	sx.writePayload(n, &schema.PayloadDouble{Value: v})
}

func (sx *SetTx) Wint(n string, v int64) {
	sx.writePayload(n, &schema.PayloadInt64{Value: v})
}

func (sx *SetTx) Wsint(n string, v int64) {
	sx.writePayload(n, &schema.PayloadSint64{Value: v})
}

func (sx *SetTx) Wuint(n string, v uint64) {
	sx.writePayload(n, &schema.PayloadUint64{Value: v})
}

func (sx *SetTx) Wbool(n string, v bool) {
	sx.writePayload(n, &schema.PayloadBool{Value: v})
}

func (sx *SetTx) Wstring(n string, v string) {
	sx.writePayload(n, &schema.PayloadString{Value: v})
}

func (sx *SetTx) Wbytes(n string, v []byte) {
	sx.writePayload(n, &schema.PayloadBytes{Value: v})
}
