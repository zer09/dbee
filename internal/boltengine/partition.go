package boltengine

import bolt "go.etcd.io/bbolt"

type partition struct {
	// Name of the partition.
	Name string
	// Index file location of the partition.
	Index string
	// Store file location of the partition.
	// This will be the storage of the data.
	Store  string
	opened bool
	index  *bolt.DB
	store  *bolt.DB
}

func (p *partition) GetIndex() *bolt.DB {
	return p.index
}

func (p *partition) GetStore() *bolt.DB {
	return p.store
}
