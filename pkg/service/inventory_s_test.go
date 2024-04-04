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
	"github.com/ShatteredRealms/go-backend/pkg/model/character"
	"github.com/ShatteredRealms/go-backend/pkg/service"
)

var _ = Describe("Inventory service", func() {
	var (
		hook           *test.Hook
		mockController *gomock.Controller
		mockRepository *mocks.MockInventoryRepository

		invService service.InventoryService

		invItem  = &character.InventoryItem{}
		invItem2 = &character.InventoryItem{}
		invItem3 = &character.InventoryItem{}
		charInv  = &character.Inventory{}
		ctx      context.Context
		fakeErr  error
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
		Expect(faker.FakeData(invItem)).To(Succeed())
		Expect(faker.FakeData(invItem2)).To(Succeed())
		Expect(faker.FakeData(invItem3)).To(Succeed())

		ints, err := faker.RandomInt(1, 1e3, 1)
		Expect(err).To(BeNil())
		charInv.CharacterId = uint(ints[0])
		charInv.Inventory = []*character.InventoryItem{invItem}
		charInv.Bank = []*character.InventoryItem{invItem2}
	})

	Describe("GetInventory", func() {
		It("should work", func() {
			mockRepository.EXPECT().GetInventory(ctx, charInv.CharacterId).Return(charInv, fakeErr)
			out, err := invService.GetInventory(ctx, charInv.CharacterId)
			Expect(err).To(MatchError(fakeErr))
			Expect(out).To(Equal(charInv))
		})
	})

	Describe("UpdateInventory", func() {
		It("should work", func() {
			mockRepository.EXPECT().UpdateInventory(ctx, charInv).Return(fakeErr)
			err := invService.UpdateInventory(ctx, charInv)
			Expect(err).To(MatchError(fakeErr))
		})
	})
})
