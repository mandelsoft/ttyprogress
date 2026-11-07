package ttyprogress_test

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mandelsoft/ttyprogress"
)

var _ = Describe("Bar", func() {
	It("uses the configured total for a total of 0", func() {
		var buf bytes.Buffer
		p := ttyprogress.For(&buf)
		defer p.Close()
		bar, err := ttyprogress.NewBar().SetTotal(50).AddWithTotal(p, 0)
		Expect(err).To(Succeed())
		bar.Set(25)
		Expect(bar.CompletedPercent()).To(Equal(50.0))
		bar, err = ttyprogress.NewBar().SetTotal(50).AddWithTotal(p, 100)
		Expect(err).To(Succeed())
		bar.Set(25)
		Expect(bar.CompletedPercent()).To(Equal(25.0))
	})
})
