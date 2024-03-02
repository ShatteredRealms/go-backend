// Code generated by MockGen. DO NOT EDIT.
// Source: /home/wil/sro/git/go-backend/pkg/pb/character.pb.go
//
// Generated by this command:
//
//	mockgen -package=mocks -source=/home/wil/sro/git/go-backend/pkg/pb/character.pb.go -destination=/home/wil/sro/git/go-backend/pkg/mocks/character.pb_mock.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockisCharacterTarget_Type is a mock of isCharacterTarget_Type interface.
type MockisCharacterTarget_Type struct {
	ctrl     *gomock.Controller
	recorder *MockisCharacterTarget_TypeMockRecorder
}

// MockisCharacterTarget_TypeMockRecorder is the mock recorder for MockisCharacterTarget_Type.
type MockisCharacterTarget_TypeMockRecorder struct {
	mock *MockisCharacterTarget_Type
}

// NewMockisCharacterTarget_Type creates a new mock instance.
func NewMockisCharacterTarget_Type(ctrl *gomock.Controller) *MockisCharacterTarget_Type {
	mock := &MockisCharacterTarget_Type{ctrl: ctrl}
	mock.recorder = &MockisCharacterTarget_TypeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockisCharacterTarget_Type) EXPECT() *MockisCharacterTarget_TypeMockRecorder {
	return m.recorder
}

// isCharacterTarget_Type mocks base method.
func (m *MockisCharacterTarget_Type) isCharacterTarget_Type() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "isCharacterTarget_Type")
}

// isCharacterTarget_Type indicates an expected call of isCharacterTarget_Type.
func (mr *MockisCharacterTarget_TypeMockRecorder) isCharacterTarget_Type() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isCharacterTarget_Type", reflect.TypeOf((*MockisCharacterTarget_Type)(nil).isCharacterTarget_Type))
}

// MockisEditCharacterRequest_OptionalOwnerId is a mock of isEditCharacterRequest_OptionalOwnerId interface.
type MockisEditCharacterRequest_OptionalOwnerId struct {
	ctrl     *gomock.Controller
	recorder *MockisEditCharacterRequest_OptionalOwnerIdMockRecorder
}

// MockisEditCharacterRequest_OptionalOwnerIdMockRecorder is the mock recorder for MockisEditCharacterRequest_OptionalOwnerId.
type MockisEditCharacterRequest_OptionalOwnerIdMockRecorder struct {
	mock *MockisEditCharacterRequest_OptionalOwnerId
}

// NewMockisEditCharacterRequest_OptionalOwnerId creates a new mock instance.
func NewMockisEditCharacterRequest_OptionalOwnerId(ctrl *gomock.Controller) *MockisEditCharacterRequest_OptionalOwnerId {
	mock := &MockisEditCharacterRequest_OptionalOwnerId{ctrl: ctrl}
	mock.recorder = &MockisEditCharacterRequest_OptionalOwnerIdMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockisEditCharacterRequest_OptionalOwnerId) EXPECT() *MockisEditCharacterRequest_OptionalOwnerIdMockRecorder {
	return m.recorder
}

// isEditCharacterRequest_OptionalOwnerId mocks base method.
func (m *MockisEditCharacterRequest_OptionalOwnerId) isEditCharacterRequest_OptionalOwnerId() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "isEditCharacterRequest_OptionalOwnerId")
}

// isEditCharacterRequest_OptionalOwnerId indicates an expected call of isEditCharacterRequest_OptionalOwnerId.
func (mr *MockisEditCharacterRequest_OptionalOwnerIdMockRecorder) isEditCharacterRequest_OptionalOwnerId() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isEditCharacterRequest_OptionalOwnerId", reflect.TypeOf((*MockisEditCharacterRequest_OptionalOwnerId)(nil).isEditCharacterRequest_OptionalOwnerId))
}

// MockisEditCharacterRequest_OptionalNewName is a mock of isEditCharacterRequest_OptionalNewName interface.
type MockisEditCharacterRequest_OptionalNewName struct {
	ctrl     *gomock.Controller
	recorder *MockisEditCharacterRequest_OptionalNewNameMockRecorder
}

// MockisEditCharacterRequest_OptionalNewNameMockRecorder is the mock recorder for MockisEditCharacterRequest_OptionalNewName.
type MockisEditCharacterRequest_OptionalNewNameMockRecorder struct {
	mock *MockisEditCharacterRequest_OptionalNewName
}

