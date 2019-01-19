package boltengine

import (
	"reflect"
)

type indexMap struct {
	name  map[uint64]string
	index map[string]uint64
}

func (i *indexMap) indexable(propID uint64) bool {
	_, ok := i.name[propID]
	return ok
}

func (i *indexMap) add(propID uint64, propName string) {
	i.name[propID] = propName
	i.index[propName] = propID
}

func (i *indexMap) indexes() []string {
	keys := reflect.ValueOf(i.index).MapKeys()
	ks := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		ks[i] = keys[i].String()
	}

	return ks
}
