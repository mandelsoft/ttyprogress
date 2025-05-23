package blocks_test

import (
	"bytes"
	"time"

	. "github.com/mandelsoft/goutils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mandelsoft/ttyprogress/blocks"
)

var _ = Describe("Blocks Test Environment", func() {
	var blks *blocks.Blocks
	var buf *bytes.Buffer

	BeforeEach(func() {
		buf = bytes.NewBuffer(nil)
		blks = blocks.New(buf)
	})

	It("assigns block", func() {
		b := blocks.NewBlock(3).SetAuto()

		Expect(b.Write([]byte("test\n"))).To(Equal(5))
		MustBeSuccessful(blks.AddBlock(b))
		ExpectError(blks.AddBlock(b)).To(Equal(blocks.ErrAlreadyAssigned))
		MustBeSuccessful(blks.Flush())
		time.Sleep(blocks.MIN_UPDATE_INTERVAL * 2)

		s := buf.String()
		Expect(s).To(Equal("test\n"))
	})
})
