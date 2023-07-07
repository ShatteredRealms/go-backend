package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

var (
	gamebackendTracer = otel.Tracer("Inner-GamebackendService")
)

type gamebackendService struct {
	gamebackendRepo repository.GamebackendRepository
}

type GamebackendService interface {
	CreatePendingConnection(ctx context.Context, character string, serverName string) (*model.PendingConnection, error)
	CheckPlayerConnection(ctx context.Context, id *uuid.UUID, serverName string) (*model.PendingConnection, error)
}

func NewGamebackendService(
	ctx context.Context,
	r repository.GamebackendRepository,
) (GamebackendService, error) {
	err := r.Migrate(ctx)
	if err != nil {
		return nil, fmt.Errorf("migrate db: %w", err)
	}

	return &gamebackendService{
		gamebackendRepo: r,
	}, nil
}

// CreatePendingConnection implements GamebackendService.
func (s *gamebackendService) CreatePendingConnection(
	ctx context.Context,
	character string,
	serverName string,
) (*model.PendingConnection, error) {
	return s.gamebackendRepo.CreatePendingConnection(ctx, character, serverName)
}

// DeletePendingConnection implements GamebackendService.
func (s *gamebackendService) CheckPlayerConnection(ctx context.Context, id *uuid.UUID, serverName string) (*model.PendingConnection, error) {
	pc := s.gamebackendRepo.FindPendingConnection(ctx, id)
	if pc == nil {
		return nil, fmt.Errorf("invalid id")
	}

	if pc.ServerName != serverName {
		logrus.WithContext(ctx).Warningf("%s requested: %s, but required: %s", pc.Character, serverName, pc.ServerName)
		return nil, fmt.Errorf("invalid server")
	}

	// @TODO(wil): Make expiration time a configuration variable
	expireTime := pc.CreatedAt.Add(30 * time.Second)
	if expireTime.Unix() < time.Now().Unix() {
		logrus.WithContext(ctx).Infof("connection expired for %s", pc.Character)
		s.gamebackendRepo.DeletePendingConnection(ctx, id)
		return nil, fmt.Errorf("expired")
	}

	s.gamebackendRepo.DeletePendingConnection(ctx, id)
	return pc, nil
}
