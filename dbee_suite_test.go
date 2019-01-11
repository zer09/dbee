package dbee_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDbee(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dbee Suite")
}
