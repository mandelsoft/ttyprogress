package ttyprogress_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTTYProgress(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TTYProgress Test Suite")
}
