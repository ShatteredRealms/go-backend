package helpers_test

import (
	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/go-backend/pkg/helpers"
)

var _ = Describe("Conversion helpers", func() {
	Describe("ArrayOfUint64ToUint", func() {
		It("should work", func() {
			ints, err := faker.RandomInt(10, 100, 1)
			Expect(err).To(Succeed())
			in := make([]uint64, ints[0])
			for idx := range in {
				in[idx] = uint64(faker.RandomUnixTime())
			}
			out := *helpers.ArrayOfUint64ToUint(&in)

			for idx := range out {
				Expect(out[idx]).To(Equal(uint(in[idx])), "each index should match")
			}
		})
	})
})
