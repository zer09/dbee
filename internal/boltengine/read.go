package boltengine

import "dbee/internal/boltengine/schema"

func (sx *SetTx) Rfloat(n string) float32 {
	m := &schema.PayloadFloat{}
	sx.readPayload(n, m)
	return m.Value
}

func (sx *SetTx) Rdouble(n string) float64 {
	m := &schema.PayloadDouble{}
	sx.readPayload(n, m)
	return m.Value
}

func (sx *SetTx) Rint(n string) int64 {
	m := &schema.PayloadInt64{}
	sx.readPayload(n, m)
	return m.Value
}

func (sx *SetTx) Rsint(n string) int64 {
	m := &schema.PayloadSint64{}
	sx.readPayload(n, m)
	return m.Value
}

func (sx *SetTx) Ruint(n string) uint64 {
	m := &schema.PayloadUint64{}
	sx.readPayload(n, m)
	return m.Value
}

func (sx *SetTx) Rbool(n string) bool {
	m := &schema.PayloadBool{}
	sx.readPayload(n, m)
	return m.Value
}

func (sx *SetTx) Rstring(n string) string {
	m := &schema.PayloadString{}
	sx.readPayload(n, m)
	return m.Value
}

func (sx *SetTx) Rbytes(n string) []byte {
	m := &schema.PayloadBytes{}
	sx.readPayload(n, m)
	return m.Value
}
