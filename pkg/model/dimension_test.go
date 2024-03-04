package model_test

import (
	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
)

var _ = Describe("Dimension model", func() {
	var (
		dimension = &model.Dimension{}
	)

	BeforeEach(func() {
		Expect(faker.FakeData(dimension)).To(Succeed())
		dimension.Name = "name"
		dimension.Location = "us-central"
	})
	Describe("Validation", func() {
		Context("issues", func() {
			It("should error with empty name", func() {
				dimension.Name = ""
				Expect(dimension.ValidateName()).To(MatchError(model.ErrDimensionNameToShort))
			})
			It("should error with to long name", func() {
				dimension.Name = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
				Expect(dimension.ValidateName()).To(MatchError(model.ErrDimensionNameToLong))
			})

			It("should error with invalid location", func() {
				dimension.Location = faker.Email()
				Expect(dimension.ValidateLocation()).To(MatchError(model.ErrInvalidServerLocation))
			})
		})
		It("should not error for valid dimension", func() {
			Expect(dimension.ValidateLocation()).To(Succeed())
			Expect(dimension.ValidateName()).To(Succeed())
		})
	})

	validateDimension := (func(dimension *model.Dimension, pb *pb.Dimension) {
		Expect(pb.Id).To(Equal(dimension.Id.String()))
		Expect(pb.Name).To(Equal(dimension.Name))
		Expect(pb.Version).To(Equal(dimension.Version))
		Expect(pb.Maps).NotTo(BeEmpty())
		Expect(pb.Location).To(Equal(dimension.Location))
	})

	Describe("ToPb", func() {
		It("should convert single dimension to protobuf and retain all fields", func() {
			out := dimension.ToPb()
			validateDimension(dimension, out)
		})

		It("should convert array of dimensions to protobuf and retain all fields", func() {
			var dimensions model.Dimensions
			dimensions = make([]*model.Dimension, 10)
			for idx := range dimensions {
				dimensions[idx] = &model.Dimension{}
				faker.FakeData(dimensions[idx])
			}
			out := dimensions.ToPb()
			Expect(out.Dimensions).To(HaveLen(len(dimensions)))
			for idx := range dimensions {
				validateDimension(dimensions[idx], out.Dimensions[idx])
			}
		})
	})
})
