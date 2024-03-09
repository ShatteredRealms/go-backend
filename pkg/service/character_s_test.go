package service_test

import (
	"context"
	"fmt"
	"time"

	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus/hooks/test"
	"go.uber.org/mock/gomock"

	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/mocks"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/service"
)

var _ = Describe("Character service", func() {

	var (
		hook *test.Hook

		mockController *gomock.Controller
		mockRepository *mocks.MockCharacterRepository
		charService    service.CharacterService
		ctx            context.Context
		character      *model.Character

		fakeError = fmt.Errorf("error")
	)

	BeforeEach(func() {
		log.Logger, hook = test.NewNullLogger()

		var err error
		ctx = context.Background()
		mockController = gomock.NewController(GinkgoT())
		mockRepository = mocks.NewMockCharacterRepository(mockController)
		mockRepository.EXPECT().Migrate(ctx).Return(nil)
		charService, err = service.NewCharacterService(ctx, mockRepository)
		Expect(err).NotTo(HaveOccurred())
		Expect(charService).NotTo(BeNil())

		character = &model.Character{
			ID:        0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: 0,
			OwnerId:   faker.Username(),
			Name:      "unreal",
			Gender:    "Male",
			Realm:     "Human",
			PlayTime:  100,
			Location: model.Location{
				World: faker.Username(),
				X:     1.1,
				Y:     1.2,
				Z:     1.3,
				Roll:  1.4,
				Pitch: 1.5,
				Yaw:   1.6,
			},
		}

		hook.Reset()
	})

	Describe("NewCharacterService", func() {
		When("given invalid input", func() {
			It("should fail due to migration fail", func() {
				mockRepository.EXPECT().Migrate(ctx).Return(fakeError)
				s, err := service.NewCharacterService(ctx, mockRepository)
				Expect(err).To(MatchError(fakeError))
				Expect(s).To(BeNil())
			})
		})
	})

	Describe("FindByName", func() {
		It("should call repo directly", func() {
			mockRepository.EXPECT().FindByName(ctx, character.Name).Return(character, fakeError)
			out, err := charService.FindByName(ctx, character.Name)
			Expect(err).To(MatchError(fakeError))
			Expect(out).To(Equal(character))
		})
	})

	Describe("Create", func() {
		When("given valid input", func() {
			It("should succeed", func() {
				mockRepository.EXPECT().Create(ctx, gomock.Any()).Return(character, fakeError)
				out, err := charService.Create(ctx, character.OwnerId, character.Name, character.Gender, character.Realm)
				Expect(err).To(MatchError(fakeError))
				Expect(out).To(Equal(character))
			})
		})

		When("given invalid input", func() {
			It("should fail on invalid character", func() {
				out, err := charService.Create(ctx, character.OwnerId, character.Name, "", character.Realm)
				Expect(err).To(MatchError(model.ErrInvalidGender))
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("FindById", func() {
		It("should call repo directly", func() {
			mockRepository.EXPECT().FindById(ctx, character.ID).Return(character, fakeError)
			out, err := charService.FindById(ctx, character.ID)
			Expect(err).To(MatchError(fakeError))
			Expect(out).To(Equal(character))
		})
	})

	Describe("Edit", func() {
		When("given valid input", func() {
			It("should be able to edit by name", func() {
				editReq := &pb.EditCharacterRequest{
					Target: &pb.CharacterTarget{
						Type: &pb.CharacterTarget_Name{
							Name: character.Name,
						},
					},
					OptionalOwnerId: &pb.EditCharacterRequest_OwnerId{
						OwnerId: faker.Username(),
					},
					OptionalNewName: &pb.EditCharacterRequest_NewName{
						NewName: faker.Username(),
					},
					OptionalGender: &pb.EditCharacterRequest_Gender{
						Gender: "Female",
					},
					OptionalRealm: &pb.EditCharacterRequest_Realm{
						Realm: "Cyborg",
					},
					OptionalPlayTime: &pb.EditCharacterRequest_PlayTime{
						PlayTime: 6,
					},
					OptionalLocation: &pb.EditCharacterRequest_Location{
						Location: &pb.Location{
							World: faker.Username(),
							X:     5.0,
							Y:     4.1,
							Z:     3.2,
							Roll:  2.3,
							Pitch: 1.4,
							Yaw:   0.5,
						},
					},
				}
				expectCharacter := new(model.Character)
				*expectCharacter = *character
				expectCharacter.OwnerId = editReq.GetOwnerId()
				expectCharacter.Name = editReq.GetNewName()
				expectCharacter.Gender = editReq.GetGender()
				expectCharacter.Realm = editReq.GetRealm()
				expectCharacter.PlayTime = editReq.GetPlayTime()
				expectCharacter.Location = *model.LocationFromPb(editReq.GetLocation())
				mockRepository.EXPECT().FindByName(ctx, editReq.Target.GetName()).Return(character, nil)
				mockRepository.EXPECT().Save(ctx, expectCharacter).Return(character, fakeError)

				out, err := charService.Edit(ctx, editReq)
				Expect(err).To(MatchError(fakeError))
				Expect(out).To(Equal(character))
			})

			It("should be able to edit by di", func() {
				editReq := &pb.EditCharacterRequest{
					Target: &pb.CharacterTarget{
						Type: &pb.CharacterTarget_Id{
							Id: uint64(character.ID),
						},
					},
					OptionalOwnerId: &pb.EditCharacterRequest_OwnerId{
						OwnerId: faker.Username(),
					},
					OptionalNewName: &pb.EditCharacterRequest_NewName{
						NewName: faker.Username(),
					},
					OptionalGender: &pb.EditCharacterRequest_Gender{
						Gender: "Female",
					},
					OptionalRealm: &pb.EditCharacterRequest_Realm{
						Realm: "Cyborg",
					},
					OptionalPlayTime: &pb.EditCharacterRequest_PlayTime{
						PlayTime: 6,
					},
					OptionalLocation: &pb.EditCharacterRequest_Location{
						Location: &pb.Location{
							World: faker.Username(),
							X:     5.0,
							Y:     4.1,
							Z:     3.2,
							Roll:  2.3,
							Pitch: 1.4,
							Yaw:   0.5,
						},
					},
				}
				expectCharacter := new(model.Character)
				*expectCharacter = *character
				expectCharacter.OwnerId = editReq.GetOwnerId()
				expectCharacter.Name = editReq.GetNewName()
				expectCharacter.Gender = editReq.GetGender()
				expectCharacter.Realm = editReq.GetRealm()
				expectCharacter.PlayTime = editReq.GetPlayTime()
				expectCharacter.Location = *model.LocationFromPb(editReq.GetLocation())
				mockRepository.EXPECT().FindById(ctx, uint(editReq.Target.GetId())).Return(character, nil)
				mockRepository.EXPECT().Save(ctx, expectCharacter).Return(character, fakeError)

				out, err := charService.Edit(ctx, editReq)
				Expect(err).To(MatchError(fakeError))
				Expect(out).To(Equal(character))
			})
		})

		When("given invalid input", func() {
			It("should fail on invalid new character details", func() {
				editReq := &pb.EditCharacterRequest{
					Target: &pb.CharacterTarget{
						Type: &pb.CharacterTarget_Name{
							Name: character.Name,
						},
					},
					OptionalOwnerId: &pb.EditCharacterRequest_OwnerId{
						OwnerId: faker.Username(),
					},
					OptionalNewName: &pb.EditCharacterRequest_NewName{
						NewName: faker.Username(),
					},
					OptionalGender: &pb.EditCharacterRequest_Gender{
						Gender: faker.Username(),
					},
					OptionalRealm: &pb.EditCharacterRequest_Realm{
						Realm: "Cyborg",
					},
					OptionalPlayTime: &pb.EditCharacterRequest_PlayTime{
						PlayTime: 6,
					},
					OptionalLocation: &pb.EditCharacterRequest_Location{
						Location: &pb.Location{
							World: faker.Username(),
							X:     5.0,
							Y:     4.1,
							Z:     3.2,
							Roll:  2.3,
							Pitch: 1.4,
							Yaw:   0.5,
						},
					},
				}
				expectCharacter := new(model.Character)
				*expectCharacter = *character
				expectCharacter.OwnerId = editReq.GetOwnerId()
				expectCharacter.Name = editReq.GetNewName()
				expectCharacter.Gender = editReq.GetGender()
				expectCharacter.Realm = editReq.GetRealm()
				expectCharacter.PlayTime = editReq.GetPlayTime()
				expectCharacter.Location = *model.LocationFromPb(editReq.GetLocation())
				mockRepository.EXPECT().FindByName(ctx, editReq.Target.GetName()).Return(character, nil)
				out, err := charService.Edit(ctx, editReq)
				Expect(err).To(MatchError(model.ErrInvalidGender))
				Expect(out).To(BeNil())
			})

			It("should fail if finding character fails", func() {
				editReq := &pb.EditCharacterRequest{
					Target: &pb.CharacterTarget{
						Type: &pb.CharacterTarget_Id{
							Id: uint64(character.ID),
						},
					},
					OptionalOwnerId: &pb.EditCharacterRequest_OwnerId{
						OwnerId: faker.Username(),
					},
					OptionalNewName: &pb.EditCharacterRequest_NewName{
						NewName: faker.Username(),
					},
					OptionalGender: &pb.EditCharacterRequest_Gender{
						Gender: "Female",
					},
					OptionalRealm: &pb.EditCharacterRequest_Realm{
						Realm: "Cyborg",
					},
					OptionalPlayTime: &pb.EditCharacterRequest_PlayTime{
						PlayTime: 6,
					},
					OptionalLocation: &pb.EditCharacterRequest_Location{
						Location: &pb.Location{
							World: faker.Username(),
							X:     5.0,
							Y:     4.1,
							Z:     3.2,
							Roll:  2.3,
							Pitch: 1.4,
							Yaw:   0.5,
						},
					},
				}
				expectCharacter := new(model.Character)
				*expectCharacter = *character
				expectCharacter.OwnerId = editReq.GetOwnerId()
				expectCharacter.Name = editReq.GetNewName()
				expectCharacter.Gender = editReq.GetGender()
				expectCharacter.Realm = editReq.GetRealm()
				expectCharacter.PlayTime = editReq.GetPlayTime()
				expectCharacter.Location = *model.LocationFromPb(editReq.GetLocation())
				mockRepository.EXPECT().FindById(ctx, uint(editReq.Target.GetId())).Return(character, fakeError)
				out, err := charService.Edit(ctx, editReq)
				Expect(err).To(MatchError(fakeError))
				Expect(out).To(BeNil())
			})

			It("should fail if unknown target", func() {
				editReq := &pb.EditCharacterRequest{
					Target: &pb.CharacterTarget{},
				}
				out, err := charService.Edit(ctx, editReq)
				Expect(err).To(MatchError(model.ErrHandleRequest))
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("Delete", func() {
		When("given valid input", func() {
			It("should try to delete", func() {
				mockRepository.EXPECT().FindById(ctx, character.ID).Return(character, nil)
				mockRepository.EXPECT().Delete(ctx, character)
				err := charService.Delete(ctx, character.ID)
				Expect(err).To(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should error on find error", func() {
				mockRepository.EXPECT().FindById(ctx, character.ID).Return(nil, fakeError)
				err := charService.Delete(ctx, character.ID)
				Expect(err).To(HaveOccurred())
			})

			It("should error on no character found", func() {
				mockRepository.EXPECT().FindById(ctx, character.ID).Return(nil, nil)
				err := charService.Delete(ctx, character.ID)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("FindAll", func() {
		It("should call repo directly", func() {
			characters := []*model.Character{character}
			mockRepository.EXPECT().FindAll(ctx).Return(characters, fakeError)
			out, err := charService.FindAll(ctx)
			Expect(err).To(MatchError(fakeError))
			Expect(out).To(ContainElements(characters))
		})
	})

	Describe("FindAllByOwner", func() {
		It("should call repo directly", func() {
			characters := []*model.Character{character}
			mockRepository.EXPECT().FindAllByOwner(ctx, character.OwnerId).Return(characters, fakeError)
			out, err := charService.FindAllByOwner(ctx, character.OwnerId)
			Expect(err).To(MatchError(fakeError))
			Expect(out).To(ContainElements(characters))
		})
	})

	Describe("AddPlayTime", func() {
		var amount uint64

		BeforeEach(func() {
			nums, err := faker.RandomInt(1, 1e5, 1)
			Expect(err).NotTo(HaveOccurred())
			amount = uint64(nums[0])
		})

		When("given valid input", func() {
			It("should try to update playtime", func() {
				mockRepository.EXPECT().FindById(ctx, character.ID).Return(character, nil)
				charOut := new(model.Character)
				*charOut = *character
				charOut.PlayTime += amount
				mockRepository.EXPECT().Save(ctx, gomock.Any()).Return(charOut, fakeError)
				out, err := charService.AddPlayTime(ctx, character.ID, amount)
				Expect(err).To(MatchError(fakeError))
				Expect(out).To(BeEquivalentTo(charOut.PlayTime))
			})
		})

		When("given invalid input", func() {
			It("should error on find error", func() {
				mockRepository.EXPECT().FindById(ctx, character.ID).Return(nil, fakeError)
				out, err := charService.AddPlayTime(ctx, character.ID, amount)
				Expect(err).To(MatchError(fakeError))
				Expect(out).To(BeEquivalentTo(0))
			})
		})
	})

})
