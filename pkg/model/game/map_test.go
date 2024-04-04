package game_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/go-backend/pkg/model/game"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
)

var _ = Describe("Map game model", func() {

	var (
		m *game.Map
	)

	BeforeEach(func() {
		_, m = randomDimensionAndMap()
	})

	validateMap := (func(m *game.Map, pb *pb.Map) {
		Expect(pb.Id).To(Equal(m.Id.String()))
		Expect(pb.Name).To(Equal(m.Name))
		Expect(pb.Path).To(Equal(m.Path))
		Expect(pb.MaxPlayers).To(Equal(m.MaxPlayers))
		Expect(pb.Instanced).To(Equal(m.Instanced))
	})

	Describe("ToPb", func() {
		It("should convert single m to protobuf and retain all fields", func() {
			out := m.ToPb()
			validateMap(m, out)
		})

		It("should convert array of ms to protobuf and retain all fields", func() {
			var maps game.Maps
			maps = make([]*game.Map, 10)
			for idx := range maps {
				_, m := randomDimensionAndMap()
				maps[idx] = m
			}

			out := maps.ToPb()
			Expect(out).To(HaveLen(len(maps)))
			for idx := range maps {
				validateMap(maps[idx], out[idx])
			}
		})
	})
})
