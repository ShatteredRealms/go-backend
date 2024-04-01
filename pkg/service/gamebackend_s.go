package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
)

var (
	gamebackendTracer = otel.Tracer("Inner-GamebackendService")
)

type dimensionService interface {
	CreateDimension(
		ctx context.Context,
		name string,
		location string,
		version string,
		mapIds []*uuid.UUID,
	) (*model.Dimension, error)
	DuplicateDimension(ctx context.Context, target *pb.DimensionTarget, name string) (*model.Dimension, error)
	FindDimensionByName(ctx context.Context, name string) (*model.Dimension, error)
	FindDimensionById(ctx context.Context, id *uuid.UUID) (*model.Dimension, error)
	FindDimension(ctx context.Context, target *pb.DimensionTarget) (*model.Dimension, error)
	FindDimensionsByNames(ctx context.Context, names []string) (model.Dimensions, error)
	FindDimensionsByIds(ctx context.Context, ids []*uuid.UUID) (model.Dimensions, error)
	FindDimensionsWithMapIds(ctx context.Context, ids []*uuid.UUID) (model.Dimensions, error)
	FindAllDimensions(ctx context.Context) (model.Dimensions, error)
	EditDimension(ctx context.Context, request *pb.EditDimensionRequest) (*model.Dimension, error)
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
	) (*model.Map, error)
	FindMapByName(ctx context.Context, name string) (*model.Map, error)
	FindMapById(ctx context.Context, id *uuid.UUID) (*model.Map, error)
	FindMap(ctx context.Context, target *pb.MapTarget) (*model.Map, error)
	FindMapsByNames(ctx context.Context, names []string) (model.Maps, error)
	FindMapsByIds(ctx context.Context, ids []*uuid.UUID) (model.Maps, error)
	FindAllMaps(ctx context.Context) (model.Maps, error)
	EditMap(ctx context.Context, request *pb.EditMapRequest) (*model.Map, error)
	DeleteMapByName(ctx context.Context, name string) error
	DeleteMapById(ctx context.Context, id *uuid.UUID) error
	DeleteMap(ctx context.Context, target *pb.MapTarget) error
}

type connectionService interface {
	CreatePendingConnection(ctx context.Context, character string, serverName string) (*model.PendingConnection, error)
	CheckPlayerConnection(ctx context.Context, id *uuid.UUID, serverName string) (*model.PendingConnection, error)
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
func (s *gamebackendService) FindDimension(ctx context.Context, target *pb.DimensionTarget) (*model.Dimension, error) {
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
		return nil, model.ErrHandleRequest.Err()

	}
}

// FindMap implements GamebackendService.
func (s *gamebackendService) FindMap(ctx context.Context, target *pb.MapTarget) (*model.Map, error) {
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
		return nil, model.ErrHandleRequest.Err()

	}
}

// DuplicateDimension implements GamebackendService.
func (s *gamebackendService) DuplicateDimension(ctx context.Context, target *pb.DimensionTarget, name string) (*model.Dimension, error) {
	dimension, err := s.FindDimension(ctx, target)
	if err != nil {
		return nil, err
	}
	if dimension == nil {
		return nil, model.ErrDoesNotExist.Err()

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
		return nil, errors.Wrap(err, "migrate db")
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
		s.gamebackendRepo.DeletePendingConnection(ctx, id)
		return nil, fmt.Errorf("expired")
	}

	s.gamebackendRepo.DeletePendingConnection(ctx, id)
	return pc, nil
}

// CreateDimension implements GamebackendService.
func (s *gamebackendService) CreateDimension(ctx context.Context, name string, location string, version string, mapIds []*uuid.UUID) (*model.Dimension, error) {
	return s.gamebackendRepo.CreateDimension(ctx, name, location, version, mapIds)
}

