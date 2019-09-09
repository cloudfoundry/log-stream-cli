package presentation_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPresentation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Presentation Suite")
}
