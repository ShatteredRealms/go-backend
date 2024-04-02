package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ShatteredRealms/go-backend/pkg/common"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/model/game"
	"github.com/ShatteredRealms/go-backend/pkg/model/gamebackend"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	"github.com/google/uuid"
)

type dimensionService interface {
	CreateDimension(
		ctx context.Context,
		name string,
		location string,
		version string,
		mapIds []*uuid.UUID,
	) (*game.Dimension, error)
	DuplicateDimension(ctx context.Context, target *pb.DimensionTarget, name string) (*game.Dimension, error)
	FindDimensionByName(ctx context.Context, name string) (*game.Dimension, error)
	FindDimensionById(ctx context.Context, id *uuid.UUID) (*game.Dimension, error)
	FindDimension(ctx context.Context, target *pb.DimensionTarget) (*game.Dimension, error)
	FindDimensionsByNames(ctx context.Context, names []string) (game.Dimensions, error)
	FindDimensionsByIds(ctx context.Context, ids []*uuid.UUID) (game.Dimensions, error)
	FindDimensionsWithMapIds(ctx context.Context, ids []*uuid.UUID) (game.Dimensions, error)
	FindAllDimensions(ctx context.Context) (game.Dimensions, error)
	EditDimension(ctx context.Context, request *pb.EditDimensionRequest) (*game.Dimension, error)
	DeleteDimensionByName(ctx context.Context, name string) error
	DeleteDimensionById(ctx context.Context, id *uuid.UUID) error
	DeleteDimension(ctx context.Context, target *pb.DimensionTarget) error
}

type mapService interface {
	CreateMap(
		ctx context.Context,
		name string,
		path string,
		maxPlayers uint64,
		instanced bool,
	) (*game.Map, error)
	FindMapByName(ctx context.Context, name string) (*game.Map, error)
	FindMapById(ctx context.Context, id *uuid.UUID) (*game.Map, error)
	FindMap(ctx context.Context, target *pb.MapTarget) (*game.Map, error)
	FindMapsByNames(ctx context.Context, names []string) (game.Maps, error)
	FindMapsByIds(ctx context.Context, ids []*uuid.UUID) (game.Maps, error)
	FindAllMaps(ctx context.Context) (game.Maps, error)
	EditMap(ctx context.Context, request *pb.EditMapRequest) (*game.Map, error)
	DeleteMapByName(ctx context.Context, name string) error
	DeleteMapById(ctx context.Context, id *uuid.UUID) error
	DeleteMap(ctx context.Context, target *pb.MapTarget) error
}

type connectionService interface {
	CreatePendingConnection(ctx context.Context, character string, serverName string) (*gamebackend.PendingConnection, error)
	CheckPlayerConnection(ctx context.Context, id *uuid.UUID, serverName string) (*gamebackend.PendingConnection, error)
}

type GamebackendService interface {
	connectionService
	dimensionService
	mapService
}

type gamebackendService struct {
	gamebackendRepo repository.GamebackendRepository
}

// FindDimension implements GamebackendService.
func (s *gamebackendService) FindDimension(ctx context.Context, target *pb.DimensionTarget) (*game.Dimension, error) {
	switch t := target.FindBy.(type) {
	case *pb.DimensionTarget_Name:
		return s.gamebackendRepo.FindDimensionByName(ctx, t.Name)

	case *pb.DimensionTarget_Id:
		id, err := uuid.Parse(t.Id)
		if err != nil {
			return nil, fmt.Errorf("invalid id: %s", t.Id)
		}
		return s.gamebackendRepo.FindDimensionById(ctx, &id)

	default:
		log.Logger.WithContext(ctx).Errorf("target type unknown: %+v", target)
		return nil, common.ErrHandleRequest.Err()

	}
}

// FindMap implements GamebackendService.
func (s *gamebackendService) FindMap(ctx context.Context, target *pb.MapTarget) (*game.Map, error) {
	switch t := target.FindBy.(type) {
	case *pb.MapTarget_Name:
		return s.gamebackendRepo.FindMapByName(ctx, t.Name)

	case *pb.MapTarget_Id:
		id, err := uuid.Parse(t.Id)
		if err != nil {
			return nil, fmt.Errorf("invalid id: %s", t.Id)
		}
		return s.gamebackendRepo.FindMapById(ctx, &id)

	default:
		log.Logger.WithContext(ctx).Errorf("target type unknown: %+v", target)
		return nil, common.ErrHandleRequest.Err()

	}
}

// DuplicateDimension implements GamebackendService.
func (s *gamebackendService) DuplicateDimension(ctx context.Context, target *pb.DimensionTarget, name string) (*game.Dimension, error) {
	dimension, err := s.FindDimension(ctx, target)
	if err != nil {
		return nil, err
	}
	if dimension == nil {
		return nil, common.ErrDoesNotExist.Err()

	}

	ids := make([]*uuid.UUID, len(dimension.Maps))
	for i, m := range dimension.Maps {
		ids[i] = m.Id
	}

	return s.gamebackendRepo.CreateDimension(ctx, name, dimension.Location, dimension.Version, ids)
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
) (*gamebackend.PendingConnection, error) {
	return s.gamebackendRepo.CreatePendingConnection(ctx, character, serverName)
}

