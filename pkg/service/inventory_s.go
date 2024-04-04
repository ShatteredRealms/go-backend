package service

import (
	"context"

	"github.com/ShatteredRealms/go-backend/pkg/model/character"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
)

type InventoryService interface {
	GetInventory(ctx context.Context, characterId uint) (*character.Inventory, error)
	UpdateInventory(ctx context.Context, inventory *character.Inventory) error
}

type inventoryService struct {
	repo repository.InventoryRepository
}

func NewInventoryService(repo repository.InventoryRepository) InventoryService {
	return &inventoryService{
		repo: repo,
	}
}

// GetInventory implements InventoryService.
func (s *inventoryService) GetInventory(ctx context.Context, characterId uint) (*character.Inventory, error) {
	return s.repo.GetInventory(ctx, characterId)
}

// UpdateInventory implements InventoryService.
func (s *inventoryService) UpdateInventory(ctx context.Context, inventory *character.Inventory) error {
	return s.repo.UpdateInventory(ctx, inventory)
}
