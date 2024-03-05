package repository_test

import (
	"context"
	"time"

	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"github.com/ShatteredRealms/go-backend/pkg/model"
)

var _ = Describe("Chat repository", func() {

	createChannel := func() *model.ChatChannel {
		channel := &model.ChatChannel{
			Model: gorm.Model{
				ID:        0,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: gorm.DeletedAt{},
			},
			Name:      faker.Username(),
			Dimension: faker.Username(),
		}

		out, err := chatRepo.CreateChannel(nil, channel)
		Expect(err).To(BeNil())
		Expect(out).NotTo(BeNil())
		Expect(out.ID).To(BeEquivalentTo(channel.ID))

		return out
	}

	Describe("AllChannels", func() {
		allChannels := (func(ctx context.Context) {
			Expect(createChannel()).NotTo(BeNil())
			out, err := chatRepo.AllChannels(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(out) >= 1).To(BeTrue())
		})

		It("should not error with invalid input", func() {
			allChannels(nil)
		})

		It("should work", func() {
			allChannels(context.Background())
		})
	})

	Describe("FindChanelById", func() {
		Context("invalid input", func() {
			It("should not error with invalid context", func() {
				out, err := chatRepo.FindChannelById(nil, 0)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error on invalid id", func() {
				out, err := chatRepo.FindChannelById(context.Background(), 1e19)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})

		Context("valid input", func() {
			It("should return matching chat channel", func() {
				channel := createChannel()
				out, err := chatRepo.FindChannelById(context.Background(), channel.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out.ID).To(Equal(channel.ID))
				Expect(out.Name).To(Equal(channel.Name))
			})
		})
	})

	Describe("CreateChannel", func() {
		It("should work", func() {
			channel := createChannel()
			channel.Name += "a"
			outChannel, err := chatRepo.CreateChannel(nil, channel)
			Expect(err).To(HaveOccurred())

			channel.ID = 0
			outChannel, err = chatRepo.CreateChannel(nil, channel)
			Expect(err).NotTo(HaveOccurred())
			Expect(outChannel).NotTo(BeNil())
			Expect(outChannel.ID).NotTo(BeEquivalentTo(0))
		})
	})

	Describe("UpdateChannel", func() {
		It("should work", func() {
			channel := createChannel()
			channel.Name += "a"
			outChannel, err := chatRepo.UpdateChannel(nil, channel)
			Expect(err).NotTo(HaveOccurred())
			Expect(outChannel).NotTo(BeNil())
			Expect(outChannel.ID).To(BeEquivalentTo(channel.ID))
			Expect(outChannel.Name).To(Equal(channel.Name))
		})
	})

	Describe("Delete", func() {
		It("should work with valid channel", func() {
			channel := createChannel()
			Expect(chatRepo.DeleteChannel(nil, channel)).To(Succeed())
			out, err := chatRepo.FindChannelById(context.Background(), channel.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(BeNil())
		})

		It("should not error with non-existant channel", func() {
			channel := createChannel()
			channel.ID += 1000
			Expect(chatRepo.DeleteChannel(nil, channel)).To(Succeed())
			out, err := chatRepo.FindChannelById(context.Background(), channel.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(BeNil())
		})

		It("should not error with invalid channel id", func() {
			channel := createChannel()
			channel.ID = 1e19
			Expect(chatRepo.DeleteChannel(nil, channel)).NotTo(Succeed())
		})
	})

	Describe("FullDelete", func() {
		It("should work with valid channel", func() {
			channel := createChannel()
			Expect(chatRepo.FullDeleteChannel(nil, channel)).To(Succeed())
			out, err := chatRepo.FindChannelById(context.Background(), channel.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(BeNil())
		})

		It("should not error with non-existant channel", func() {
			channel := createChannel()
			channel.ID += 2000
			Expect(chatRepo.FullDeleteChannel(nil, channel)).To(Succeed())
			out, err := chatRepo.FindChannelById(context.Background(), channel.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(BeNil())
		})

		It("should not error with invalid channel id", func() {
			channel := createChannel()
			channel.ID = 1e19
			Expect(chatRepo.FullDeleteChannel(nil, channel)).NotTo(Succeed())
		})
	})

	Describe("FindDeleteWithName", func() {
		var channel *model.ChatChannel
		BeforeEach(func() {
			channel = createChannel()
			Expect(chatRepo.DeleteChannel(nil, channel)).To(Succeed())
			out, err := chatRepo.FindChannelById(context.Background(), channel.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(BeNil())
		})

		Context("invalid input", func() {
			It("should work with invalid ctx", func() {
				out, err := chatRepo.FindDeletedWithName(nil, channel.Name)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out.ID).To(Equal(channel.ID))
			})

			It("should not error with non-existant channel", func() {
				channel.Name += "a"
				out, err := chatRepo.FindDeletedWithName(context.Background(), channel.Name)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})

		Context("valid input", func() {
			It("should work with valid channel", func() {
				out, err := chatRepo.FindDeletedWithName(context.Background(), channel.Name)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out.ID).To(Equal(channel.ID))
			})
		})
	})

	Describe("Authorization", func() {
		It("should allow adding and removing authorization", func() {
			channel := createChannel()
			exampleId := uint(1)
			exampleChannelIds := []uint{channel.ID}

			channels, err := chatRepo.AuthorizedChannelsForCharacter(nil, uint(exampleId))
			Expect(err).NotTo(HaveOccurred())
			Expect(channels).To(BeEmpty())

			Expect(chatRepo.ChangeAuthorizationForCharacter(nil, exampleId, exampleChannelIds, true)).To(Succeed())
			channels, err = chatRepo.AuthorizedChannelsForCharacter(nil, uint(exampleId))
			Expect(err).NotTo(HaveOccurred())
			Expect(channels).To(HaveLen(len(exampleChannelIds)))

			Expect(chatRepo.ChangeAuthorizationForCharacter(nil, exampleId, exampleChannelIds, false)).To(Succeed())
			channels, err = chatRepo.AuthorizedChannelsForCharacter(nil, uint(exampleId))
			Expect(err).NotTo(HaveOccurred())
			Expect(channels).To(HaveLen(0))
		})

		It("should not commit changes on error", func() {
			exampleId := uint(1)
			exampleChannelIds := []uint{1, 2, 1e19}
			Expect(chatRepo.ChangeAuthorizationForCharacter(nil, exampleId, exampleChannelIds, true)).NotTo(Succeed())

			channels, err := chatRepo.AuthorizedChannelsForCharacter(nil, uint(exampleId))
			Expect(err).NotTo(HaveOccurred())
			Expect(channels).To(BeEmpty())
		})

		It("should not error on deleting non-existing", func() {
			Expect(chatRepo.ChangeAuthorizationForCharacter(nil, 1, []uint{1, 2}, false)).To(Succeed())
		})
	})
})