// DeletePendingConnection implements GamebackendService.
func (s *gamebackendService) CheckPlayerConnection(ctx context.Context, id *uuid.UUID, serverName string) (*gamebackend.PendingConnection, error) {
	pc := s.gamebackendRepo.FindPendingConnection(ctx, id)
	if pc == nil {
		return nil, fmt.Errorf("connection not found")
	}

	if pc.ServerName != serverName {
		log.Logger.WithContext(ctx).Warningf("%s requested: %s, but required: %s", pc.Character, serverName, pc.ServerName)
		return nil, fmt.Errorf("invalid server")
	}

	// @TODO(wil): Make expiration time a configuration variable
	expireTime := pc.CreatedAt.Add(30 * time.Second)
	if expireTime.Unix() < time.Now().Unix() {
		log.Logger.WithContext(ctx).Infof("connection expired for %s", pc.Character)
		err := s.gamebackendRepo.DeletePendingConnection(ctx, id)
		return nil, errors.Join(fmt.Errorf("expired"), err)
	}

	s.gamebackendRepo.DeletePendingConnection(ctx, id)
	return pc, nil
}

// CreateDimension implements GamebackendService.
func (s *gamebackendService) CreateDimension(ctx context.Context, name string, location string, version string, mapIds []*uuid.UUID) (*game.Dimension, error) {
	return s.gamebackendRepo.CreateDimension(ctx, name, location, version, mapIds)
}

// CreateMap implements GamebackendService.
func (s *gamebackendService) CreateMap(ctx context.Context, name string, path string, maxPlayers uint64, instanced bool) (*game.Map, error) {
	return s.gamebackendRepo.CreateMap(ctx, name, path, maxPlayers, instanced)
}

// DeleteDimension implements GamebackendService.
func (s *gamebackendService) DeleteDimensionById(ctx context.Context, id *uuid.UUID) error {
	return s.gamebackendRepo.DeleteDimensionById(ctx, id)
}

// DeleteDimensionByName implements GamebackendService.
func (s *gamebackendService) DeleteDimensionByName(ctx context.Context, name string) error {
	return s.gamebackendRepo.DeleteDimensionByName(ctx, name)
}

// DeleteMapById implements GamebackendService.
func (s *gamebackendService) DeleteMapById(ctx context.Context, id *uuid.UUID) error {
	return s.gamebackendRepo.DeleteMapById(ctx, id)
}

// DeleteMapByName implements GamebackendService.
func (s *gamebackendService) DeleteMapByName(ctx context.Context, name string) error {
	return s.gamebackendRepo.DeleteMapByName(ctx, name)
}

// EditDimension implements GamebackendService.
func (s *gamebackendService) EditDimension(ctx context.Context, request *pb.EditDimensionRequest) (*game.Dimension, error) {
	currentDimension, err := s.FindDimension(ctx, request.Target)

	if err != nil {
		return nil, err
	}

	if currentDimension == nil {
		return nil, common.ErrDoesNotExist.Err()

	}

	if request.OptionalName != nil {
		currentDimension.Name = request.OptionalName.(*pb.EditDimensionRequest_Name).Name
	}

	if request.OptionalVersion != nil {
		currentDimension.Version = request.OptionalVersion.(*pb.EditDimensionRequest_Version).Version
	}

	if request.OptionalLocation != nil {
		currentDimension.Location = request.OptionalLocation.(*pb.EditDimensionRequest_Location).Location
	}

	if request.EditMaps {
		ids, err := helpers.ParseUUIDs(request.MapIds)
		if err != nil {
			return nil, fmt.Errorf("invalid map ids: %w", err)
		}

		maps, err := s.gamebackendRepo.FindMapsByIds(ctx, ids)
		if err != nil {
			return nil, fmt.Errorf("getting maps: %w", err)
		}

		if len(ids) != len(maps) {
			missingIds := make([]*uuid.UUID, len(ids))
			count := 0
			mapsSet := make(map[uuid.UUID]struct{})
			for _, currentMap := range maps {
				mapsSet[*currentMap.Id] = struct{}{}
			}
			for _, id := range ids {
				if _, ok := mapsSet[*id]; !ok {
					missingIds[count] = id
					count++
				}
			}
			return nil, fmt.Errorf("could not find map ids: %v", missingIds)
		}

		currentDimension.Maps = maps
	}

	return s.gamebackendRepo.SaveDimension(ctx, currentDimension)
}

