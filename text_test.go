package ttyprogress_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mandelsoft/ttyprogress"
)

var _ = Describe("Text", func() {
	It("uses the default view for 0 lines", func() {
		Expect(ttyprogress.NewText(0).GetView()).To(Equal(ttyprogress.NewText().GetView()))
		Expect(ttyprogress.NewText(5).GetView()).To(Equal(5))
	})
})
