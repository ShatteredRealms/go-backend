package service

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
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
		chatTemplateIds []*uuid.UUID,
	) (*model.Dimension, error)
	DuplicateDimension(ctx context.Context, refId *uuid.UUID, name string) (*model.Dimension, error)
	FindDimensionByName(ctx context.Context, name string) (*model.Dimension, error)
	FindDimensionById(ctx context.Context, id *uuid.UUID) (*model.Dimension, error)
	FindDimensionsByNames(ctx context.Context, names []string) (model.Dimensions, error)
	FindDimensionsByIds(ctx context.Context, ids []*uuid.UUID) (model.Dimensions, error)
	FindAllDimensions(ctx context.Context) (model.Dimensions, error)
	EditDimension(ctx context.Context, request *pb.EditDimensionRequest) (*model.Dimension, error)
	DeleteDimensionByName(ctx context.Context, name string) error
	DeleteDimensionById(ctx context.Context, id *uuid.UUID) error
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
	FindMapsByNames(ctx context.Context, names []string) (model.Maps, error)
	FindMapsByIds(ctx context.Context, ids []*uuid.UUID) (model.Maps, error)
	FindAllMaps(ctx context.Context) (model.Maps, error)
	EditMap(ctx context.Context, request *pb.EditMapRequest) (*model.Map, error)
	DeleteMapByName(ctx context.Context, name string) error
	DeleteMapById(ctx context.Context, id *uuid.UUID) error
}

type chatTemplateService interface {
	CreateChatTemplate(
		ctx context.Context,
		name string,
	) (*model.ChatTemplate, error)
	FindChatTemplateByName(ctx context.Context, name string) (*model.ChatTemplate, error)
	FindChatTemplateById(ctx context.Context, id *uuid.UUID) (*model.ChatTemplate, error)
	FindChatTemplatesByNames(ctx context.Context, names []string) (model.ChatTemplates, error)
	FindChatTemplatesByIds(ctx context.Context, ids []*uuid.UUID) (model.ChatTemplates, error)
	FindAllChatTemplates(ctx context.Context) (model.ChatTemplates, error)
	EditChatTemplate(ctx context.Context, request *pb.EditChatTemplateRequest) (*model.ChatTemplate, error)
	DeleteChatTemplateByName(ctx context.Context, name string) error
	DeleteChatTemplateById(ctx context.Context, id *uuid.UUID) error
}

type connectionService interface {
	CreatePendingConnection(ctx context.Context, character string, serverName string) (*model.PendingConnection, error)
	CheckPlayerConnection(ctx context.Context, id *uuid.UUID, serverName string) (*model.PendingConnection, error)
}

type GamebackendService interface {
	connectionService
	dimensionService
	mapService
	chatTemplateService
}

