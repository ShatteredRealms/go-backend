// Code generated by MockGen. DO NOT EDIT.
// Source: /home/wil/sro/git/go-backend/pkg/pb/servermanager.pb.go
//
// Generated by this command:
//
//	mockgen -package=mocks -source=/home/wil/sro/git/go-backend/pkg/pb/servermanager.pb.go -destination=/home/wil/sro/git/go-backend/pkg/mocks/servermanager.pb_mock.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockisDimensionTarget_FindBy is a mock of isDimensionTarget_FindBy interface.
type MockisDimensionTarget_FindBy struct {
	ctrl     *gomock.Controller
	recorder *MockisDimensionTarget_FindByMockRecorder
}

// MockisDimensionTarget_FindByMockRecorder is the mock recorder for MockisDimensionTarget_FindBy.
type MockisDimensionTarget_FindByMockRecorder struct {
	mock *MockisDimensionTarget_FindBy
}

// NewMockisDimensionTarget_FindBy creates a new mock instance.
func NewMockisDimensionTarget_FindBy(ctrl *gomock.Controller) *MockisDimensionTarget_FindBy {
	mock := &MockisDimensionTarget_FindBy{ctrl: ctrl}
	mock.recorder = &MockisDimensionTarget_FindByMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockisDimensionTarget_FindBy) EXPECT() *MockisDimensionTarget_FindByMockRecorder {
	return m.recorder
}

// isDimensionTarget_FindBy mocks base method.
func (m *MockisDimensionTarget_FindBy) isDimensionTarget_FindBy() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "isDimensionTarget_FindBy")
}

// isDimensionTarget_FindBy indicates an expected call of isDimensionTarget_FindBy.
func (mr *MockisDimensionTarget_FindByMockRecorder) isDimensionTarget_FindBy() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isDimensionTarget_FindBy", reflect.TypeOf((*MockisDimensionTarget_FindBy)(nil).isDimensionTarget_FindBy))
}

// MockisMapTarget_FindBy is a mock of isMapTarget_FindBy interface.
type MockisMapTarget_FindBy struct {
	ctrl     *gomock.Controller
	recorder *MockisMapTarget_FindByMockRecorder
}

// MockisMapTarget_FindByMockRecorder is the mock recorder for MockisMapTarget_FindBy.
type MockisMapTarget_FindByMockRecorder struct {
	mock *MockisMapTarget_FindBy
}

// NewMockisMapTarget_FindBy creates a new mock instance.
func NewMockisMapTarget_FindBy(ctrl *gomock.Controller) *MockisMapTarget_FindBy {
	mock := &MockisMapTarget_FindBy{ctrl: ctrl}
	mock.recorder = &MockisMapTarget_FindByMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockisMapTarget_FindBy) EXPECT() *MockisMapTarget_FindByMockRecorder {
	return m.recorder
}

// isMapTarget_FindBy mocks base method.
func (m *MockisMapTarget_FindBy) isMapTarget_FindBy() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "isMapTarget_FindBy")
}

// isMapTarget_FindBy indicates an expected call of isMapTarget_FindBy.
func (mr *MockisMapTarget_FindByMockRecorder) isMapTarget_FindBy() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isMapTarget_FindBy", reflect.TypeOf((*MockisMapTarget_FindBy)(nil).isMapTarget_FindBy))
}

// MockisEditDimensionRequest_OptionalName is a mock of isEditDimensionRequest_OptionalName interface.
type MockisEditDimensionRequest_OptionalName struct {
	ctrl     *gomock.Controller
	recorder *MockisEditDimensionRequest_OptionalNameMockRecorder
}

// MockisEditDimensionRequest_OptionalNameMockRecorder is the mock recorder for MockisEditDimensionRequest_OptionalName.
type MockisEditDimensionRequest_OptionalNameMockRecorder struct {
	mock *MockisEditDimensionRequest_OptionalName
}

