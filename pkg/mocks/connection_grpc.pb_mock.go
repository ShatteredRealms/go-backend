// Code generated by MockGen. DO NOT EDIT.
// Source: /home/wil/sro/git/go-backend/pkg/pb/connection_grpc.pb.go
//
// Generated by this command:
//
//	mockgen -package=mocks -source=/home/wil/sro/git/go-backend/pkg/pb/connection_grpc.pb.go -destination=/home/wil/sro/git/go-backend/pkg/mocks/connection_grpc.pb_mock.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	pb "github.com/ShatteredRealms/go-backend/pkg/pb"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockConnectionServiceClient is a mock of ConnectionServiceClient interface.
type MockConnectionServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockConnectionServiceClientMockRecorder
}

// MockConnectionServiceClientMockRecorder is the mock recorder for MockConnectionServiceClient.
type MockConnectionServiceClientMockRecorder struct {
	mock *MockConnectionServiceClient
}

// NewMockConnectionServiceClient creates a new mock instance.
func NewMockConnectionServiceClient(ctrl *gomock.Controller) *MockConnectionServiceClient {
	mock := &MockConnectionServiceClient{ctrl: ctrl}
	mock.recorder = &MockConnectionServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConnectionServiceClient) EXPECT() *MockConnectionServiceClientMockRecorder {
	return m.recorder
}

