package character

import "github.com/ShatteredRealms/go-backend/pkg/pb"

type InventoryItem struct {
	Id       string
	Slot     uint32
	Quantity uint64
}

type InventoryItems []*InventoryItem

type Inventory struct {
	CharacterId uint           `json:"_id" bson:"_id"`
	Inventory   InventoryItems `json:"inventory"`
	Bank        InventoryItems `json:"bank"`
}

func (item *InventoryItem) ToPb() *pb.InventoryItem {
	return &pb.InventoryItem{
		Id:       item.Id,
		Slot:     item.Slot,
		Quantity: item.Quantity,
	}
}

func (items InventoryItems) ToPb() []*pb.InventoryItem {
	out := make([]*pb.InventoryItem, len(items))
	for idx, item := range items {
		out[idx] = item.ToPb()
	}

	return out
}

func (inventory *Inventory) ToPb() *pb.Inventory {
	return &pb.Inventory{
		InventoryItems: inventory.Inventory.ToPb(),
		BankItems:      inventory.Bank.ToPb(),
	}
}

func InventoryItemFromPb(item *pb.InventoryItem) *InventoryItem {
	return &InventoryItem{
		Id:       item.Id,
		Slot:     item.Slot,
		Quantity: item.Quantity,
	}
}

func InventoryItemsFromPb(items []*pb.InventoryItem) InventoryItems {
	out := make(InventoryItems, len(items))
	for idx, item := range items {
		out[idx] = InventoryItemFromPb(item)
	}

	return out
}
