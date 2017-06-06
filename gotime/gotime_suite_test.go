package gotime_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGotime(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gotime Suite")
}
