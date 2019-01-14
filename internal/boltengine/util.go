package boltengine

import "strings"

func sanitizeProp(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
