package boltengine

import (
	"path/filepath"
	"strings"
)

func sanitizeProp(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func openPartition(part *Partition, dir string) error {
	var err error

	if part.store == nil {
		part.store, err = open(filepath.Join(dir, part.storeName))
		if err != nil {
			return err
		}
	}

	return nil
}
