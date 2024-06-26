// Code generated by MockGen. DO NOT EDIT.
// Source: /home/wil/sro/git/go-backend/pkg/repository/gamebackend_r.go
//
// Generated by this command:
//
//	mockgen -package=mocks -source=/home/wil/sro/git/go-backend/pkg/repository/gamebackend_r.go -destination=/home/wil/sro/git/go-backend/pkg/mocks/gamebackend_r_mock.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	game "github.com/ShatteredRealms/go-backend/pkg/model/game"
	gamebackend "github.com/ShatteredRealms/go-backend/pkg/model/gamebackend"
	repository "github.com/ShatteredRealms/go-backend/pkg/repository"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
	gorm "gorm.io/gorm"
)

// MockdimensionRepository is a mock of dimensionRepository interface.
type MockdimensionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockdimensionRepositoryMockRecorder
}

// MockdimensionRepositoryMockRecorder is the mock recorder for MockdimensionRepository.
type MockdimensionRepositoryMockRecorder struct {
	mock *MockdimensionRepository
}

// NewMockdimensionRepository creates a new mock instance.
func NewMockdimensionRepository(ctrl *gomock.Controller) *MockdimensionRepository {
	mock := &MockdimensionRepository{ctrl: ctrl}
	mock.recorder = &MockdimensionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockdimensionRepository) EXPECT() *MockdimensionRepositoryMockRecorder {
	return m.recorder
}