// NewMockisEditDimensionRequest_OptionalName creates a new mock instance.
func NewMockisEditDimensionRequest_OptionalName(ctrl *gomock.Controller) *MockisEditDimensionRequest_OptionalName {
	mock := &MockisEditDimensionRequest_OptionalName{ctrl: ctrl}
	mock.recorder = &MockisEditDimensionRequest_OptionalNameMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockisEditDimensionRequest_OptionalName) EXPECT() *MockisEditDimensionRequest_OptionalNameMockRecorder {
	return m.recorder
}

// isEditDimensionRequest_OptionalName mocks base method.
func (m *MockisEditDimensionRequest_OptionalName) isEditDimensionRequest_OptionalName() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "isEditDimensionRequest_OptionalName")
}

// isEditDimensionRequest_OptionalName indicates an expected call of isEditDimensionRequest_OptionalName.
func (mr *MockisEditDimensionRequest_OptionalNameMockRecorder) isEditDimensionRequest_OptionalName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isEditDimensionRequest_OptionalName", reflect.TypeOf((*MockisEditDimensionRequest_OptionalName)(nil).isEditDimensionRequest_OptionalName))
}

// MockisEditDimensionRequest_OptionalVersion is a mock of isEditDimensionRequest_OptionalVersion interface.
type MockisEditDimensionRequest_OptionalVersion struct {
	ctrl     *gomock.Controller
	recorder *MockisEditDimensionRequest_OptionalVersionMockRecorder
}

// MockisEditDimensionRequest_OptionalVersionMockRecorder is the mock recorder for MockisEditDimensionRequest_OptionalVersion.
type MockisEditDimensionRequest_OptionalVersionMockRecorder struct {
	mock *MockisEditDimensionRequest_OptionalVersion
}

// NewMockisEditDimensionRequest_OptionalVersion creates a new mock instance.
func NewMockisEditDimensionRequest_OptionalVersion(ctrl *gomock.Controller) *MockisEditDimensionRequest_OptionalVersion {
	mock := &MockisEditDimensionRequest_OptionalVersion{ctrl: ctrl}
	mock.recorder = &MockisEditDimensionRequest_OptionalVersionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockisEditDimensionRequest_OptionalVersion) EXPECT() *MockisEditDimensionRequest_OptionalVersionMockRecorder {
	return m.recorder
}

// isEditDimensionRequest_OptionalVersion mocks base method.
func (m *MockisEditDimensionRequest_OptionalVersion) isEditDimensionRequest_OptionalVersion() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "isEditDimensionRequest_OptionalVersion")
}

// isEditDimensionRequest_OptionalVersion indicates an expected call of isEditDimensionRequest_OptionalVersion.
func (mr *MockisEditDimensionRequest_OptionalVersionMockRecorder) isEditDimensionRequest_OptionalVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isEditDimensionRequest_OptionalVersion", reflect.TypeOf((*MockisEditDimensionRequest_OptionalVersion)(nil).isEditDimensionRequest_OptionalVersion))
}

// MockisEditDimensionRequest_OptionalLocation is a mock of isEditDimensionRequest_OptionalLocation interface.
type MockisEditDimensionRequest_OptionalLocation struct {
	ctrl     *gomock.Controller
	recorder *MockisEditDimensionRequest_OptionalLocationMockRecorder
}

// MockisEditDimensionRequest_OptionalLocationMockRecorder is the mock recorder for MockisEditDimensionRequest_OptionalLocation.
type MockisEditDimensionRequest_OptionalLocationMockRecorder struct {
	mock *MockisEditDimensionRequest_OptionalLocation
}

// NewMockisEditDimensionRequest_OptionalLocation creates a new mock instance.
func NewMockisEditDimensionRequest_OptionalLocation(ctrl *gomock.Controller) *MockisEditDimensionRequest_OptionalLocation {
	mock := &MockisEditDimensionRequest_OptionalLocation{ctrl: ctrl}
	mock.recorder = &MockisEditDimensionRequest_OptionalLocationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockisEditDimensionRequest_OptionalLocation) EXPECT() *MockisEditDimensionRequest_OptionalLocationMockRecorder {
	return m.recorder
}

