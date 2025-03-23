package blocks_test

import (
	"bytes"

	. "github.com/mandelsoft/goutils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mandelsoft/ttyprogress/blocks"
)

var _ = Describe("UIBlock Test Environment", func() {
	var blks *blocks.Blocks
	var buf *bytes.Buffer

	BeforeEach(func() {
		buf = bytes.NewBuffer(nil)
		blks = blocks.New(buf)
	})

	It("assigns block", func() {
		b := blocks.NewBlock(3)

		Expect(b.Write([]byte("test\n"))).To(Equal(5))
		MustBeSuccessful(blks.AddBlock(b))
		ExpectError(blks.AddBlock(b)).To(Equal(blocks.ErrAlreadyAssigned))
		MustBeSuccessful(blks.Flush())
		Expect(buf.String()).To(Equal("test\n"))
	})
})
