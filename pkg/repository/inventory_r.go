package repository

import (
	"context"

	"github.com/ShatteredRealms/go-backend/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type InventoryRepository interface {
	GetInventory(ctx context.Context, characterId uint) (*model.CharacterInventory, error)
	UpdateInventory(ctx context.Context, inventory *model.CharacterInventory) error

	Migrate(ctx context.Context) error
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
func (r *inventoryRepository) GetInventory(ctx context.Context, characterId uint) (inventory *model.CharacterInventory, err error) {
	err = r.inventoryCollection().FindOne(ctx, bson.D{{"characterId", characterId}}).Decode(&inventory)
	if err != nil {
		return nil, err
	}

	return inventory, nil
}

// UpdateInventory implements InventoryRepository.
func (r *inventoryRepository) UpdateInventory(ctx context.Context, inventory *model.CharacterInventory) error {
	_, err := r.inventoryCollection().InsertOne(ctx, inventory)
	return err
}

func (r *inventoryRepository) inventoryCollection() *mongo.Collection {
	return r.db.Collection("inventories")
}

// Migrate implements InventoryRepository.
func (*inventoryRepository) Migrate(ctx context.Context) error {
	return nil
}
