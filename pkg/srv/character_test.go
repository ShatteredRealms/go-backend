package srv_test

import (
	"context"
	"time"

	app "github.com/ShatteredRealms/go-backend/cmd/character/app"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/mocks"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/srv"
	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus/hooks/test"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ = Describe("Character server", func() {
	var (
		hook            *test.Hook
		mockController  *gomock.Controller
		mockCharService *mocks.MockCharacterService
		mockInvService  *mocks.MockInventoryService
		charCtx         *app.CharactersServerContext

		server pb.CharacterServiceServer
		ctx    = context.Background()

		character *model.Character
	)

	BeforeEach(func() {
		log.Logger, hook = test.NewNullLogger()
		mockController = gomock.NewController(GinkgoT())

		mockCharService = mocks.NewMockCharacterService(mockController)
		mockInvService = mocks.NewMockInventoryService(mockController)

		charCtx = &app.CharactersServerContext{
			GlobalConfig:     conf,
			CharacterService: mockCharService,
			InventoryService: mockInvService,
			KeycloakClient:   keycloak,
			Tracer:           otel.Tracer("test-character"),
		}

		var err error
		server, err = srv.NewCharacterServiceServer(ctx, charCtx)
		Expect(err).NotTo(HaveOccurred())
		Expect(server).NotTo(BeNil())

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

	Describe("AddCharacterPlayTime", func() {
		When("given valid input", func() {
			It("should work given character name", func() {
				req := &pb.AddPlayTimeRequest{
					Character: &pb.CharacterTarget{
						Type: &pb.CharacterTarget_Name{Name: character.Name},
					},
					Time: 100,
				}
				mockCharService.EXPECT().FindByName(gomock.Any(), character.Name).Return(character, nil)
				mockCharService.EXPECT().AddPlayTime(gomock.Any(), character.ID, req.Time).Return(character, nil)
				out, err := server.AddCharacterPlayTime(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out.Time).To(BeEquivalentTo(character.PlayTime))
			})
			It("should work given character id", func() {
				req := &pb.AddPlayTimeRequest{
					Character: &pb.CharacterTarget{
						Type: &pb.CharacterTarget_Id{Id: uint64(character.ID)},
					},
					Time: 100,
				}
				mockCharService.EXPECT().AddPlayTime(gomock.Any(), character.ID, req.Time).Return(character, nil)
				out, err := server.AddCharacterPlayTime(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out.Time).To(BeEquivalentTo(character.PlayTime))
			})
		})
		When("given invalid input", func() {
			It("should error if adding playtime fails", func() {
				req := &pb.AddPlayTimeRequest{
					Character: &pb.CharacterTarget{
						Type: &pb.CharacterTarget_Id{Id: uint64(character.ID)},
					},
					Time: 100,
				}
				mockCharService.EXPECT().AddPlayTime(gomock.Any(), character.ID, req.Time).Return(character, fakeErr)
				out, err := server.AddCharacterPlayTime(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
			It("should error if unable to lookup target", func() {
				req := &pb.AddPlayTimeRequest{
					Character: &pb.CharacterTarget{
						Type: &pb.CharacterTarget_Name{Name: character.Name},
					},
					Time: 100,
				}
				mockCharService.EXPECT().FindByName(gomock.Any(), character.Name).Return(character, fakeErr)
				out, err := server.AddCharacterPlayTime(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
			It("should error if does not have correct privledges", func() {
				_ = character
				req := &pb.AddPlayTimeRequest{
					Character: &pb.CharacterTarget{
						Type: &pb.CharacterTarget_Id{Id: uint64(character.ID)},
					},
					Time: 100,
				}
				out, err := server.AddCharacterPlayTime(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
			It("should error if claims are invalid", func() {
				_ = character
				req := &pb.AddPlayTimeRequest{
					Character: &pb.CharacterTarget{
						Type: &pb.CharacterTarget_Id{Id: uint64(character.ID)},
					},
					Time: 100,
				}
				out, err := server.AddCharacterPlayTime(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("CreateCharacter", func() {
		var req *pb.CreateCharacterRequest
		BeforeEach(func() {
			req = &pb.CreateCharacterRequest{
				Owner: &pb.UserTarget{
					Target: &pb.UserTarget_Id{
						Id: *player.ID,
					},
				},
				Name:   faker.Username(),
				Gender: character.Gender,
				Realm:  character.Realm,
			}
		})
		When("given valid input", func() {
			It("should work for players creating for themselves", func() {
				mockCharService.EXPECT().Create(gomock.Any(), *player.ID, req.Name, req.Gender, req.Realm).Return(character, nil)
				out, err := server.CreateCharacter(incPlayerCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())

				req.Owner.Target = &pb.UserTarget_Username{Username: *player.Username}
				mockCharService.EXPECT().Create(gomock.Any(), *player.ID, req.Name, req.Gender, req.Realm).Return(character, nil)
				out, err = server.CreateCharacter(incPlayerCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})

			It("should work for admins creating for themselves", func() {
				req.Owner.Target = &pb.UserTarget_Id{Id: *admin.ID}
				mockCharService.EXPECT().Create(gomock.Any(), *admin.ID, req.Name, req.Gender, req.Realm).Return(character, nil)
				out, err := server.CreateCharacter(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())

				req.Owner.Target = &pb.UserTarget_Username{Username: *admin.Username}
				mockCharService.EXPECT().Create(gomock.Any(), *admin.ID, req.Name, req.Gender, req.Realm).Return(character, nil)
				out, err = server.CreateCharacter(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})

			It("should work for admins creating for others", func() {
				mockCharService.EXPECT().Create(gomock.Any(), *player.ID, req.Name, req.Gender, req.Realm).Return(character, nil)
				out, err := server.CreateCharacter(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())

				req.Owner.Target = &pb.UserTarget_Username{Username: *player.Username}
				mockCharService.EXPECT().Create(gomock.Any(), *player.ID, req.Name, req.Gender, req.Realm).Return(character, nil)
				out, err = server.CreateCharacter(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should error if requesting to create for other without permission", func() {
				req.Owner.Target = &pb.UserTarget_Id{Id: *admin.ID}
				out, err := server.CreateCharacter(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if the request is missing a target", func() {
				req.Owner.Target = nil
				out, err := server.CreateCharacter(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error with invalid claims", func() {
				out, err := server.CreateCharacter(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error with invalid permissions", func() {
				out, err := server.CreateCharacter(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if creating fails", func() {
				mockCharService.EXPECT().Create(gomock.Any(), *player.ID, req.Name, req.Gender, req.Realm).Return(character, fakeErr)
				out, err := server.CreateCharacter(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())

				mockCharService.EXPECT().Create(gomock.Any(), *player.ID, req.Name, req.Gender, req.Realm).Return(nil, nil)
				out, err = server.CreateCharacter(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("Delete Character", func() {
		var req *pb.CharacterTarget
		BeforeEach(func() {
			req = &pb.CharacterTarget{
				Type: &pb.CharacterTarget_Id{
					Id: uint64(character.ID),
				},
			}
		})
		When("Given invalid input", func() {
			It("should work given", func() {
				mockCharService.EXPECT().Delete(gomock.Any(), character.ID).Return(nil)
				mockCharService.EXPECT().FindById(gomock.Any(), character.ID).Return(character, nil)
				out, err := server.DeleteCharacter(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())

				req.Type = &pb.CharacterTarget_Name{
					Name: character.Name,
				}
				mockCharService.EXPECT().Delete(gomock.Any(), character.ID).Return(nil)
				mockCharService.EXPECT().FindByName(gomock.Any(), character.Name).Return(character, nil)
				out, err = server.DeleteCharacter(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})
		})

		When("Given invalid input", func() {
			It("should error if given invalid context", func() {
				out, err := server.DeleteCharacter(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if user has invalid permission to delete", func() {
				out, err := server.DeleteCharacter(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if error finding by name", func() {
				mockCharService.EXPECT().FindById(gomock.Any(), character.ID).Return(nil, nil)
				out, err := server.DeleteCharacter(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if character not found", func() {
				mockCharService.EXPECT().FindById(gomock.Any(), character.ID).Return(nil, nil)
				out, err := server.DeleteCharacter(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if user has invalid permission for deleting other", func() {
				character.OwnerId = *admin.ID
				mockCharService.EXPECT().FindById(gomock.Any(), character.ID).Return(character, nil)
				out, err := server.DeleteCharacter(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("EditCharacter", func() {
		var req *pb.EditCharacterRequest
		BeforeEach(func() {
			req = &pb.EditCharacterRequest{
				Target: &pb.CharacterTarget{
					Type: &pb.CharacterTarget_Id{
						Id: uint64(character.ID),
					},
				},
				OptionalOwnerId:  nil,
				OptionalNewName:  nil,
				OptionalGender:   nil,
				OptionalRealm:    nil,
				OptionalPlayTime: nil,
				OptionalLocation: nil,
			}
		})
		When("Given invalid input (moderator+ permissions)", func() {
			It("should work with id target", func() {
				mockCharService.EXPECT().Edit(gomock.Any(), req).Return(character, nil)
				out, err := server.EditCharacter(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})

			It("should work with name target", func() {
				req.Target.Type = &pb.CharacterTarget_Name{
					Name: character.Name,
				}
				mockCharService.EXPECT().Edit(gomock.Any(), req).Return(character, nil)
				out, err := server.EditCharacter(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})
		})

		When("Given invalid input", func() {
			It("should error if given invalid context", func() {
				out, err := server.EditCharacter(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error when invalid permissions (player)", func() {
				out, err := server.EditCharacter(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error when invalid permissions (guest)", func() {
				out, err := server.EditCharacter(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
			It("should error when editing service fails", func() {
				mockCharService.EXPECT().Edit(gomock.Any(), req).Return(nil, fakeErr)
				out, err := server.EditCharacter(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("GetCharacter", func() {
		var req *pb.CharacterTarget
		BeforeEach(func() {
			req = &pb.CharacterTarget{
				Type: &pb.CharacterTarget_Id{
					Id: uint64(character.ID),
				},
			}
		})
		When("given valid input", func() {
			It("should work getting self by id (admin)", func() {
				character.OwnerId = *admin.ID
				mockCharService.EXPECT().FindById(gomock.Any(), character.ID).Return(character, nil)
				out, err := server.GetCharacter(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})

			It("should work getting self by id (player)", func() {
				character.OwnerId = *player.ID
				mockCharService.EXPECT().FindById(gomock.Any(), character.ID).Return(character, nil)
				out, err := server.GetCharacter(incPlayerCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})

			It("should work getting other by id (admin)", func() {
				character.OwnerId = *player.ID
				mockCharService.EXPECT().FindById(gomock.Any(), character.ID).Return(character, nil)
				out, err := server.GetCharacter(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})

			It("should work getting other by name (admin)", func() {
				req.Type = &pb.CharacterTarget_Name{
					Name: character.Name,
				}
				character.OwnerId = *player.ID
				mockCharService.EXPECT().FindByName(gomock.Any(), character.Name).Return(character, nil)
				out, err := server.GetCharacter(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should error if given invalid context", func() {
				out, err := server.GetCharacter(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error getting other (player)", func() {
				character.OwnerId = *admin.ID
				mockCharService.EXPECT().FindById(gomock.Any(), character.ID).Return(character, nil)
				out, err := server.GetCharacter(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if getting character fails", func() {
				mockCharService.EXPECT().FindById(gomock.Any(), character.ID).Return(nil, fakeErr)
				out, err := server.GetCharacter(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if no character exists", func() {
				mockCharService.EXPECT().FindById(gomock.Any(), character.ID).Return(nil, nil)
				out, err := server.GetCharacter(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("GetAllCharactersForUser", func() {
		var req *pb.UserTarget
		var characters model.Characters
		BeforeEach(func() {
			req = &pb.UserTarget{
				Target: &pb.UserTarget_Id{
					Id: *player.ID,
				},
			}
			characters = []*model.Character{character}
		})

		When("given valid input", func() {
			It("should work getting self by id (admin)", func() {
				req.Target = &pb.UserTarget_Id{
					Id: *admin.ID,
				}
				mockCharService.EXPECT().FindAllByOwner(gomock.Any(), *admin.ID).Return(characters, nil)
				out, err := server.GetAllCharactersForUser(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})

			It("should work getting self by id (player)", func() {
				mockCharService.EXPECT().FindAllByOwner(gomock.Any(), *player.ID).Return(characters, nil)
				out, err := server.GetAllCharactersForUser(incPlayerCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})

			It("should work getting other by id (admin)", func() {
				req.Target = &pb.UserTarget_Id{
					Id: *admin.ID,
				}
				mockCharService.EXPECT().FindAllByOwner(gomock.Any(), *admin.ID).Return(characters, nil)
				out, err := server.GetAllCharactersForUser(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})

			It("should work getting other by name (admin)", func() {
				req.Target = &pb.UserTarget_Username{
					Username: *admin.Username,
				}
				mockCharService.EXPECT().FindAllByOwner(gomock.Any(), *admin.ID).Return(characters, nil)
				out, err := server.GetAllCharactersForUser(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should error if given invalid context", func() {
				out, err := server.GetAllCharactersForUser(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error getting other (player)", func() {
				req.Target = &pb.UserTarget_Id{
					Id: *admin.ID,
				}
				out, err := server.GetAllCharactersForUser(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if getting character fails", func() {
				mockCharService.EXPECT().FindAllByOwner(gomock.Any(), *player.ID).Return(nil, fakeErr)
				out, err := server.GetAllCharactersForUser(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error getting self (guest)", func() {
				out, err := server.GetAllCharactersForUser(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("GetCharacters", func() {
		var req *emptypb.Empty
		var characters model.Characters
		BeforeEach(func() {
			req = &emptypb.Empty{}
			characters = []*model.Character{character}
		})

		When("given valid input", func() {
			It("should work (admin)", func() {
				mockCharService.EXPECT().FindAll(gomock.Any()).Return(characters, nil)
				out, err := server.GetCharacters(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should error if given invalid context", func() {
				out, err := server.GetCharacters(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error (player)", func() {
				out, err := server.GetCharacters(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error (guest)", func() {
				out, err := server.GetCharacters(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if getting character fails", func() {
				mockCharService.EXPECT().FindAll(gomock.Any()).Return(nil, fakeErr)
				out, err := server.GetCharacters(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("GetInventory", func() {
		var req *pb.CharacterTarget
		var inv *model.CharacterInventory
		BeforeEach(func() {
			req = &pb.CharacterTarget{
				Type: &pb.CharacterTarget_Id{
					Id: uint64(character.ID),
				},
			}
			inv = &model.CharacterInventory{}
			Expect(faker.FakeData(inv)).To(Succeed())
		})

		When("given valid input", func() {
			It("should work if no inventory exists yet (admin)", func() {
				character.OwnerId = *admin.ID
				mockCharService.EXPECT().FindById(gomock.Any(), character.ID).Return(character, nil)
				mockInvService.EXPECT().GetInventory(gomock.Any(), character.ID).Return(nil, mongo.ErrNoDocuments)
				out, err := server.GetInventory(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})

			It("should work getting self by id (admin)", func() {
				character.OwnerId = *admin.ID
				mockCharService.EXPECT().FindById(gomock.Any(), character.ID).Return(character, nil)
				mockInvService.EXPECT().GetInventory(gomock.Any(), character.ID).Return(inv, nil)
				out, err := server.GetInventory(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})

			It("should work getting other by name (admin)", func() {
				req.Type = &pb.CharacterTarget_Name{
					Name: character.Name,
				}
				character.OwnerId = *player.ID
				mockCharService.EXPECT().FindByName(gomock.Any(), character.Name).Return(character, nil)
				mockInvService.EXPECT().GetInventory(gomock.Any(), character.ID).Return(inv, nil)
				out, err := server.GetInventory(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should error if given invalid context", func() {
				out, err := server.GetInventory(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error (player)", func() {
				out, err := server.GetInventory(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error (guest)", func() {
				out, err := server.GetInventory(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if getting character fails", func() {
				mockCharService.EXPECT().FindById(gomock.Any(), character.ID).Return(nil, fakeErr)
				out, err := server.GetInventory(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if no character exists", func() {
				mockCharService.EXPECT().FindById(gomock.Any(), character.ID).Return(nil, nil)
				out, err := server.GetInventory(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})

	Describe("SetInventory", func() {
		var req *pb.UpdateInventoryRequest
		var inv *model.CharacterInventory
		BeforeEach(func() {
			req = &pb.UpdateInventoryRequest{
				Target: &pb.CharacterTarget{
					Type: &pb.CharacterTarget_Id{
						Id: uint64(character.ID),
					},
				},
				InventoryItems: []*pb.InventoryItem{
					{
						Id:       faker.Username(),
						Slot:     1,
						Quantity: 2,
					},
					{
						Id:       faker.Username(),
						Slot:     3,
						Quantity: 4,
					},
				},
				BankItems: []*pb.InventoryItem{
					{
						Id:       faker.Username(),
						Slot:     1,
						Quantity: 2,
					},
					{
						Id:       faker.Username(),
						Slot:     3,
						Quantity: 4,
					},
				},
			}
			inv = &model.CharacterInventory{}
			Expect(faker.FakeData(inv)).To(Succeed())
		})

		When("given valid input", func() {
			It("should work if no inventory exists yet (admin)", func() {
				mockCharService.EXPECT().FindById(gomock.Any(), character.ID).Return(character, nil)
				mockInvService.EXPECT().UpdateInventory(gomock.Any(), gomock.Any()).Return(mongo.ErrNoDocuments)
				out, err := server.SetInventory(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})

			It("should work getting self by id (admin)", func() {
				mockCharService.EXPECT().FindById(gomock.Any(), character.ID).Return(character, nil)
				mockInvService.EXPECT().UpdateInventory(gomock.Any(), gomock.Any()).Return(nil)
				out, err := server.SetInventory(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})

			It("should work getting other by name (admin)", func() {
				req.Target.Type = &pb.CharacterTarget_Name{
					Name: character.Name,
				}
				mockCharService.EXPECT().FindByName(gomock.Any(), character.Name).Return(character, nil)
				mockInvService.EXPECT().UpdateInventory(gomock.Any(), gomock.Any()).Return(nil)
				out, err := server.SetInventory(incAdminCtx, req)
				Expect(err).NotTo(HaveOccurred())
				Expect(out).NotTo(BeNil())
			})
		})

		When("given invalid input", func() {
			It("should error if given invalid context", func() {
				out, err := server.SetInventory(context.Background(), req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error (player)", func() {
				out, err := server.SetInventory(incPlayerCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error (guest)", func() {
				out, err := server.SetInventory(incGuestCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if getting character fails", func() {
				mockCharService.EXPECT().FindById(gomock.Any(), character.ID).Return(nil, fakeErr)
				out, err := server.SetInventory(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if no character exists", func() {
				mockCharService.EXPECT().FindById(gomock.Any(), character.ID).Return(nil, nil)
				out, err := server.SetInventory(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})

			It("should error if update inventory fails", func() {
				mockCharService.EXPECT().FindById(gomock.Any(), character.ID).Return(character, nil)
				mockInvService.EXPECT().UpdateInventory(gomock.Any(), gomock.Any()).Return(fakeErr)
				out, err := server.SetInventory(incAdminCtx, req)
				Expect(err).To(HaveOccurred())
				Expect(out).To(BeNil())
			})
		})
	})
})
