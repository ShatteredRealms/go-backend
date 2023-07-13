package service

import (
	"context"

	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
)

type InventoryService interface {
	GetInventory(ctx context.Context, characterId uint) (*model.CharacterInventory, error)
	UpdateInventory(ctx context.Context, inventory *model.CharacterInventory) error
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
func (s *inventoryService) GetInventory(ctx context.Context, characterId uint) (*model.CharacterInventory, error) {
	return s.repo.GetInventory(ctx, characterId)
}

// UpdateInventory implements InventoryService.
func (s *inventoryService) UpdateInventory(ctx context.Context, inventory *model.CharacterInventory) error {
	return s.repo.UpdateInventory(ctx, inventory)
}
