// Code generated by MockGen. DO NOT EDIT.
// Source: /home/wil/sro/git/go-backend/pkg/service/character_s.go
//
// Generated by this command:
//
//	mockgen -package=mocks -source=/home/wil/sro/git/go-backend/pkg/service/character_s.go -destination=/home/wil/sro/git/go-backend/pkg/mocks/character_s_mock.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	character "github.com/ShatteredRealms/go-backend/pkg/model/character"
	pb "github.com/ShatteredRealms/go-backend/pkg/pb"
	gomock "go.uber.org/mock/gomock"
)

// MockCharacterService is a mock of CharacterService interface.
type MockCharacterService struct {
	ctrl     *gomock.Controller
	recorder *MockCharacterServiceMockRecorder
}

// MockCharacterServiceMockRecorder is the mock recorder for MockCharacterService.
type MockCharacterServiceMockRecorder struct {
	mock *MockCharacterService
}

// NewMockCharacterService creates a new mock instance.
func NewMockCharacterService(ctrl *gomock.Controller) *MockCharacterService {
	mock := &MockCharacterService{ctrl: ctrl}
	mock.recorder = &MockCharacterServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCharacterService) EXPECT() *MockCharacterServiceMockRecorder {
	return m.recorder
}

// AddPlayTime mocks base method.
func (m *MockCharacterService) AddPlayTime(ctx context.Context, characterId uint, amount uint64) (*character.Character, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPlayTime", ctx, characterId, amount)
	ret0, _ := ret[0].(*character.Character)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddPlayTime indicates an expected call of AddPlayTime.
func (mr *MockCharacterServiceMockRecorder) AddPlayTime(ctx, characterId, amount any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPlayTime", reflect.TypeOf((*MockCharacterService)(nil).AddPlayTime), ctx, characterId, amount)
}

// Create mocks base method.
func (m *MockCharacterService) Create(ctx context.Context, ownerId, name, gender, realm, dimension string) (*character.Character, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, ownerId, name, gender, realm, dimension)
	ret0, _ := ret[0].(*character.Character)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockCharacterServiceMockRecorder) Create(ctx, ownerId, name, gender, realm, dimension any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCharacterService)(nil).Create), ctx, ownerId, name, gender, realm, dimension)
}

// Delete mocks base method.
func (m *MockCharacterService) Delete(ctx context.Context, id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCharacterServiceMockRecorder) Delete(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCharacterService)(nil).Delete), ctx, id)
}

// Edit mocks base method.
func (m *MockCharacterService) Edit(ctx context.Context, char *pb.EditCharacterRequest) (*character.Character, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Edit", ctx, char)
	ret0, _ := ret[0].(*character.Character)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Edit indicates an expected call of Edit.
func (mr *MockCharacterServiceMockRecorder) Edit(ctx, char any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Edit", reflect.TypeOf((*MockCharacterService)(nil).Edit), ctx, char)
}

// FindAll mocks base method.
func (m *MockCharacterService) FindAll(arg0 context.Context) (character.Characters, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", arg0)
	ret0, _ := ret[0].(character.Characters)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockCharacterServiceMockRecorder) FindAll(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockCharacterService)(nil).FindAll), arg0)
}

// FindAllByOwner mocks base method.
func (m *MockCharacterService) FindAllByOwner(ctx context.Context, ownerId string) (character.Characters, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAllByOwner", ctx, ownerId)
	ret0, _ := ret[0].(character.Characters)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAllByOwner indicates an expected call of FindAllByOwner.
func (mr *MockCharacterServiceMockRecorder) FindAllByOwner(ctx, ownerId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAllByOwner", reflect.TypeOf((*MockCharacterService)(nil).FindAllByOwner), ctx, ownerId)
}

// FindById mocks base method.
func (m *MockCharacterService) FindById(ctx context.Context, id uint) (*character.Character, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindById", ctx, id)
	ret0, _ := ret[0].(*character.Character)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindById indicates an expected call of FindById.
func (mr *MockCharacterServiceMockRecorder) FindById(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockCharacterService)(nil).FindById), ctx, id)
}

// FindByName mocks base method.
func (m *MockCharacterService) FindByName(ctx context.Context, name string) (*character.Character, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByName", ctx, name)
	ret0, _ := ret[0].(*character.Character)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByName indicates an expected call of FindByName.
func (mr *MockCharacterServiceMockRecorder) FindByName(ctx, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByName", reflect.TypeOf((*MockCharacterService)(nil).FindByName), ctx, name)
}