// CreateMap implements GamebackendService.
func (s *gamebackendService) CreateMap(ctx context.Context, name string, path string, maxPlayers uint64, instanced bool) (*model.Map, error) {
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
func (s *gamebackendService) EditDimension(ctx context.Context, request *pb.EditDimensionRequest) (*model.Dimension, error) {
	currentDimension, err := s.FindDimension(ctx, request.Target)

	if err != nil {
		return nil, err
	}

	if currentDimension == nil {
		return nil, model.ErrDoesNotExist.Err()

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
			return nil, errors.Wrap(err, "invalid map ids")
		}

		maps, err := s.gamebackendRepo.FindMapsByIds(ctx, ids)
		if err != nil {
			return nil, errors.Wrap(err, "getting maps")
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
func (s *gamebackendService) EditMap(ctx context.Context, request *pb.EditMapRequest) (*model.Map, error) {
	var currentMap *model.Map
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
		err = model.ErrHandleRequest.Err()

	}

	if err != nil {
		return nil, err
	}

	if currentMap == nil {
		return nil, model.ErrDoesNotExist.Err()

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
func (s *gamebackendService) FindAllDimensions(ctx context.Context) (model.Dimensions, error) {
	return s.gamebackendRepo.FindAllDimensions(ctx)
}

// FindAllMaps implements GamebackendService.
func (s *gamebackendService) FindAllMaps(ctx context.Context) (model.Maps, error) {
	return s.gamebackendRepo.FindAllMaps(ctx)
}

// FindDimensionById implements GamebackendService.
func (s *gamebackendService) FindDimensionById(ctx context.Context, id *uuid.UUID) (*model.Dimension, error) {
	return s.gamebackendRepo.FindDimensionById(ctx, id)
}

// FindDimensionByName implements GamebackendService.
func (s *gamebackendService) FindDimensionByName(ctx context.Context, name string) (*model.Dimension, error) {
	return s.gamebackendRepo.FindDimensionByName(ctx, name)
}

// FindDimensionsByIds implements GamebackendService.
func (s *gamebackendService) FindDimensionsByIds(ctx context.Context, ids []*uuid.UUID) (model.Dimensions, error) {
	return s.gamebackendRepo.FindDimensionsByIds(ctx, ids)
}

// FindDimensionsByNames implements GamebackendService.
func (s *gamebackendService) FindDimensionsByNames(ctx context.Context, names []string) (model.Dimensions, error) {
	return s.gamebackendRepo.FindDimensionsByNames(ctx, names)
}

// FindDimensionsWithMapIds implements GamebackendService.
func (s *gamebackendService) FindDimensionsWithMapIds(ctx context.Context, ids []*uuid.UUID) (model.Dimensions, error) {
	return s.gamebackendRepo.FindDimensionsWithMapIds(ctx, ids)
}

// FindMapById implements GamebackendService.
func (s *gamebackendService) FindMapById(ctx context.Context, id *uuid.UUID) (*model.Map, error) {
	return s.gamebackendRepo.FindMapById(ctx, id)
}

// FindMapByName implements GamebackendService.
func (s *gamebackendService) FindMapByName(ctx context.Context, name string) (*model.Map, error) {
	return s.gamebackendRepo.FindMapByName(ctx, name)
}

// FindMapsByIds implements GamebackendService.
func (s *gamebackendService) FindMapsByIds(ctx context.Context, ids []*uuid.UUID) (model.Maps, error) {
	return s.gamebackendRepo.FindMapsByIds(ctx, ids)
}

// FindMapsByNames implements GamebackendService.
func (s *gamebackendService) FindMapsByNames(ctx context.Context, names []string) (model.Maps, error) {
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
			return errors.Errorf("invalid id: %s", t.Id)
		}
		return s.DeleteDimensionById(ctx, &id)

	default:
		log.Logger.WithContext(ctx).Errorf("target type unknown: %+v", target)
		return model.ErrHandleRequest.Err()

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
			return errors.Errorf("invalid id: %s", t.Id)
		}
		return s.DeleteMapById(ctx, &id)

	default:
		log.Logger.WithContext(ctx).Errorf("target type unknown: %+v", target)
		return model.ErrHandleRequest.Err()

	}
}
