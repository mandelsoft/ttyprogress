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

	It("uses the default view for a view of 0", func() {
		render := func(b *blocks.Block) string {
			b.SetAuto()
			buf := bytes.NewBuffer(nil)
			blks := blocks.New(buf)
			for i := 0; i < 15; i++ {
				b.Write([]byte("line\n"))
			}
			MustBeSuccessful(blks.AddBlock(b))
			MustBeSuccessful(blks.Flush())
			time.Sleep(blocks.MIN_UPDATE_INTERVAL * 2)
			return buf.String()
		}
		Expect(render(blocks.NewBlock(0))).To(Equal(render(blocks.NewBlock())))
		Expect(render(blocks.NewBlock(0))).NotTo(Equal(render(blocks.NewBlock(3))))
	})
})