// ConnectGameServer mocks base method.
func (m *MockConnectionServiceClient) ConnectGameServer(ctx context.Context, in *pb.CharacterTarget, opts ...grpc.CallOption) (*pb.ConnectGameServerResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ConnectGameServer", varargs...)
	ret0, _ := ret[0].(*pb.ConnectGameServerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConnectGameServer indicates an expected call of ConnectGameServer.
func (mr *MockConnectionServiceClientMockRecorder) ConnectGameServer(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnectGameServer", reflect.TypeOf((*MockConnectionServiceClient)(nil).ConnectGameServer), varargs...)
}

// IsPlaying mocks base method.
func (m *MockConnectionServiceClient) IsPlaying(ctx context.Context, in *pb.CharacterTarget, opts ...grpc.CallOption) (*pb.ConnectionStatus, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "IsPlaying", varargs...)
	ret0, _ := ret[0].(*pb.ConnectionStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsPlaying indicates an expected call of IsPlaying.
func (mr *MockConnectionServiceClientMockRecorder) IsPlaying(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPlaying", reflect.TypeOf((*MockConnectionServiceClient)(nil).IsPlaying), varargs...)
}

// TransferPlayer mocks base method.
func (m *MockConnectionServiceClient) TransferPlayer(ctx context.Context, in *pb.TransferPlayerRequest, opts ...grpc.CallOption) (*pb.ConnectGameServerResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "TransferPlayer", varargs...)
	ret0, _ := ret[0].(*pb.ConnectGameServerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TransferPlayer indicates an expected call of TransferPlayer.
func (mr *MockConnectionServiceClientMockRecorder) TransferPlayer(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransferPlayer", reflect.TypeOf((*MockConnectionServiceClient)(nil).TransferPlayer), varargs...)
}

// VerifyConnect mocks base method.
func (m *MockConnectionServiceClient) VerifyConnect(ctx context.Context, in *pb.VerifyConnectRequest, opts ...grpc.CallOption) (*pb.CharacterDetails, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "VerifyConnect", varargs...)
	ret0, _ := ret[0].(*pb.CharacterDetails)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyConnect indicates an expected call of VerifyConnect.
func (mr *MockConnectionServiceClientMockRecorder) VerifyConnect(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyConnect", reflect.TypeOf((*MockConnectionServiceClient)(nil).VerifyConnect), varargs...)
}

// MockConnectionServiceServer is a mock of ConnectionServiceServer interface.
type MockConnectionServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockConnectionServiceServerMockRecorder
}

// MockConnectionServiceServerMockRecorder is the mock recorder for MockConnectionServiceServer.
type MockConnectionServiceServerMockRecorder struct {
	mock *MockConnectionServiceServer
}

// NewMockConnectionServiceServer creates a new mock instance.
func NewMockConnectionServiceServer(ctrl *gomock.Controller) *MockConnectionServiceServer {
	mock := &MockConnectionServiceServer{ctrl: ctrl}
	mock.recorder = &MockConnectionServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConnectionServiceServer) EXPECT() *MockConnectionServiceServerMockRecorder {
	return m.recorder
}

// ConnectGameServer mocks base method.
func (m *MockConnectionServiceServer) ConnectGameServer(arg0 context.Context, arg1 *pb.CharacterTarget) (*pb.ConnectGameServerResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConnectGameServer", arg0, arg1)
	ret0, _ := ret[0].(*pb.ConnectGameServerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConnectGameServer indicates an expected call of ConnectGameServer.
func (mr *MockConnectionServiceServerMockRecorder) ConnectGameServer(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnectGameServer", reflect.TypeOf((*MockConnectionServiceServer)(nil).ConnectGameServer), arg0, arg1)
}

// IsPlaying mocks base method.
func (m *MockConnectionServiceServer) IsPlaying(arg0 context.Context, arg1 *pb.CharacterTarget) (*pb.ConnectionStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsPlaying", arg0, arg1)
	ret0, _ := ret[0].(*pb.ConnectionStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsPlaying indicates an expected call of IsPlaying.
func (mr *MockConnectionServiceServerMockRecorder) IsPlaying(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPlaying", reflect.TypeOf((*MockConnectionServiceServer)(nil).IsPlaying), arg0, arg1)
}

// TransferPlayer mocks base method.
func (m *MockConnectionServiceServer) TransferPlayer(arg0 context.Context, arg1 *pb.TransferPlayerRequest) (*pb.ConnectGameServerResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransferPlayer", arg0, arg1)
	ret0, _ := ret[0].(*pb.ConnectGameServerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TransferPlayer indicates an expected call of TransferPlayer.
func (mr *MockConnectionServiceServerMockRecorder) TransferPlayer(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransferPlayer", reflect.TypeOf((*MockConnectionServiceServer)(nil).TransferPlayer), arg0, arg1)
}

// VerifyConnect mocks base method.
func (m *MockConnectionServiceServer) VerifyConnect(arg0 context.Context, arg1 *pb.VerifyConnectRequest) (*pb.CharacterDetails, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyConnect", arg0, arg1)
	ret0, _ := ret[0].(*pb.CharacterDetails)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyConnect indicates an expected call of VerifyConnect.
func (mr *MockConnectionServiceServerMockRecorder) VerifyConnect(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyConnect", reflect.TypeOf((*MockConnectionServiceServer)(nil).VerifyConnect), arg0, arg1)
}

// mustEmbedUnimplementedConnectionServiceServer mocks base method.
func (m *MockConnectionServiceServer) mustEmbedUnimplementedConnectionServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedConnectionServiceServer")
}

// mustEmbedUnimplementedConnectionServiceServer indicates an expected call of mustEmbedUnimplementedConnectionServiceServer.
func (mr *MockConnectionServiceServerMockRecorder) mustEmbedUnimplementedConnectionServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedConnectionServiceServer", reflect.TypeOf((*MockConnectionServiceServer)(nil).mustEmbedUnimplementedConnectionServiceServer))
}

// MockUnsafeConnectionServiceServer is a mock of UnsafeConnectionServiceServer interface.
type MockUnsafeConnectionServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeConnectionServiceServerMockRecorder
}

// MockUnsafeConnectionServiceServerMockRecorder is the mock recorder for MockUnsafeConnectionServiceServer.
type MockUnsafeConnectionServiceServerMockRecorder struct {
	mock *MockUnsafeConnectionServiceServer
}

// NewMockUnsafeConnectionServiceServer creates a new mock instance.
func NewMockUnsafeConnectionServiceServer(ctrl *gomock.Controller) *MockUnsafeConnectionServiceServer {
	mock := &MockUnsafeConnectionServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeConnectionServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeConnectionServiceServer) EXPECT() *MockUnsafeConnectionServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedConnectionServiceServer mocks base method.
func (m *MockUnsafeConnectionServiceServer) mustEmbedUnimplementedConnectionServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedConnectionServiceServer")
}

// mustEmbedUnimplementedConnectionServiceServer indicates an expected call of mustEmbedUnimplementedConnectionServiceServer.
func (mr *MockUnsafeConnectionServiceServerMockRecorder) mustEmbedUnimplementedConnectionServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedConnectionServiceServer", reflect.TypeOf((*MockUnsafeConnectionServiceServer)(nil).mustEmbedUnimplementedConnectionServiceServer))
}
