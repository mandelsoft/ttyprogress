package blocks_test

import (
	"bytes"

	. "github.com/mandelsoft/goutils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UIBlock Test Environment", func() {
	var blocks *UIBlocks
	var buf *bytes.Buffer

	BeforeEach(func() {
		buf = bytes.NewBuffer(nil)
		blocks = blocks.New(buf)
	})

	It("assigns block", func() {
		b := blocks.NewBlock(3)

		Expect(b.Write([]byte("test\n"))).To(Equal(5))
		MustBeSuccessful(blocks.AddBlock(b))
		ExpectError(blocks.AddBlock(b)).To(Equal(blocks.ErrAlreadyAssigned))
		MustBeSuccessful(blocks.Flush())
		Expect(buf.String()).To(Equal("test\n"))
	})
})
