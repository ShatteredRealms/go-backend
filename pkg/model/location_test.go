package model_test

import (
	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
)

var _ = Describe("Location model", func() {
	validateLocation := (func(location *model.Location, pb *pb.Location) {
		Expect(pb.World).To(Equal(location.World))
		Expect(pb.X).To(Equal(location.X))
		Expect(pb.Y).To(Equal(location.Y))
		Expect(pb.Z).To(Equal(location.Z))
		Expect(pb.Roll).To(Equal(location.Roll))
		Expect(pb.Pitch).To(Equal(location.Pitch))
		Expect(pb.Yaw).To(Equal(location.Yaw))
	})

	Describe("ToPb", func() {
		It("should convert single location to protobuf and retain all fields", func() {
			location := &model.Location{}
			Expect(faker.FakeData(location)).To(Succeed())
			out := location.ToPb()
			validateLocation(location, out)
		})
	})

	Describe("ToPb", func() {
		It("should convert single location to protobuf and retain all fields", func() {
			location := &pb.Location{}
			Expect(faker.FakeData(location)).To(Succeed())
			out := model.LocationFromPb(location)
			validateLocation(out, location)
		})
	})
})
