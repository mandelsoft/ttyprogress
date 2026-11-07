package blocks_test

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mandelsoft/ttyprogress/blocks"
)

var _ = Describe("Blocks", func() {
	It("provides a disabled output context for writers, which are no terminal", func() {
		blks := blocks.New(bytes.NewBuffer(nil))
		Expect(blks.GetTTYGontext()).NotTo(BeNil())
		Expect(blks.IsColorsEnabled()).To(BeFalse())
		blks.EnableColors()
		Expect(blks.IsColorsEnabled()).To(BeTrue())
	})
})