type gamebackendService struct {
	gamebackendRepo repository.GamebackendRepository
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
		return nil, fmt.Errorf("connection not found")
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

// CreateChatTemplate implements GamebackendService.
func (s *gamebackendService) CreateChatTemplate(ctx context.Context, name string) (*model.ChatTemplate, error) {
	return s.gamebackendRepo.CreateChatTemplate(ctx, name)
}

// CreateDimension implements GamebackendService.
func (s *gamebackendService) CreateDimension(ctx context.Context, name string, location string, version string, mapIds []*uuid.UUID, chatTemplateIds []*uuid.UUID) (*model.Dimension, error) {
	return s.gamebackendRepo.CreateDimension(ctx, name, location, version, mapIds, chatTemplateIds)
}

// CreateMap implements GamebackendService.
func (s *gamebackendService) CreateMap(ctx context.Context, name string, path string, maxPlayers uint64, instanced bool) (*model.Map, error) {
	return s.gamebackendRepo.CreateMap(ctx, name, path, maxPlayers, instanced)
}

// DeleteChatTemplateById implements GamebackendService.
func (s *gamebackendService) DeleteChatTemplateById(ctx context.Context, id *uuid.UUID) error {
	return s.gamebackendRepo.DeleteChatTemplateById(ctx, id)
}

// DeleteChatTemplateByName implements GamebackendService.
func (s *gamebackendService) DeleteChatTemplateByName(ctx context.Context, name string) error {
	return s.gamebackendRepo.DeleteChatTemplateByName(ctx, name)
}

// DeleteDimensionById implements GamebackendService.
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

// DuplicateDimension implements GamebackendService.
func (s *gamebackendService) DuplicateDimension(ctx context.Context, refId *uuid.UUID, name string) (*model.Dimension, error) {
	return s.gamebackendRepo.DuplicateDimension(ctx, refId, name)
}

// EditChatTemplate implements GamebackendService.
func (s *gamebackendService) EditChatTemplate(ctx context.Context, request *pb.EditChatTemplateRequest) (*model.ChatTemplate, error) {
	var currentChatTemplate *model.ChatTemplate
	var err error

	switch target := request.Target.FindBy.(type) {
	case *pb.ChatTemplateTarget_Id:
		id, err := uuid.Parse(target.Id)
		if err != nil {
			return nil, err
		}
		currentChatTemplate, err = s.FindChatTemplateById(ctx, &id)
	case *pb.ChatTemplateTarget_Name:
		currentChatTemplate, err = s.FindChatTemplateByName(ctx, target.Name)
	default:
		log.WithContext(ctx).Errorf("target type unknown: %s", reflect.TypeOf(target).Name())
		return nil, model.ErrHandleRequest
	}

	if err != nil {
		return nil, err
	}

	if request.OptionalName != nil {
		currentChatTemplate.Name = request.OptionalName.(*pb.EditChatTemplateRequest_Name).Name
	}

	return s.gamebackendRepo.SaveChatTemplate(ctx, currentChatTemplate)
}

// EditDimension implements GamebackendService.
func (s *gamebackendService) EditDimension(ctx context.Context, request *pb.EditDimensionRequest) (*model.Dimension, error) {
	var currentDimension *model.Dimension
	var err error

	switch target := request.Target.FindBy.(type) {
	case *pb.DimensionTarget_Id:
		id, err := uuid.Parse(target.Id)
		if err != nil {
			return nil, fmt.Errorf("invalid target id: %s", target.Id)
		}
		currentDimension, err = s.FindDimensionById(ctx, &id)
	case *pb.DimensionTarget_Name:
		currentDimension, err = s.FindDimensionByName(ctx, target.Name)
	default:
		log.WithContext(ctx).Errorf("dimension target type unknown: %s", reflect.TypeOf(target).Name())
		return nil, model.ErrHandleRequest
	}

	if err != nil {
		return nil, err
	}

	if request.OptionalName != nil {
		currentDimension.Name = request.OptionalName.(*pb.EditDimensionRequest_Name).Name
	}

	if request.OptionalVersion != nil {
		currentDimension.Version = request.OptionalVersion.(*pb.EditDimensionRequest_Version).Version
	}

	if request.OptionalLocation != nil {
		currentDimension.ServerLocation = request.OptionalLocation.(*pb.EditDimensionRequest_Location).Location
	}

	if request.EditMaps {
		ids, err := helpers.ParseUUIDs(request.MapIds)
		if err != nil {
			return nil, fmt.Errorf("map ids: %w", err)
		}

		maps, err := s.gamebackendRepo.FindMapsByIds(ctx, ids)
		if err != nil {
			return nil, fmt.Errorf("getting maps: %w", err)
		}

		currentDimension.Maps = maps
	}

	if request.EditChatTemplates {
		ids, err := helpers.ParseUUIDs(request.ChatTemplateIds)
		if err != nil {
			return nil, fmt.Errorf("chat template ids: %w", err)
		}

		chatTemplates, err := s.gamebackendRepo.FindChatTemplatesByIds(ctx, ids)
		if err != nil {
			return nil, fmt.Errorf("getting chats: %w", err)
		}

		currentDimension.ChatTemplates = chatTemplates
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
		log.WithContext(ctx).Errorf("target type unknown: %s", reflect.TypeOf(target).Name())
		return nil, model.ErrHandleRequest
	}

	if err != nil {
		return nil, err
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

// FindAllChatTemplates implements GamebackendService.
func (s *gamebackendService) FindAllChatTemplates(ctx context.Context) (model.ChatTemplates, error) {
	return s.gamebackendRepo.FindAllChatTemplates(ctx)
}

// FindAllDimensions implements GamebackendService.
func (s *gamebackendService) FindAllDimensions(ctx context.Context) (model.Dimensions, error) {
	return s.gamebackendRepo.FindAllDimensions(ctx)
}

// FindAllMaps implements GamebackendService.
func (s *gamebackendService) FindAllMaps(ctx context.Context) (model.Maps, error) {
	return s.gamebackendRepo.FindAllMaps(ctx)
}

// FindChatTemplateById implements GamebackendService.
func (s *gamebackendService) FindChatTemplateById(ctx context.Context, id *uuid.UUID) (*model.ChatTemplate, error) {
	return s.gamebackendRepo.FindChatTemplateById(ctx, id)
}

// FindChatTemplateByName implements GamebackendService.
func (s *gamebackendService) FindChatTemplateByName(ctx context.Context, name string) (*model.ChatTemplate, error) {
	return s.gamebackendRepo.FindChatTemplateByName(ctx, name)
}

// FindChatTemplatesByIds implements GamebackendService.
func (s *gamebackendService) FindChatTemplatesByIds(ctx context.Context, ids []*uuid.UUID) (model.ChatTemplates, error) {
	return s.gamebackendRepo.FindChatTemplatesByIds(ctx, ids)
}

// FindChatTemplatesByNames implements GamebackendService.
func (s *gamebackendService) FindChatTemplatesByNames(ctx context.Context, names []string) (model.ChatTemplates, error) {
	return s.gamebackendRepo.FindChatTemplatesByNames(ctx, names)
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