// isEditDimensionRequest_OptionalLocation mocks base method.
func (m *MockisEditDimensionRequest_OptionalLocation) isEditDimensionRequest_OptionalLocation() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "isEditDimensionRequest_OptionalLocation")
}

// isEditDimensionRequest_OptionalLocation indicates an expected call of isEditDimensionRequest_OptionalLocation.
func (mr *MockisEditDimensionRequest_OptionalLocationMockRecorder) isEditDimensionRequest_OptionalLocation() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isEditDimensionRequest_OptionalLocation", reflect.TypeOf((*MockisEditDimensionRequest_OptionalLocation)(nil).isEditDimensionRequest_OptionalLocation))
}

// MockisEditMapRequest_OptionalName is a mock of isEditMapRequest_OptionalName interface.
type MockisEditMapRequest_OptionalName struct {
	ctrl     *gomock.Controller
	recorder *MockisEditMapRequest_OptionalNameMockRecorder
}

// MockisEditMapRequest_OptionalNameMockRecorder is the mock recorder for MockisEditMapRequest_OptionalName.
type MockisEditMapRequest_OptionalNameMockRecorder struct {
	mock *MockisEditMapRequest_OptionalName
}

// NewMockisEditMapRequest_OptionalName creates a new mock instance.
func NewMockisEditMapRequest_OptionalName(ctrl *gomock.Controller) *MockisEditMapRequest_OptionalName {
	mock := &MockisEditMapRequest_OptionalName{ctrl: ctrl}
	mock.recorder = &MockisEditMapRequest_OptionalNameMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockisEditMapRequest_OptionalName) EXPECT() *MockisEditMapRequest_OptionalNameMockRecorder {
	return m.recorder
}

// isEditMapRequest_OptionalName mocks base method.
func (m *MockisEditMapRequest_OptionalName) isEditMapRequest_OptionalName() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "isEditMapRequest_OptionalName")
}

// isEditMapRequest_OptionalName indicates an expected call of isEditMapRequest_OptionalName.
func (mr *MockisEditMapRequest_OptionalNameMockRecorder) isEditMapRequest_OptionalName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isEditMapRequest_OptionalName", reflect.TypeOf((*MockisEditMapRequest_OptionalName)(nil).isEditMapRequest_OptionalName))
}

// MockisEditMapRequest_OptionalPath is a mock of isEditMapRequest_OptionalPath interface.
type MockisEditMapRequest_OptionalPath struct {
	ctrl     *gomock.Controller
	recorder *MockisEditMapRequest_OptionalPathMockRecorder
}

// MockisEditMapRequest_OptionalPathMockRecorder is the mock recorder for MockisEditMapRequest_OptionalPath.
type MockisEditMapRequest_OptionalPathMockRecorder struct {
	mock *MockisEditMapRequest_OptionalPath
}

// NewMockisEditMapRequest_OptionalPath creates a new mock instance.
func NewMockisEditMapRequest_OptionalPath(ctrl *gomock.Controller) *MockisEditMapRequest_OptionalPath {
	mock := &MockisEditMapRequest_OptionalPath{ctrl: ctrl}
	mock.recorder = &MockisEditMapRequest_OptionalPathMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockisEditMapRequest_OptionalPath) EXPECT() *MockisEditMapRequest_OptionalPathMockRecorder {
	return m.recorder
}

// isEditMapRequest_OptionalPath mocks base method.
func (m *MockisEditMapRequest_OptionalPath) isEditMapRequest_OptionalPath() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "isEditMapRequest_OptionalPath")
}

// isEditMapRequest_OptionalPath indicates an expected call of isEditMapRequest_OptionalPath.
func (mr *MockisEditMapRequest_OptionalPathMockRecorder) isEditMapRequest_OptionalPath() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isEditMapRequest_OptionalPath", reflect.TypeOf((*MockisEditMapRequest_OptionalPath)(nil).isEditMapRequest_OptionalPath))
}