// CreateDimension mocks base method.
func (m *MockdimensionRepository) CreateDimension(ctx context.Context, name, location, version string, mapIds []*uuid.UUID) (*game.Dimension, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDimension", ctx, name, location, version, mapIds)
	ret0, _ := ret[0].(*game.Dimension)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDimension indicates an expected call of CreateDimension.
func (mr *MockdimensionRepositoryMockRecorder) CreateDimension(ctx, name, location, version, mapIds any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDimension", reflect.TypeOf((*MockdimensionRepository)(nil).CreateDimension), ctx, name, location, version, mapIds)
}

// DeleteDimensionById mocks base method.
func (m *MockdimensionRepository) DeleteDimensionById(ctx context.Context, id *uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteDimensionById", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteDimensionById indicates an expected call of DeleteDimensionById.
func (mr *MockdimensionRepositoryMockRecorder) DeleteDimensionById(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDimensionById", reflect.TypeOf((*MockdimensionRepository)(nil).DeleteDimensionById), ctx, id)
}

// DeleteDimensionByName mocks base method.
func (m *MockdimensionRepository) DeleteDimensionByName(ctx context.Context, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteDimensionByName", ctx, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteDimensionByName indicates an expected call of DeleteDimensionByName.
func (mr *MockdimensionRepositoryMockRecorder) DeleteDimensionByName(ctx, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDimensionByName", reflect.TypeOf((*MockdimensionRepository)(nil).DeleteDimensionByName), ctx, name)
}

// FindAllDimensions mocks base method.
func (m *MockdimensionRepository) FindAllDimensions(ctx context.Context) (game.Dimensions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllDimensions", ctx)
	ret0, _ := ret[0].(game.Dimensions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllDimensions indicates an expected call of FindAllDimensions.
func (mr *MockdimensionRepositoryMockRecorder) FindAllDimensions(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllDimensions", reflect.TypeOf((*MockdimensionRepository)(nil).FindAllDimensions), ctx)
}

// FindDimensionById mocks base method.
func (m *MockdimensionRepository) FindDimensionById(ctx context.Context, id *uuid.UUID) (*game.Dimension, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindDimensionById", ctx, id)
	ret0, _ := ret[0].(*game.Dimension)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindDimensionById indicates an expected call of FindDimensionById.
func (mr *MockdimensionRepositoryMockRecorder) FindDimensionById(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindDimensionById", reflect.TypeOf((*MockdimensionRepository)(nil).FindDimensionById), ctx, id)
}

// FindDimensionByName mocks base method.
func (m *MockdimensionRepository) FindDimensionByName(ctx context.Context, name string) (*game.Dimension, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindDimensionByName", ctx, name)
	ret0, _ := ret[0].(*game.Dimension)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindDimensionByName indicates an expected call of FindDimensionByName.
func (mr *MockdimensionRepositoryMockRecorder) FindDimensionByName(ctx, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindDimensionByName", reflect.TypeOf((*MockdimensionRepository)(nil).FindDimensionByName), ctx, name)
}

// FindDimensionsByIds mocks base method.
func (m *MockdimensionRepository) FindDimensionsByIds(ctx context.Context, ids []*uuid.UUID) (game.Dimensions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindDimensionsByIds", ctx, ids)
	ret0, _ := ret[0].(game.Dimensions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindDimensionsByIds indicates an expected call of FindDimensionsByIds.
func (mr *MockdimensionRepositoryMockRecorder) FindDimensionsByIds(ctx, ids any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindDimensionsByIds", reflect.TypeOf((*MockdimensionRepository)(nil).FindDimensionsByIds), ctx, ids)
}

// FindDimensionsByNames mocks base method.
func (m *MockdimensionRepository) FindDimensionsByNames(ctx context.Context, names []string) (game.Dimensions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindDimensionsByNames", ctx, names)
	ret0, _ := ret[0].(game.Dimensions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindDimensionsByNames indicates an expected call of FindDimensionsByNames.
func (mr *MockdimensionRepositoryMockRecorder) FindDimensionsByNames(ctx, names any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindDimensionsByNames", reflect.TypeOf((*MockdimensionRepository)(nil).FindDimensionsByNames), ctx, names)
}

// FindDimensionsWithMapIds mocks base method.
func (m *MockdimensionRepository) FindDimensionsWithMapIds(ctx context.Context, ids []*uuid.UUID) (game.Dimensions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindDimensionsWithMapIds", ctx, ids)
	ret0, _ := ret[0].(game.Dimensions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindDimensionsWithMapIds indicates an expected call of FindDimensionsWithMapIds.
func (mr *MockdimensionRepositoryMockRecorder) FindDimensionsWithMapIds(ctx, ids any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindDimensionsWithMapIds", reflect.TypeOf((*MockdimensionRepository)(nil).FindDimensionsWithMapIds), ctx, ids)
}

// SaveDimension mocks base method.
func (m *MockdimensionRepository) SaveDimension(ctx context.Context, dimension *game.Dimension) (*game.Dimension, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveDimension", ctx, dimension)
	ret0, _ := ret[0].(*game.Dimension)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveDimension indicates an expected call of SaveDimension.
func (mr *MockdimensionRepositoryMockRecorder) SaveDimension(ctx, dimension any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveDimension", reflect.TypeOf((*MockdimensionRepository)(nil).SaveDimension), ctx, dimension)
}

// MockmapRepository is a mock of mapRepository interface.
type MockmapRepository struct {
	ctrl     *gomock.Controller
	recorder *MockmapRepositoryMockRecorder
}

// MockmapRepositoryMockRecorder is the mock recorder for MockmapRepository.
type MockmapRepositoryMockRecorder struct {
	mock *MockmapRepository
}

// NewMockmapRepository creates a new mock instance.
func NewMockmapRepository(ctrl *gomock.Controller) *MockmapRepository {
	mock := &MockmapRepository{ctrl: ctrl}
	mock.recorder = &MockmapRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockmapRepository) EXPECT() *MockmapRepositoryMockRecorder {
	return m.recorder
}

// CreateMap mocks base method.
func (m *MockmapRepository) CreateMap(ctx context.Context, name, path string, maxPlayers uint64, instanced bool) (*game.Map, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMap", ctx, name, path, maxPlayers, instanced)
	ret0, _ := ret[0].(*game.Map)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMap indicates an expected call of CreateMap.
func (mr *MockmapRepositoryMockRecorder) CreateMap(ctx, name, path, maxPlayers, instanced any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMap", reflect.TypeOf((*MockmapRepository)(nil).CreateMap), ctx, name, path, maxPlayers, instanced)
}

// DeleteMapById mocks base method.
func (m *MockmapRepository) DeleteMapById(ctx context.Context, id *uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMapById", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMapById indicates an expected call of DeleteMapById.
func (mr *MockmapRepositoryMockRecorder) DeleteMapById(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMapById", reflect.TypeOf((*MockmapRepository)(nil).DeleteMapById), ctx, id)
}

// DeleteMapByName mocks base method.
func (m *MockmapRepository) DeleteMapByName(ctx context.Context, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMapByName", ctx, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMapByName indicates an expected call of DeleteMapByName.
func (mr *MockmapRepositoryMockRecorder) DeleteMapByName(ctx, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMapByName", reflect.TypeOf((*MockmapRepository)(nil).DeleteMapByName), ctx, name)
}

// FindAllMaps mocks base method.
func (m *MockmapRepository) FindAllMaps(ctx context.Context) (game.Maps, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllMaps", ctx)
	ret0, _ := ret[0].(game.Maps)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllMaps indicates an expected call of FindAllMaps.
func (mr *MockmapRepositoryMockRecorder) FindAllMaps(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllMaps", reflect.TypeOf((*MockmapRepository)(nil).FindAllMaps), ctx)
}

// FindMapById mocks base method.
func (m *MockmapRepository) FindMapById(ctx context.Context, id *uuid.UUID) (*game.Map, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindMapById", ctx, id)
	ret0, _ := ret[0].(*game.Map)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindMapById indicates an expected call of FindMapById.
func (mr *MockmapRepositoryMockRecorder) FindMapById(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindMapById", reflect.TypeOf((*MockmapRepository)(nil).FindMapById), ctx, id)
}

// FindMapByName mocks base method.
func (m *MockmapRepository) FindMapByName(ctx context.Context, name string) (*game.Map, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindMapByName", ctx, name)
	ret0, _ := ret[0].(*game.Map)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindMapByName indicates an expected call of FindMapByName.
func (mr *MockmapRepositoryMockRecorder) FindMapByName(ctx, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindMapByName", reflect.TypeOf((*MockmapRepository)(nil).FindMapByName), ctx, name)
}

// FindMapsByIds mocks base method.
func (m *MockmapRepository) FindMapsByIds(ctx context.Context, ids []*uuid.UUID) (game.Maps, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindMapsByIds", ctx, ids)
	ret0, _ := ret[0].(game.Maps)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindMapsByIds indicates an expected call of FindMapsByIds.
func (mr *MockmapRepositoryMockRecorder) FindMapsByIds(ctx, ids any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindMapsByIds", reflect.TypeOf((*MockmapRepository)(nil).FindMapsByIds), ctx, ids)
}

// FindMapsByNames mocks base method.
func (m *MockmapRepository) FindMapsByNames(ctx context.Context, names []string) (game.Maps, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindMapsByNames", ctx, names)
	ret0, _ := ret[0].(game.Maps)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindMapsByNames indicates an expected call of FindMapsByNames.
func (mr *MockmapRepositoryMockRecorder) FindMapsByNames(ctx, names any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindMapsByNames", reflect.TypeOf((*MockmapRepository)(nil).FindMapsByNames), ctx, names)
}

// SaveMap mocks base method.
func (m_2 *MockmapRepository) SaveMap(ctx context.Context, m *game.Map) (*game.Map, error) {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "SaveMap", ctx, m)
	ret0, _ := ret[0].(*game.Map)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveMap indicates an expected call of SaveMap.
func (mr *MockmapRepositoryMockRecorder) SaveMap(ctx, m any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveMap", reflect.TypeOf((*MockmapRepository)(nil).SaveMap), ctx, m)
}

// MockGamebackendRepository is a mock of GamebackendRepository interface.
type MockGamebackendRepository struct {
	ctrl     *gomock.Controller
	recorder *MockGamebackendRepositoryMockRecorder
}

// MockGamebackendRepositoryMockRecorder is the mock recorder for MockGamebackendRepository.
type MockGamebackendRepositoryMockRecorder struct {
	mock *MockGamebackendRepository
}

// NewMockGamebackendRepository creates a new mock instance.
func NewMockGamebackendRepository(ctrl *gomock.Controller) *MockGamebackendRepository {
	mock := &MockGamebackendRepository{ctrl: ctrl}
	mock.recorder = &MockGamebackendRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGamebackendRepository) EXPECT() *MockGamebackendRepositoryMockRecorder {
	return m.recorder
}

// CreateDimension mocks base method.
func (m *MockGamebackendRepository) CreateDimension(ctx context.Context, name, location, version string, mapIds []*uuid.UUID) (*game.Dimension, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDimension", ctx, name, location, version, mapIds)
	ret0, _ := ret[0].(*game.Dimension)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDimension indicates an expected call of CreateDimension.
func (mr *MockGamebackendRepositoryMockRecorder) CreateDimension(ctx, name, location, version, mapIds any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDimension", reflect.TypeOf((*MockGamebackendRepository)(nil).CreateDimension), ctx, name, location, version, mapIds)
}

// CreateMap mocks base method.
func (m *MockGamebackendRepository) CreateMap(ctx context.Context, name, path string, maxPlayers uint64, instanced bool) (*game.Map, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMap", ctx, name, path, maxPlayers, instanced)
	ret0, _ := ret[0].(*game.Map)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMap indicates an expected call of CreateMap.
func (mr *MockGamebackendRepositoryMockRecorder) CreateMap(ctx, name, path, maxPlayers, instanced any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMap", reflect.TypeOf((*MockGamebackendRepository)(nil).CreateMap), ctx, name, path, maxPlayers, instanced)
}

// CreatePendingConnection mocks base method.
func (m *MockGamebackendRepository) CreatePendingConnection(ctx context.Context, character, serverName string) (*gamebackend.PendingConnection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePendingConnection", ctx, character, serverName)
	ret0, _ := ret[0].(*gamebackend.PendingConnection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePendingConnection indicates an expected call of CreatePendingConnection.
func (mr *MockGamebackendRepositoryMockRecorder) CreatePendingConnection(ctx, character, serverName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePendingConnection", reflect.TypeOf((*MockGamebackendRepository)(nil).CreatePendingConnection), ctx, character, serverName)
}

// DeleteDimensionById mocks base method.
func (m *MockGamebackendRepository) DeleteDimensionById(ctx context.Context, id *uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteDimensionById", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteDimensionById indicates an expected call of DeleteDimensionById.
func (mr *MockGamebackendRepositoryMockRecorder) DeleteDimensionById(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDimensionById", reflect.TypeOf((*MockGamebackendRepository)(nil).DeleteDimensionById), ctx, id)
}

// DeleteDimensionByName mocks base method.
func (m *MockGamebackendRepository) DeleteDimensionByName(ctx context.Context, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteDimensionByName", ctx, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteDimensionByName indicates an expected call of DeleteDimensionByName.
func (mr *MockGamebackendRepositoryMockRecorder) DeleteDimensionByName(ctx, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDimensionByName", reflect.TypeOf((*MockGamebackendRepository)(nil).DeleteDimensionByName), ctx, name)
}

// DeleteMapById mocks base method.
func (m *MockGamebackendRepository) DeleteMapById(ctx context.Context, id *uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMapById", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMapById indicates an expected call of DeleteMapById.
func (mr *MockGamebackendRepositoryMockRecorder) DeleteMapById(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMapById", reflect.TypeOf((*MockGamebackendRepository)(nil).DeleteMapById), ctx, id)
}

// DeleteMapByName mocks base method.
func (m *MockGamebackendRepository) DeleteMapByName(ctx context.Context, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMapByName", ctx, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMapByName indicates an expected call of DeleteMapByName.
func (mr *MockGamebackendRepositoryMockRecorder) DeleteMapByName(ctx, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMapByName", reflect.TypeOf((*MockGamebackendRepository)(nil).DeleteMapByName), ctx, name)
}

// DeletePendingConnection mocks base method.
func (m *MockGamebackendRepository) DeletePendingConnection(ctx context.Context, id *uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePendingConnection", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePendingConnection indicates an expected call of DeletePendingConnection.
func (mr *MockGamebackendRepositoryMockRecorder) DeletePendingConnection(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePendingConnection", reflect.TypeOf((*MockGamebackendRepository)(nil).DeletePendingConnection), ctx, id)
}

// FindAllDimensions mocks base method.
func (m *MockGamebackendRepository) FindAllDimensions(ctx context.Context) (game.Dimensions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllDimensions", ctx)
	ret0, _ := ret[0].(game.Dimensions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllDimensions indicates an expected call of FindAllDimensions.
func (mr *MockGamebackendRepositoryMockRecorder) FindAllDimensions(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllDimensions", reflect.TypeOf((*MockGamebackendRepository)(nil).FindAllDimensions), ctx)
}

// FindAllMaps mocks base method.
func (m *MockGamebackendRepository) FindAllMaps(ctx context.Context) (game.Maps, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllMaps", ctx)
	ret0, _ := ret[0].(game.Maps)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllMaps indicates an expected call of FindAllMaps.
func (mr *MockGamebackendRepositoryMockRecorder) FindAllMaps(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllMaps", reflect.TypeOf((*MockGamebackendRepository)(nil).FindAllMaps), ctx)
}

// FindDimensionById mocks base method.
func (m *MockGamebackendRepository) FindDimensionById(ctx context.Context, id *uuid.UUID) (*game.Dimension, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindDimensionById", ctx, id)
	ret0, _ := ret[0].(*game.Dimension)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindDimensionById indicates an expected call of FindDimensionById.
func (mr *MockGamebackendRepositoryMockRecorder) FindDimensionById(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindDimensionById", reflect.TypeOf((*MockGamebackendRepository)(nil).FindDimensionById), ctx, id)
}

// FindDimensionByName mocks base method.
func (m *MockGamebackendRepository) FindDimensionByName(ctx context.Context, name string) (*game.Dimension, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindDimensionByName", ctx, name)
	ret0, _ := ret[0].(*game.Dimension)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindDimensionByName indicates an expected call of FindDimensionByName.
func (mr *MockGamebackendRepositoryMockRecorder) FindDimensionByName(ctx, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindDimensionByName", reflect.TypeOf((*MockGamebackendRepository)(nil).FindDimensionByName), ctx, name)
}

// FindDimensionsByIds mocks base method.
func (m *MockGamebackendRepository) FindDimensionsByIds(ctx context.Context, ids []*uuid.UUID) (game.Dimensions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindDimensionsByIds", ctx, ids)
	ret0, _ := ret[0].(game.Dimensions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindDimensionsByIds indicates an expected call of FindDimensionsByIds.
func (mr *MockGamebackendRepositoryMockRecorder) FindDimensionsByIds(ctx, ids any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindDimensionsByIds", reflect.TypeOf((*MockGamebackendRepository)(nil).FindDimensionsByIds), ctx, ids)
}

// FindDimensionsByNames mocks base method.
func (m *MockGamebackendRepository) FindDimensionsByNames(ctx context.Context, names []string) (game.Dimensions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindDimensionsByNames", ctx, names)
	ret0, _ := ret[0].(game.Dimensions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindDimensionsByNames indicates an expected call of FindDimensionsByNames.
func (mr *MockGamebackendRepositoryMockRecorder) FindDimensionsByNames(ctx, names any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindDimensionsByNames", reflect.TypeOf((*MockGamebackendRepository)(nil).FindDimensionsByNames), ctx, names)
}

// FindDimensionsWithMapIds mocks base method.
func (m *MockGamebackendRepository) FindDimensionsWithMapIds(ctx context.Context, ids []*uuid.UUID) (game.Dimensions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindDimensionsWithMapIds", ctx, ids)
	ret0, _ := ret[0].(game.Dimensions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindDimensionsWithMapIds indicates an expected call of FindDimensionsWithMapIds.
func (mr *MockGamebackendRepositoryMockRecorder) FindDimensionsWithMapIds(ctx, ids any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindDimensionsWithMapIds", reflect.TypeOf((*MockGamebackendRepository)(nil).FindDimensionsWithMapIds), ctx, ids)
}

// FindMapById mocks base method.
func (m *MockGamebackendRepository) FindMapById(ctx context.Context, id *uuid.UUID) (*game.Map, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindMapById", ctx, id)
	ret0, _ := ret[0].(*game.Map)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindMapById indicates an expected call of FindMapById.
func (mr *MockGamebackendRepositoryMockRecorder) FindMapById(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindMapById", reflect.TypeOf((*MockGamebackendRepository)(nil).FindMapById), ctx, id)
}

// FindMapByName mocks base method.
func (m *MockGamebackendRepository) FindMapByName(ctx context.Context, name string) (*game.Map, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindMapByName", ctx, name)
	ret0, _ := ret[0].(*game.Map)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindMapByName indicates an expected call of FindMapByName.
func (mr *MockGamebackendRepositoryMockRecorder) FindMapByName(ctx, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindMapByName", reflect.TypeOf((*MockGamebackendRepository)(nil).FindMapByName), ctx, name)
}

// FindMapsByIds mocks base method.
func (m *MockGamebackendRepository) FindMapsByIds(ctx context.Context, ids []*uuid.UUID) (game.Maps, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindMapsByIds", ctx, ids)
	ret0, _ := ret[0].(game.Maps)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindMapsByIds indicates an expected call of FindMapsByIds.
func (mr *MockGamebackendRepositoryMockRecorder) FindMapsByIds(ctx, ids any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindMapsByIds", reflect.TypeOf((*MockGamebackendRepository)(nil).FindMapsByIds), ctx, ids)
}

// FindMapsByNames mocks base method.
func (m *MockGamebackendRepository) FindMapsByNames(ctx context.Context, names []string) (game.Maps, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindMapsByNames", ctx, names)
	ret0, _ := ret[0].(game.Maps)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindMapsByNames indicates an expected call of FindMapsByNames.
func (mr *MockGamebackendRepositoryMockRecorder) FindMapsByNames(ctx, names any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindMapsByNames", reflect.TypeOf((*MockGamebackendRepository)(nil).FindMapsByNames), ctx, names)
}

// FindPendingConnection mocks base method.
func (m *MockGamebackendRepository) FindPendingConnection(ctx context.Context, id *uuid.UUID) *gamebackend.PendingConnection {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindPendingConnection", ctx, id)
	ret0, _ := ret[0].(*gamebackend.PendingConnection)
	return ret0
}

// FindPendingConnection indicates an expected call of FindPendingConnection.
func (mr *MockGamebackendRepositoryMockRecorder) FindPendingConnection(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindPendingConnection", reflect.TypeOf((*MockGamebackendRepository)(nil).FindPendingConnection), ctx, id)
}

// Migrate mocks base method.
func (m *MockGamebackendRepository) Migrate(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Migrate", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Migrate indicates an expected call of Migrate.
func (mr *MockGamebackendRepositoryMockRecorder) Migrate(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Migrate", reflect.TypeOf((*MockGamebackendRepository)(nil).Migrate), ctx)
}

// SaveDimension mocks base method.
func (m *MockGamebackendRepository) SaveDimension(ctx context.Context, dimension *game.Dimension) (*game.Dimension, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveDimension", ctx, dimension)
	ret0, _ := ret[0].(*game.Dimension)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveDimension indicates an expected call of SaveDimension.
func (mr *MockGamebackendRepositoryMockRecorder) SaveDimension(ctx, dimension any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveDimension", reflect.TypeOf((*MockGamebackendRepository)(nil).SaveDimension), ctx, dimension)
}

// SaveMap mocks base method.
func (m_2 *MockGamebackendRepository) SaveMap(ctx context.Context, m *game.Map) (*game.Map, error) {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "SaveMap", ctx, m)
	ret0, _ := ret[0].(*game.Map)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveMap indicates an expected call of SaveMap.
func (mr *MockGamebackendRepositoryMockRecorder) SaveMap(ctx, m any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveMap", reflect.TypeOf((*MockGamebackendRepository)(nil).SaveMap), ctx, m)
}

// WithTrx mocks base method.
func (m *MockGamebackendRepository) WithTrx(trx *gorm.DB) repository.GamebackendRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTrx", trx)
	ret0, _ := ret[0].(repository.GamebackendRepository)
	return ret0
}

// WithTrx indicates an expected call of WithTrx.
func (mr *MockGamebackendRepositoryMockRecorder) WithTrx(trx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTrx", reflect.TypeOf((*MockGamebackendRepository)(nil).WithTrx), trx)
}
