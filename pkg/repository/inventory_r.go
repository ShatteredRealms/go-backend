package repository

import (
	"context"

	"github.com/ShatteredRealms/go-backend/pkg/model/character"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InventoryRepository interface {
	GetInventory(ctx context.Context, characterId uint) (*character.Inventory, error)
	UpdateInventory(ctx context.Context, inventory *character.Inventory) error
}

type inventoryRepository struct {
	db *mongo.Database
}

func NewInventoryRepository(db *mongo.Database) InventoryRepository {
	return &inventoryRepository{
		db: db,
	}
}

// GetInventory implements InventoryRepository.
func (r *inventoryRepository) GetInventory(ctx context.Context, characterId uint) (inventory *character.Inventory, err error) {
	err = r.inventoryCollection().FindOne(ctx, bson.D{{"_id", characterId}}).Decode(&inventory)
	if err != nil {
		return nil, err
	}

	return inventory, nil
}

// UpdateInventory implements InventoryRepository.
func (r *inventoryRepository) UpdateInventory(ctx context.Context, inventory *character.Inventory) error {
	_, err := r.inventoryCollection().ReplaceOne(ctx, bson.D{{"_id", inventory.CharacterId}}, inventory, options.Replace().SetUpsert(true))
	return err
}

func (r *inventoryRepository) inventoryCollection() *mongo.Collection {
	return r.db.Collection("inventories")
}