// EditMap implements GamebackendService.
func (s *gamebackendService) EditMap(ctx context.Context, request *pb.EditMapRequest) (*game.Map, error) {
	var currentMap *game.Map
	var err error

	switch target := request.Target.FindBy.(type) {
	case *pb.MapTarget_Id:
		id, err := uuid.Parse(target.Id)
		if err != nil {
			return nil, err
		}
		currentMap, err = s.FindMapById(ctx, &id)
	case *pb.MapTarget_Name:
		currentMap, err = s.FindMapByName(ctx, target.Name)
	default:
		log.Logger.WithContext(ctx).Errorf("map target type unknown: %+v", request.Target)
		err = common.ErrHandleRequest.Err()

	}

	if err != nil {
		return nil, err
	}

	if currentMap == nil {
		return nil, common.ErrDoesNotExist.Err()

	}

	if request.OptionalName != nil {
		currentMap.Name = request.OptionalName.(*pb.EditMapRequest_Name).Name
	}

	if request.OptionalPath != nil {
		currentMap.Path = request.OptionalPath.(*pb.EditMapRequest_Path).Path
	}

	if request.OptionalMaxPlayers != nil {
		currentMap.MaxPlayers = request.OptionalMaxPlayers.(*pb.EditMapRequest_MaxPlayers).MaxPlayers
	}

	if request.OptionalInstanced != nil {
		currentMap.Instanced = request.OptionalInstanced.(*pb.EditMapRequest_Instanced).Instanced
	}

	return s.gamebackendRepo.SaveMap(ctx, currentMap)
}

// FindAllDimensions implements GamebackendService.
func (s *gamebackendService) FindAllDimensions(ctx context.Context) (game.Dimensions, error) {
	return s.gamebackendRepo.FindAllDimensions(ctx)
}

// FindAllMaps implements GamebackendService.
func (s *gamebackendService) FindAllMaps(ctx context.Context) (game.Maps, error) {
	return s.gamebackendRepo.FindAllMaps(ctx)
}

// FindDimensionById implements GamebackendService.
func (s *gamebackendService) FindDimensionById(ctx context.Context, id *uuid.UUID) (*game.Dimension, error) {
	return s.gamebackendRepo.FindDimensionById(ctx, id)
}

// FindDimensionByName implements GamebackendService.
func (s *gamebackendService) FindDimensionByName(ctx context.Context, name string) (*game.Dimension, error) {
	return s.gamebackendRepo.FindDimensionByName(ctx, name)
}

// FindDimensionsByIds implements GamebackendService.
func (s *gamebackendService) FindDimensionsByIds(ctx context.Context, ids []*uuid.UUID) (game.Dimensions, error) {
	return s.gamebackendRepo.FindDimensionsByIds(ctx, ids)
}

// FindDimensionsByNames implements GamebackendService.
func (s *gamebackendService) FindDimensionsByNames(ctx context.Context, names []string) (game.Dimensions, error) {
	return s.gamebackendRepo.FindDimensionsByNames(ctx, names)
}

// FindDimensionsWithMapIds implements GamebackendService.
func (s *gamebackendService) FindDimensionsWithMapIds(ctx context.Context, ids []*uuid.UUID) (game.Dimensions, error) {
	return s.gamebackendRepo.FindDimensionsWithMapIds(ctx, ids)
}

// FindMapById implements GamebackendService.
func (s *gamebackendService) FindMapById(ctx context.Context, id *uuid.UUID) (*game.Map, error) {
	return s.gamebackendRepo.FindMapById(ctx, id)
}

// FindMapByName implements GamebackendService.
func (s *gamebackendService) FindMapByName(ctx context.Context, name string) (*game.Map, error) {
	return s.gamebackendRepo.FindMapByName(ctx, name)
}

// FindMapsByIds implements GamebackendService.
func (s *gamebackendService) FindMapsByIds(ctx context.Context, ids []*uuid.UUID) (game.Maps, error) {
	return s.gamebackendRepo.FindMapsByIds(ctx, ids)
}

// FindMapsByNames implements GamebackendService.
func (s *gamebackendService) FindMapsByNames(ctx context.Context, names []string) (game.Maps, error) {
	return s.gamebackendRepo.FindMapsByNames(ctx, names)
}

// DeleteDimension implements GamebackendService.
func (s *gamebackendService) DeleteDimension(ctx context.Context, target *pb.DimensionTarget) error {
	switch t := target.FindBy.(type) {
	case *pb.DimensionTarget_Name:
		return s.DeleteDimensionByName(ctx, t.Name)

	case *pb.DimensionTarget_Id:
		id, err := uuid.Parse(t.Id)
		if err != nil {
			return fmt.Errorf("invalid id: %s", t.Id)
		}
		return s.DeleteDimensionById(ctx, &id)

	default:
		log.Logger.WithContext(ctx).Errorf("target type unknown: %+v", target)
		return common.ErrHandleRequest.Err()

	}
}

// DeleteMap implements GamebackendService.
func (s *gamebackendService) DeleteMap(ctx context.Context, target *pb.MapTarget) error {
	switch t := target.FindBy.(type) {
	case *pb.MapTarget_Name:
		return s.DeleteMapByName(ctx, t.Name)

	case *pb.MapTarget_Id:
		id, err := uuid.Parse(t.Id)
		if err != nil {
			return fmt.Errorf("invalid id: %s", t.Id)
		}
		return s.DeleteMapById(ctx, &id)

	default:
		log.Logger.WithContext(ctx).Errorf("target type unknown: %+v", target)
		return common.ErrHandleRequest.Err()

	}
}
