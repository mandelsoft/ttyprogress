package units_test

import (
	"github.com/mandelsoft/ttyprogress/units"
)

var _ = Describe("Units Test Environment", func() {
	Context("bytes", func() {
		It("plain", func() {
			u := units.Bytes()
			Expect(u(1)).To(Equal("1"))
			Expect(u(1023)).To(Equal("1023"))
			Expect(u(1025)).To(Equal("1 KB"))
			Expect(u(1023 * 1024 * 1024)).To(Equal("1023 MB"))
			Expect(u(1023 * 1024 * 1024 * 1024)).To(Equal("1023 GB"))
		})

		It("scaled", func() {
			u := units.Bytes(1024)
			Expect(u(1)).To(Equal("1 KB"))
			Expect(u(1023)).To(Equal("1023 KB"))
			Expect(u(1025)).To(Equal("1 MB"))
			Expect(u(1023 * 1024 * 1024)).To(Equal("1023 GB"))
			Expect(u(1023 * 1024 * 1024 * 1024)).To(Equal("1023 TB"))
		})
	})

	Context("millimeter", func() {
		It("plain", func() {
			u := units.Millimeter()
			Expect(u(1)).To(Equal("1 mm"))
			Expect(u(999)).To(Equal("999 mm"))
			Expect(u(1001)).To(Equal("1 m"))
			Expect(u(999 * 1000 * 1000)).To(Equal("999 km"))
			Expect(u(999 * 1000 * 1000 * 1000)).To(Equal("999000 km"))
		})
	})
})