// NewMockisEditCharacterRequest_OptionalNewName creates a new mock instance.
func NewMockisEditCharacterRequest_OptionalNewName(ctrl *gomock.Controller) *MockisEditCharacterRequest_OptionalNewName {
	mock := &MockisEditCharacterRequest_OptionalNewName{ctrl: ctrl}
	mock.recorder = &MockisEditCharacterRequest_OptionalNewNameMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockisEditCharacterRequest_OptionalNewName) EXPECT() *MockisEditCharacterRequest_OptionalNewNameMockRecorder {
	return m.recorder
}

// isEditCharacterRequest_OptionalNewName mocks base method.
func (m *MockisEditCharacterRequest_OptionalNewName) isEditCharacterRequest_OptionalNewName() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "isEditCharacterRequest_OptionalNewName")
}

// isEditCharacterRequest_OptionalNewName indicates an expected call of isEditCharacterRequest_OptionalNewName.
func (mr *MockisEditCharacterRequest_OptionalNewNameMockRecorder) isEditCharacterRequest_OptionalNewName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isEditCharacterRequest_OptionalNewName", reflect.TypeOf((*MockisEditCharacterRequest_OptionalNewName)(nil).isEditCharacterRequest_OptionalNewName))
}

// MockisEditCharacterRequest_OptionalGender is a mock of isEditCharacterRequest_OptionalGender interface.
type MockisEditCharacterRequest_OptionalGender struct {
	ctrl     *gomock.Controller
	recorder *MockisEditCharacterRequest_OptionalGenderMockRecorder
}

// MockisEditCharacterRequest_OptionalGenderMockRecorder is the mock recorder for MockisEditCharacterRequest_OptionalGender.
type MockisEditCharacterRequest_OptionalGenderMockRecorder struct {
	mock *MockisEditCharacterRequest_OptionalGender
}

// NewMockisEditCharacterRequest_OptionalGender creates a new mock instance.
func NewMockisEditCharacterRequest_OptionalGender(ctrl *gomock.Controller) *MockisEditCharacterRequest_OptionalGender {
	mock := &MockisEditCharacterRequest_OptionalGender{ctrl: ctrl}
	mock.recorder = &MockisEditCharacterRequest_OptionalGenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockisEditCharacterRequest_OptionalGender) EXPECT() *MockisEditCharacterRequest_OptionalGenderMockRecorder {
	return m.recorder
}

// isEditCharacterRequest_OptionalGender mocks base method.
func (m *MockisEditCharacterRequest_OptionalGender) isEditCharacterRequest_OptionalGender() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "isEditCharacterRequest_OptionalGender")
}

// isEditCharacterRequest_OptionalGender indicates an expected call of isEditCharacterRequest_OptionalGender.
func (mr *MockisEditCharacterRequest_OptionalGenderMockRecorder) isEditCharacterRequest_OptionalGender() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isEditCharacterRequest_OptionalGender", reflect.TypeOf((*MockisEditCharacterRequest_OptionalGender)(nil).isEditCharacterRequest_OptionalGender))
}

// MockisEditCharacterRequest_OptionalRealm is a mock of isEditCharacterRequest_OptionalRealm interface.
type MockisEditCharacterRequest_OptionalRealm struct {
	ctrl     *gomock.Controller
	recorder *MockisEditCharacterRequest_OptionalRealmMockRecorder
}

// MockisEditCharacterRequest_OptionalRealmMockRecorder is the mock recorder for MockisEditCharacterRequest_OptionalRealm.
type MockisEditCharacterRequest_OptionalRealmMockRecorder struct {
	mock *MockisEditCharacterRequest_OptionalRealm
}

// NewMockisEditCharacterRequest_OptionalRealm creates a new mock instance.
func NewMockisEditCharacterRequest_OptionalRealm(ctrl *gomock.Controller) *MockisEditCharacterRequest_OptionalRealm {
	mock := &MockisEditCharacterRequest_OptionalRealm{ctrl: ctrl}
	mock.recorder = &MockisEditCharacterRequest_OptionalRealmMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockisEditCharacterRequest_OptionalRealm) EXPECT() *MockisEditCharacterRequest_OptionalRealmMockRecorder {
	return m.recorder
}

// isEditCharacterRequest_OptionalRealm mocks base method.
func (m *MockisEditCharacterRequest_OptionalRealm) isEditCharacterRequest_OptionalRealm() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "isEditCharacterRequest_OptionalRealm")
}

