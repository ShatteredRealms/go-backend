// Code generated by MockGen. DO NOT EDIT.
// Source: /home/wil/sro/git/go-backend/pkg/service/inventory_s.go
//
// Generated by this command:
//
//	mockgen -package=mocks -source=/home/wil/sro/git/go-backend/pkg/service/inventory_s.go -destination=/home/wil/sro/git/go-backend/pkg/mocks/inventory_s_mock.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	character "github.com/ShatteredRealms/go-backend/pkg/model/character"
	gomock "go.uber.org/mock/gomock"
)

// MockInventoryService is a mock of InventoryService interface.
type MockInventoryService struct {
	ctrl     *gomock.Controller
	recorder *MockInventoryServiceMockRecorder
}

// MockInventoryServiceMockRecorder is the mock recorder for MockInventoryService.
type MockInventoryServiceMockRecorder struct {
	mock *MockInventoryService
}

// NewMockInventoryService creates a new mock instance.
func NewMockInventoryService(ctrl *gomock.Controller) *MockInventoryService {
	mock := &MockInventoryService{ctrl: ctrl}
	mock.recorder = &MockInventoryServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInventoryService) EXPECT() *MockInventoryServiceMockRecorder {
	return m.recorder
}

// GetInventory mocks base method.
func (m *MockInventoryService) GetInventory(ctx context.Context, characterId uint) (*character.Inventory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInventory", ctx, characterId)
	ret0, _ := ret[0].(*character.Inventory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInventory indicates an expected call of GetInventory.
func (mr *MockInventoryServiceMockRecorder) GetInventory(ctx, characterId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInventory", reflect.TypeOf((*MockInventoryService)(nil).GetInventory), ctx, characterId)
}

// UpdateInventory mocks base method.
func (m *MockInventoryService) UpdateInventory(ctx context.Context, inventory *character.Inventory) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInventory", ctx, inventory)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateInventory indicates an expected call of UpdateInventory.
func (mr *MockInventoryServiceMockRecorder) UpdateInventory(ctx, inventory any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInventory", reflect.TypeOf((*MockInventoryService)(nil).UpdateInventory), ctx, inventory)
}
