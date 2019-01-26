package boltengine

import "dbee/internal/boltengine/schema"

func (sx *SetTx) Rint(n string) int64 {
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
