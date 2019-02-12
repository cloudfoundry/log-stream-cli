package rlp_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRlp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rlp Suite")
}
