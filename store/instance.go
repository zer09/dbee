package store

import "time"

type Instance interface {
	// Close the instance.
	Close() error

	// Dir of the instance.
	Dir() string

	// GetPropName get the property name using an index number.
	GetPropName(index uint64) (string, error)

	// GetPropIndex will get the property index using the property name.
	// If the name is not indexed then it will store the index and property name.
	GetPropIndex(name string) (uint64, error)

	// Set will get an instance of set.
	// name is the name of the set.
	Set(name string) (Set, error)
}

type Set interface {
	// Name of the set.
	Name() string

	// Partitions list.
	Partitions() []string

	// Partition will get a partition by name.
	Partition(name string) (Partition, error)

	// Index will create an index for the set.
	// propName will be the name of the property.
	// Every partition will have it's own indexing.
	// for example if you index the property of name, then all name property
	// created on partition 1 will on partition 1.
	// If the property is index after there is already data, then the index will
	// re-scan the set to find data to be indexed.
	Index(propName string) error

	// ListIndexes will list all indexed property.
	ListIndexes() []string

	// Get set transaction
	// The id is optional and will a SetTx id, that was auto generated.
	// If you provide the id, it means that the data is already existed,
	// otherwise just leave it and it will be treated as a new data.
	// This will use the default partition.
	Get(id ...string) (SetTx, error)
}

type Partition interface {
	// Name of the partition
	Name() string

	// Get set transaction
	// The id is optional and will a SetTx id, that was auto generated.
	// If you provide the id, it means that the data is already existed,
	// otherwise just leave it and it will be treated as a new data.
	Get(id ...string) (SetTx, error)
	One(propName string, value interface{}) (SetTx, error)
}

type SetTx interface {
	// ID of the transaction.
	ID() string

	// CreatedOn the time when the item is created.
	CreatedOn() time.Time

	// LastUpdate the time when the item is updated.
	// This will include when the item is soft deleted.
	LastUpdate() time.Time

	// Partition on where this transaction belongs to.
	Partition() Partition

	// IsSoftDeleted the item.
	IsSoftDeleted() bool

	// Delete the data softly.
	// you still need to execute commit before the item will be deleted.
	Delete()

	// HardDelete the data.
	// This will auto commit the delete.
	HardDelete() error

	// Commit the the item.
	Commit() error

	// OnDisk identify if the data is saved in disk
	OnDisk() bool

	Wint(n string, v int64)
	Wuint(n string, v uint64)
	Wbool(n string, v bool)
	Wstring(n string, v string)
	Wbytes(n string, v []byte)

	Rint(n string) int64
	Ruint(n string) uint64
	Rbool(n string) bool
	Rstring(n string) string
	Rbytes(n string) []byte
}