// MockisEditMapRequest_OptionalMaxPlayers is a mock of isEditMapRequest_OptionalMaxPlayers interface.
type MockisEditMapRequest_OptionalMaxPlayers struct {
	ctrl     *gomock.Controller
	recorder *MockisEditMapRequest_OptionalMaxPlayersMockRecorder
}

// MockisEditMapRequest_OptionalMaxPlayersMockRecorder is the mock recorder for MockisEditMapRequest_OptionalMaxPlayers.
type MockisEditMapRequest_OptionalMaxPlayersMockRecorder struct {
	mock *MockisEditMapRequest_OptionalMaxPlayers
}

// NewMockisEditMapRequest_OptionalMaxPlayers creates a new mock instance.
func NewMockisEditMapRequest_OptionalMaxPlayers(ctrl *gomock.Controller) *MockisEditMapRequest_OptionalMaxPlayers {
	mock := &MockisEditMapRequest_OptionalMaxPlayers{ctrl: ctrl}
	mock.recorder = &MockisEditMapRequest_OptionalMaxPlayersMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockisEditMapRequest_OptionalMaxPlayers) EXPECT() *MockisEditMapRequest_OptionalMaxPlayersMockRecorder {
	return m.recorder
}

// isEditMapRequest_OptionalMaxPlayers mocks base method.
func (m *MockisEditMapRequest_OptionalMaxPlayers) isEditMapRequest_OptionalMaxPlayers() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "isEditMapRequest_OptionalMaxPlayers")
}

// isEditMapRequest_OptionalMaxPlayers indicates an expected call of isEditMapRequest_OptionalMaxPlayers.
func (mr *MockisEditMapRequest_OptionalMaxPlayersMockRecorder) isEditMapRequest_OptionalMaxPlayers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isEditMapRequest_OptionalMaxPlayers", reflect.TypeOf((*MockisEditMapRequest_OptionalMaxPlayers)(nil).isEditMapRequest_OptionalMaxPlayers))
}

// MockisEditMapRequest_OptionalInstanced is a mock of isEditMapRequest_OptionalInstanced interface.
type MockisEditMapRequest_OptionalInstanced struct {
	ctrl     *gomock.Controller
	recorder *MockisEditMapRequest_OptionalInstancedMockRecorder
}

// MockisEditMapRequest_OptionalInstancedMockRecorder is the mock recorder for MockisEditMapRequest_OptionalInstanced.
type MockisEditMapRequest_OptionalInstancedMockRecorder struct {
	mock *MockisEditMapRequest_OptionalInstanced
}

// NewMockisEditMapRequest_OptionalInstanced creates a new mock instance.
func NewMockisEditMapRequest_OptionalInstanced(ctrl *gomock.Controller) *MockisEditMapRequest_OptionalInstanced {
	mock := &MockisEditMapRequest_OptionalInstanced{ctrl: ctrl}
	mock.recorder = &MockisEditMapRequest_OptionalInstancedMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockisEditMapRequest_OptionalInstanced) EXPECT() *MockisEditMapRequest_OptionalInstancedMockRecorder {
	return m.recorder
}

// isEditMapRequest_OptionalInstanced mocks base method.
func (m *MockisEditMapRequest_OptionalInstanced) isEditMapRequest_OptionalInstanced() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "isEditMapRequest_OptionalInstanced")
}

// isEditMapRequest_OptionalInstanced indicates an expected call of isEditMapRequest_OptionalInstanced.
func (mr *MockisEditMapRequest_OptionalInstancedMockRecorder) isEditMapRequest_OptionalInstanced() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isEditMapRequest_OptionalInstanced", reflect.TypeOf((*MockisEditMapRequest_OptionalInstanced)(nil).isEditMapRequest_OptionalInstanced))
}
