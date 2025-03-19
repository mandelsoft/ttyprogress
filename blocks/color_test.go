package blocks_test

import (
	"github.com/fatih/color"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Color Test Environment", func() {

	c := color.New(color.FgYellow, color.Bold)
	c.EnableColor()
	It("", func() {
		s := c.Sprintf("test")

		Expect(ColorLength([]byte(s))).To(Equal(7))
	})
})
