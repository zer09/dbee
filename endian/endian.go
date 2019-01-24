package endian

import (
	"encoding/binary"
	"math"
)

// I32toB converts int32 to 32 byte big endian.
func I32toB(v uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, v)
	return b
}

// I64toB converts int64 to 64 byte big endian.
func I64toB(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func Float32toB(v float32) []byte {
	bits := math.Float32bits(v)
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, bits)
	return b
}

func Float64toB(v float64) []byte {
	bits := math.Float64bits(v)
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, bits)
	return b
}

// BtoI32 converts 32 byte binary big endian to uint32.
func BtoI32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

// BtoI64 converts 64 byte binary big endian to uint64.
func BtoI64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func BtoFloat32(b []byte) float32 {
	bits := binary.BigEndian.Uint32(b)
	return math.Float32frombits(bits)
}

func BtoFloat64(b []byte) float64 {
	bits := binary.BigEndian.Uint64(b)
	return math.Float64frombits(bits)
}
