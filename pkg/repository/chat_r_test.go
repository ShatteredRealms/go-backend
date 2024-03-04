package repository_test

import (
	"context"
	"time"

	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus/hooks/test"
	"gorm.io/gorm"

	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	testdb "github.com/ShatteredRealms/go-backend/test/db"
)

var _ = Describe("Chat repository", func() {
	var (
		hook *test.Hook

		db          *gorm.DB
		dbCloseFunc func()

		repo repository.ChatRepository

		channel *model.ChatChannel
	)

	BeforeEach(func() {
		log.Logger, hook = test.NewNullLogger()

		db, dbCloseFunc = testdb.SetupGormWithDocker()
		Expect(db).NotTo(BeNil())

		repo = repository.NewChatRepository(db)
		Expect(repo).NotTo(BeNil())
		Expect(repo.Migrate(context.Background())).To(Succeed())

		channel = &model.ChatChannel{
			Model: gorm.Model{
				ID:        0,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: gorm.DeletedAt{},
			},
			Name:      faker.Username(),
			Dimension: faker.Username(),
		}

		outChannel, err := repo.CreateChannel(nil, channel)
		Expect(err).To(BeNil())
		Expect(outChannel).NotTo(BeNil())
		Expect(outChannel.ID).To(BeEquivalentTo(1))

		hook.Reset()
	})

	Describe("AllChannels", func() {
		allChannels := (func(ctx context.Context) {
			out, err := repo.AllChannels(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(HaveLen(1))
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
				out, err := repo.FindChannelById(nil, 0)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error on invalid id", func() {
				out, err := repo.FindChannelById(context.Background(), 1e19)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})

		Context("valid input", func() {
			It("should return matching chat channel", func() {
				out, err := repo.FindChannelById(context.Background(), channel.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out.ID).To(Equal(channel.ID))
				Expect(out.Name).To(Equal(channel.Name))
			})
		})
	})

	Describe("CreateChannel", func() {
		It("should work", func() {
			channel.Name += "a"
			outChannel, err := repo.CreateChannel(nil, channel)
			Expect(err).To(HaveOccurred())

			channel.ID = 0
			outChannel, err = repo.CreateChannel(nil, channel)
			Expect(err).NotTo(HaveOccurred())
			Expect(outChannel).NotTo(BeNil())
			Expect(outChannel.ID).To(BeEquivalentTo(2))
		})
	})

	Describe("UpdateChannel", func() {
		It("should work", func() {
			channel.Name += "a"
			outChannel, err := repo.UpdateChannel(nil, channel)
			Expect(err).NotTo(HaveOccurred())
			Expect(outChannel).NotTo(BeNil())
			Expect(outChannel.ID).To(BeEquivalentTo(channel.ID))
			Expect(outChannel.Name).To(Equal(channel.Name))
		})
	})

	Describe("Delete", func() {
		It("should work with valid channel", func() {
			Expect(repo.DeleteChannel(nil, channel)).To(Succeed())
			out, err := repo.AllChannels(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(BeEmpty())
		})

		It("should not error with non-existant channel", func() {
			channel.ID += 1
			Expect(repo.DeleteChannel(nil, channel)).To(Succeed())
			out, err := repo.AllChannels(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(HaveLen(1))
		})

		It("should not error with invalid channel id", func() {
			channel.ID = 1e19
			Expect(repo.DeleteChannel(nil, channel)).NotTo(Succeed())
			out, err := repo.AllChannels(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(HaveLen(1))
		})
	})

	Describe("FullDelete", func() {
		It("should work with valid channel", func() {
			Expect(repo.FullDeleteChannel(nil, channel)).To(Succeed())
			out, err := repo.AllChannels(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(BeEmpty())
		})

		It("should not error with non-existant channel", func() {
			channel.ID += 1
			Expect(repo.FullDeleteChannel(nil, channel)).To(Succeed())
			out, err := repo.AllChannels(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(HaveLen(1))
		})

		It("should not error with invalid channel id", func() {
			channel.ID = 1e19
			Expect(repo.FullDeleteChannel(nil, channel)).NotTo(Succeed())
			out, err := repo.AllChannels(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(HaveLen(1))
		})
	})

	Describe("FindDeleteWithName", func() {
		BeforeEach(func() {
			Expect(repo.DeleteChannel(nil, channel)).To(Succeed())
			out, err := repo.AllChannels(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(BeEmpty())
		})

		Context("invalid input", func() {
			It("should work with invalid ctx", func() {
				out, err := repo.FindDeletedWithName(nil, channel.Name)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out.ID).To(Equal(channel.ID))
			})

			It("should not error with non-existant channel", func() {
				channel.Name += "a"
				out, err := repo.FindDeletedWithName(context.Background(), channel.Name)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})

		Context("valid input", func() {
			It("should work with valid channel", func() {
				out, err := repo.FindDeletedWithName(context.Background(), channel.Name)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
				Expect(out.ID).To(Equal(channel.ID))
			})
		})
	})

	Describe("Authorization", func() {
		It("should allow adding and removing authorization", func() {
			exampleId := uint(1)
			exampleChannelIds := []uint{channel.ID}

			channels, err := repo.AuthorizedChannelsForCharacter(nil, uint(exampleId))
			Expect(err).NotTo(HaveOccurred())
			Expect(channels).To(BeEmpty())

			Expect(repo.ChangeAuthorizationForCharacter(nil, exampleId, exampleChannelIds, true)).To(Succeed())
			channels, err = repo.AuthorizedChannelsForCharacter(nil, uint(exampleId))
			Expect(err).NotTo(HaveOccurred())
			Expect(channels).To(HaveLen(len(exampleChannelIds)))

			Expect(repo.ChangeAuthorizationForCharacter(nil, exampleId, exampleChannelIds, false)).To(Succeed())
			channels, err = repo.AuthorizedChannelsForCharacter(nil, uint(exampleId))
			Expect(err).NotTo(HaveOccurred())
			Expect(channels).To(HaveLen(0))
		})

		It("should not commit changes on error", func() {
			exampleId := uint(1)
			exampleChannelIds := []uint{1, 2, 1e19}
			Expect(repo.ChangeAuthorizationForCharacter(nil, exampleId, exampleChannelIds, true)).NotTo(Succeed())

			channels, err := repo.AuthorizedChannelsForCharacter(nil, uint(exampleId))
			Expect(err).NotTo(HaveOccurred())
			Expect(channels).To(BeEmpty())
		})

		It("should not error on deleting non-existing", func() {
			Expect(repo.ChangeAuthorizationForCharacter(nil, 1, []uint{1, 2}, false)).To(Succeed())
		})
	})

	AfterEach(func() {
		dbCloseFunc()
	})
})
