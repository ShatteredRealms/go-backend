package model_test

import (
	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
)

var _ = Describe("Inventory model", func() {
	var (
		inventoryItem      = &model.InventoryItem{}
		inventoryItem2     = &model.InventoryItem{}
		inventoryItem3     = &model.InventoryItem{}
		characterInventory = &model.CharacterInventory{}
	)

	BeforeEach(func() {
		Expect(faker.FakeData(inventoryItem)).To(Succeed())
		Expect(faker.FakeData(inventoryItem2)).To(Succeed())
		Expect(faker.FakeData(inventoryItem3)).To(Succeed())

		ints, err := faker.RandomInt(1, 1e3, 1)
		Expect(err).To(BeNil())
		characterInventory.CharacterId = uint(ints[0])
		characterInventory.Inventory = []*model.InventoryItem{inventoryItem}
		characterInventory.Bank = []*model.InventoryItem{inventoryItem2}
	})

	validateInventoryItem := (func(invItem *model.InventoryItem, pb *pb.InventoryItem) {
		Expect(pb.Id).To(Equal(invItem.Id))
		Expect(pb.Slot).To(Equal(invItem.Slot))
		Expect(pb.Quantity).To(Equal(invItem.Quantity))
	})

	Describe("ToPb", func() {
		Context("InventoryItem", func() {
			It("should convert single inventory item to protobuf and retain all fields", func() {
				validateInventoryItem(inventoryItem, inventoryItem.ToPb())
				validateInventoryItem(inventoryItem2, inventoryItem2.ToPb())
				validateInventoryItem(inventoryItem3, inventoryItem3.ToPb())
			})

			It("should convert array of inventory items to protobuf and retain all fields", func() {
				var inventoryItems model.InventoryItems
				inventoryItems = make([]*model.InventoryItem, 10)
				for idx := range inventoryItems {
					inventoryItems[idx] = &model.InventoryItem{}
					Expect(faker.FakeData(inventoryItems[idx])).To(Succeed())
				}
				out := inventoryItems.ToPb()
				Expect(out).To(HaveLen(len(inventoryItems)))
				for idx := range inventoryItems {
					validateInventoryItem(inventoryItems[idx], out[idx])
				}
			})
		})
		It("should convert character inventory to protobuf and retain all fields", func() {
			out := characterInventory.ToPb()
			Expect(out.InventoryItems).To(HaveLen(len(characterInventory.Inventory)))
			Expect(out.BankItems).To(HaveLen(len(characterInventory.Bank)))
		})
	})

	Describe("FromPb", func() {
		It("should convert single inventory item to protobuf and retain all fields", func() {
			pbInvItem := &pb.InventoryItem{}
			Expect(faker.FakeData(pbInvItem)).To(Succeed())
			validateInventoryItem(model.InventoryItemFromPb(pbInvItem), pbInvItem)
		})

		It("should convert array of inventory items to protobuf and retain all fields", func() {
			var pbInventoryItems []*pb.InventoryItem
			pbInventoryItems = make([]*pb.InventoryItem, 10)
			for idx := range pbInventoryItems {
				pbInventoryItems[idx] = &pb.InventoryItem{}
				faker.FakeData(pbInventoryItems[idx])
			}
			out := model.InventoryItemsFromPb(pbInventoryItems)
			Expect(out).To(HaveLen(len(pbInventoryItems)))
			for idx := range pbInventoryItems {
				validateInventoryItem(out[idx], pbInventoryItems[idx])
			}
		})
	})
})
