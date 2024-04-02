package game_test

import (
	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/go-backend/pkg/model/game"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
)

var _ = Describe("Location game model", func() {
	validateLocation := (func(location *game.Location, pb *pb.Location) {
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
			location := &game.Location{}
			Expect(faker.FakeData(location)).To(Succeed())
			out := location.ToPb()
			validateLocation(location, out)
		})
	})

	Describe("ToPb", func() {
		It("should convert single location to protobuf and retain all fields", func() {
			location := &pb.Location{}
			Expect(faker.FakeData(location)).To(Succeed())
			out := game.LocationFromPb(location)
			validateLocation(out, location)
		})
	})
})
