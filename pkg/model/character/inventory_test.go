package character_test

import (
	"github.com/bxcodec/faker/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ShatteredRealms/go-backend/pkg/model/character"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
)

var _ = Describe("Inventory model", func() {
	var (
		invItem  = &character.InventoryItem{}
		invItem2 = &character.InventoryItem{}
		invItem3 = &character.InventoryItem{}
		charInv  = &character.Inventory{}
	)

	BeforeEach(func() {
		Expect(faker.FakeData(invItem)).To(Succeed())
		Expect(faker.FakeData(invItem2)).To(Succeed())
		Expect(faker.FakeData(invItem3)).To(Succeed())

		ints, err := faker.RandomInt(1, 1e3, 1)
		Expect(err).To(BeNil())
		charInv.CharacterId = uint(ints[0])
		charInv.Inventory = []*character.InventoryItem{invItem}
		charInv.Bank = []*character.InventoryItem{invItem2}
	})

	validateInventoryItem := (func(invItem *character.InventoryItem, pb *pb.InventoryItem) {
		Expect(pb.Id).To(Equal(invItem.Id))
		Expect(pb.Slot).To(Equal(invItem.Slot))
		Expect(pb.Quantity).To(Equal(invItem.Quantity))
	})

	Describe("ToPb", func() {
		Context("InventoryItem", func() {
			It("should convert single inventory item to protobuf and retain all fields", func() {
				validateInventoryItem(invItem, invItem.ToPb())
				validateInventoryItem(invItem2, invItem2.ToPb())
				validateInventoryItem(invItem3, invItem3.ToPb())
			})

			It("should convert array of inventory items to protobuf and retain all fields", func() {
				var inventoryItems character.InventoryItems
				inventoryItems = make([]*character.InventoryItem, 10)
				for idx := range inventoryItems {
					inventoryItems[idx] = &character.InventoryItem{}
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
			out := charInv.ToPb()
			Expect(out.InventoryItems).To(HaveLen(len(charInv.Inventory)))
			Expect(out.BankItems).To(HaveLen(len(charInv.Bank)))
		})
	})

	Describe("FromPb", func() {
		It("should convert single inventory item to protobuf and retain all fields", func() {
			pbInvItem := &pb.InventoryItem{}
			Expect(faker.FakeData(pbInvItem)).To(Succeed())
			validateInventoryItem(character.InventoryItemFromPb(pbInvItem), pbInvItem)
		})

		It("should convert array of inventory items to protobuf and retain all fields", func() {
			var pbInventoryItems []*pb.InventoryItem
			pbInventoryItems = make([]*pb.InventoryItem, 10)
			for idx := range pbInventoryItems {
				pbInventoryItems[idx] = &pb.InventoryItem{}
				faker.FakeData(pbInventoryItems[idx])
			}
			out := character.InventoryItemsFromPb(pbInventoryItems)
			Expect(out).To(HaveLen(len(pbInventoryItems)))
			for idx := range pbInventoryItems {
				validateInventoryItem(out[idx], pbInventoryItems[idx])
			}
		})
	})
})
