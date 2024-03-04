package model_test

import (
	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
)

var _ = Describe("Chat model", func() {
	Describe("ToPb", func() {
		validateCharacter := (func(channel *model.ChatChannel, pb *pb.ChatChannel) {
			Expect(pb.Id).To(Equal(uint64(channel.ID)))
			Expect(pb.Name).To(Equal(channel.Name))
			Expect(pb.Dimension).To(Equal(channel.Dimension))
		})

		It("should convert a single channel", func() {
			channel := &model.ChatChannel{}
			Expect(faker.FakeData(channel)).To(Succeed())

			out := channel.ToPb()
			validateCharacter(channel, out)
		})

		It("should convert channel arrays", func() {
			var channels model.ChatChannels
			channels = make([]*model.ChatChannel, 10)
			for idx := range channels {
				channels[idx] = &model.ChatChannel{}
				faker.FakeData(channels[idx])
			}
			out := channels.ToPb()
			Expect(out.Channels).To(HaveLen(len(channels)))
			for idx := range channels {
				validateCharacter(channels[idx], out.Channels[idx])
			}
		})
	})
})