// isEditCharacterRequest_OptionalRealm indicates an expected call of isEditCharacterRequest_OptionalRealm.
func (mr *MockisEditCharacterRequest_OptionalRealmMockRecorder) isEditCharacterRequest_OptionalRealm() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isEditCharacterRequest_OptionalRealm", reflect.TypeOf((*MockisEditCharacterRequest_OptionalRealm)(nil).isEditCharacterRequest_OptionalRealm))
}

// MockisEditCharacterRequest_OptionalPlayTime is a mock of isEditCharacterRequest_OptionalPlayTime interface.
type MockisEditCharacterRequest_OptionalPlayTime struct {
	ctrl     *gomock.Controller
	recorder *MockisEditCharacterRequest_OptionalPlayTimeMockRecorder
}

// MockisEditCharacterRequest_OptionalPlayTimeMockRecorder is the mock recorder for MockisEditCharacterRequest_OptionalPlayTime.
type MockisEditCharacterRequest_OptionalPlayTimeMockRecorder struct {
	mock *MockisEditCharacterRequest_OptionalPlayTime
}

// NewMockisEditCharacterRequest_OptionalPlayTime creates a new mock instance.
func NewMockisEditCharacterRequest_OptionalPlayTime(ctrl *gomock.Controller) *MockisEditCharacterRequest_OptionalPlayTime {
	mock := &MockisEditCharacterRequest_OptionalPlayTime{ctrl: ctrl}
	mock.recorder = &MockisEditCharacterRequest_OptionalPlayTimeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockisEditCharacterRequest_OptionalPlayTime) EXPECT() *MockisEditCharacterRequest_OptionalPlayTimeMockRecorder {
	return m.recorder
}

// isEditCharacterRequest_OptionalPlayTime mocks base method.
func (m *MockisEditCharacterRequest_OptionalPlayTime) isEditCharacterRequest_OptionalPlayTime() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "isEditCharacterRequest_OptionalPlayTime")
}

// isEditCharacterRequest_OptionalPlayTime indicates an expected call of isEditCharacterRequest_OptionalPlayTime.
func (mr *MockisEditCharacterRequest_OptionalPlayTimeMockRecorder) isEditCharacterRequest_OptionalPlayTime() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isEditCharacterRequest_OptionalPlayTime", reflect.TypeOf((*MockisEditCharacterRequest_OptionalPlayTime)(nil).isEditCharacterRequest_OptionalPlayTime))
}

// MockisEditCharacterRequest_OptionalLocation is a mock of isEditCharacterRequest_OptionalLocation interface.
type MockisEditCharacterRequest_OptionalLocation struct {
	ctrl     *gomock.Controller
	recorder *MockisEditCharacterRequest_OptionalLocationMockRecorder
}

// MockisEditCharacterRequest_OptionalLocationMockRecorder is the mock recorder for MockisEditCharacterRequest_OptionalLocation.
type MockisEditCharacterRequest_OptionalLocationMockRecorder struct {
	mock *MockisEditCharacterRequest_OptionalLocation
}

// NewMockisEditCharacterRequest_OptionalLocation creates a new mock instance.
func NewMockisEditCharacterRequest_OptionalLocation(ctrl *gomock.Controller) *MockisEditCharacterRequest_OptionalLocation {
	mock := &MockisEditCharacterRequest_OptionalLocation{ctrl: ctrl}
	mock.recorder = &MockisEditCharacterRequest_OptionalLocationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockisEditCharacterRequest_OptionalLocation) EXPECT() *MockisEditCharacterRequest_OptionalLocationMockRecorder {
	return m.recorder
}

// isEditCharacterRequest_OptionalLocation mocks base method.
func (m *MockisEditCharacterRequest_OptionalLocation) isEditCharacterRequest_OptionalLocation() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "isEditCharacterRequest_OptionalLocation")
}

// isEditCharacterRequest_OptionalLocation indicates an expected call of isEditCharacterRequest_OptionalLocation.
func (mr *MockisEditCharacterRequest_OptionalLocationMockRecorder) isEditCharacterRequest_OptionalLocation() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isEditCharacterRequest_OptionalLocation", reflect.TypeOf((*MockisEditCharacterRequest_OptionalLocation)(nil).isEditCharacterRequest_OptionalLocation))
}
