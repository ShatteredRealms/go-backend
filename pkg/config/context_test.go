package config_test

import (
	"context"
	"fmt"

	"github.com/ShatteredRealms/go-backend/pkg/common"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/mocks"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus/hooks/test"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Context", func() {
	var (
		mockCtrl         *gomock.Controller
		mockCSC          *mocks.MockCharacterServiceClient
		ctx              context.Context
		srvContext       *config.ServerContext
		target           *pb.CharacterTarget
		characterDetails *pb.CharacterDetails
		hook             *test.Hook
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockCSC = mocks.NewMockCharacterServiceClient(mockCtrl)
		ctx = context.Background()
		target = &pb.CharacterTarget{
			Type: nil,
		}

		randIds, err := faker.RandomInt(1, 100000, 1)
		Expect(err).ShouldNot(HaveOccurred())
		characterDetails = &pb.CharacterDetails{
			Id:   uint64(randIds[0]),
			Name: faker.Name(),
		}

		log.Logger, hook = test.NewNullLogger()
	})

	Context("GetCharacterIdFromTarget", func() {
		Context("with invalid inputs", func() {
			It("should require a target", func() {
				id, err := srvContext.GetCharacterIdFromTarget(ctx, nil)
				Expect(err).Should(HaveOccurred())
				Expect(id).NotTo(Equal(uint(characterDetails.Id)))
			})

			It("should error on invalid target", func() {
				id, err := srvContext.GetCharacterIdFromTarget(ctx, target)
				Expect(err).Should(MatchError(common.ErrHandleRequest.Err()))
				Expect(id).NotTo(Equal(uint(characterDetails.Id)))
				Expect(hook.AllEntries()).To(HaveLen(1))
				Expect(hook.LastEntry().String()).To(ContainSubstring("target type unknown"))
			})
		})

		Context("with valid inputs", func() {
			It("should return existing id if provided", func() {
				target.Type = &pb.CharacterTarget_Id{
					Id: characterDetails.Id,
				}
				id, err := srvContext.GetCharacterIdFromTarget(ctx, target)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(id).To(Equal(uint(target.GetId())))
			})

			Context("looking up character", func() {
				It("should if given a name", func() {
					target.Type = &pb.CharacterTarget_Name{
						Name: characterDetails.Name,
					}
					mockCSC.EXPECT().GetCharacter(gomock.Eq(ctx), gomock.Eq(target)).Return(characterDetails, nil)
					id, err := srvContext.GetCharacterIdFromTarget(ctx, target)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(id).To(Equal(uint(characterDetails.Id)))
				})

				It("should error if no charater is found", func() {
					err := fmt.Errorf("")
					target.Type = &pb.CharacterTarget_Name{
						Name: characterDetails.Name,
					}
					mockCSC.EXPECT().GetCharacter(gomock.Eq(ctx), gomock.Eq(target)).Return(nil, err)
					id, err := srvContext.GetCharacterIdFromTarget(ctx, target)
					Expect(err).To(MatchError(err))
					Expect(id).NotTo(Equal(uint(characterDetails.Id)))
				})
			})
		})
	})

	Context("GetCharacterNameFromTarget", func() {
		Context("with invalid inputs", func() {
			It("should require a target", func() {
				name, err := srvContext.GetCharacterNameFromTarget(ctx, nil)
				Expect(err).Should(HaveOccurred())
				Expect(name).NotTo(Equal(characterDetails.Name))
			})

			It("should error on invalid target", func() {
				name, err := srvContext.GetCharacterNameFromTarget(ctx, target)
				Expect(err).Should(MatchError(common.ErrHandleRequest.Err()))
				Expect(name).NotTo(Equal(characterDetails.Name))
				Expect(hook.AllEntries()).To(HaveLen(1))
				Expect(hook.LastEntry().String()).To(ContainSubstring("target type unknown"))
			})
		})

		Context("with valid inputs", func() {
			It("should return existing name if provided", func() {
				target.Type = &pb.CharacterTarget_Name{
					Name: characterDetails.Name,
				}
				name, err := srvContext.GetCharacterNameFromTarget(ctx, target)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(name).To(Equal(target.GetName()))
			})

			Context("looking up character names", func() {
				It("should if given an id", func() {
					target.Type = &pb.CharacterTarget_Id{
						Id: characterDetails.Id,
					}
					mockCSC.EXPECT().GetCharacter(gomock.Eq(ctx), gomock.Eq(target)).Return(characterDetails, nil)
					name, err := srvContext.GetCharacterNameFromTarget(ctx, target)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(name).To(Equal(characterDetails.Name))
				})

				It("should error if no charater is found", func() {
					err := fmt.Errorf("")
					target.Type = &pb.CharacterTarget_Id{
						Id: characterDetails.Id,
					}
					mockCSC.EXPECT().GetCharacter(gomock.Eq(ctx), gomock.Eq(target)).Return(nil, err)
					name, err := srvContext.GetCharacterNameFromTarget(ctx, target)
					Expect(err).To(MatchError(err))
					Expect(name).NotTo(Equal(characterDetails.Name))
				})
			})
		})
	})
})
