package service_test

import (
	"context"
	"fmt"

	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus/hooks/test"
	"go.uber.org/mock/gomock"

	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/mocks"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/service"
)

var _ = Describe("Inventory service", func() {
	var (
		hook           *test.Hook
		mockController *gomock.Controller
		mockRepository *mocks.MockInventoryRepository

		invService service.InventoryService

		inventoryItem      = &model.InventoryItem{}
		inventoryItem2     = &model.InventoryItem{}
		inventoryItem3     = &model.InventoryItem{}
		characterInventory = &model.CharacterInventory{}
		ctx                context.Context
		fakeErr            error
	)

	BeforeEach(func() {
		log.Logger, hook = test.NewNullLogger()
		mockController = gomock.NewController(GinkgoT())
		mockRepository = mocks.NewMockInventoryRepository(mockController)
		hook.Reset()

		var err error
		invService = service.NewInventoryService(mockRepository)
		Expect(invService).NotTo(BeNil())
		hook.Reset()

		ctx = context.Background()
		fakeErr = fmt.Errorf("error: %s", faker.Username())
		Expect(faker.FakeData(inventoryItem)).To(Succeed())
		Expect(faker.FakeData(inventoryItem2)).To(Succeed())
		Expect(faker.FakeData(inventoryItem3)).To(Succeed())

		ints, err := faker.RandomInt(1, 1e3, 1)
		Expect(err).To(BeNil())
		characterInventory.CharacterId = uint(ints[0])
		characterInventory.Inventory = []*model.InventoryItem{inventoryItem}
		characterInventory.Bank = []*model.InventoryItem{inventoryItem2}
	})

	Describe("GetInventory", func() {
		It("should work", func() {
			mockRepository.EXPECT().GetInventory(ctx, characterInventory.CharacterId).Return(characterInventory, fakeErr)
			out, err := invService.GetInventory(ctx, characterInventory.CharacterId)
			Expect(err).To(MatchError(fakeErr))
			Expect(out).To(Equal(characterInventory))
		})
	})

	Describe("UpdateInventory", func() {
		It("should work", func() {
			mockRepository.EXPECT().UpdateInventory(ctx, characterInventory).Return(fakeErr)
			err := invService.UpdateInventory(ctx, characterInventory)
			Expect(err).To(MatchError(fakeErr))
		})
	})
})
