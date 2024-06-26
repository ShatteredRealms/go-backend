package game_test

import (
	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/go-backend/pkg/common"
	"github.com/ShatteredRealms/go-backend/pkg/model/game"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
)

var _ = Describe("Dimension game model", func() {
	var (
		dimension *game.Dimension
	)

	BeforeEach(func() {
		dimension, _ = randomDimensionAndMap()
	})

	Describe("Validation", func() {
		Context("issues", func() {
			It("should error with empty name", func() {
				dimension.Name = ""
				Expect(dimension.ValidateName()).To(MatchError(game.ErrDimensionNameToShort))
			})
			It("should error with to long name", func() {
				dimension.Name = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
				Expect(dimension.ValidateName()).To(MatchError(game.ErrDimensionNameToLong))
			})

			It("should error with invalid location", func() {
				dimension.Location = faker.Email()
				Expect(dimension.ValidateLocation()).To(MatchError(common.ErrInvalidServerLocation))
			})
		})
		It("should not error for valid dimension", func() {
			Expect(dimension.ValidateLocation()).To(Succeed())
			Expect(dimension.ValidateName()).To(Succeed())
		})
	})

	validateDimension := (func(dimension *game.Dimension, pb *pb.Dimension) {
		Expect(pb.Id).To(Equal(dimension.Id.String()))
		Expect(pb.Name).To(Equal(dimension.Name))
		Expect(pb.Version).To(Equal(dimension.Version))
		Expect(pb.Location).To(Equal(dimension.Location))
	})

	Describe("ToPb", func() {
		It("should convert single dimension to protobuf and retain all fields", func() {
			out := dimension.ToPb()
			validateDimension(dimension, out)
		})

		It("should convert array of dimensions to protobuf and retain all fields", func() {
			var dimensions game.Dimensions
			dimensions = make([]*game.Dimension, 10)
			for idx := range dimensions {
				dim, _ := randomDimensionAndMap()
				dimensions[idx] = dim
			}
			out := dimensions.ToPb()
			Expect(out.Dimensions).To(HaveLen(len(dimensions)))
			for idx := range dimensions {
				validateDimension(dimensions[idx], out.Dimensions[idx])
			}
		})
	})

	Describe("GetImageName", func() {
		It("should work", func() {
			Expect(dimension.GetImageName()).To(HaveSuffix(":" + dimension.Version))
			dimension.Version = ""
			Expect(dimension.GetImageName()).To(HaveSuffix(":latest"))
		})
	})
})
