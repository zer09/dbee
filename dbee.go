package dbee

import (
	"dbee/internal/boltengine"
	"dbee/store"
)

type Engine int

//go:generate stringer -type=Engine
const (
	Bolt Engine = iota + 1
)

func Open(dir string, engine Engine) (store.Instance, error) {
	return boltengine.New(dir)
}
