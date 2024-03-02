// Code generated by MockGen. DO NOT EDIT.
// Source: /home/wil/sro/git/go-backend/pkg/pb/servermanager_grpc.pb.go
//
// Generated by this command:
//
//	mockgen -package=mocks -source=/home/wil/sro/git/go-backend/pkg/pb/servermanager_grpc.pb.go -destination=/home/wil/sro/git/go-backend/pkg/mocks/servermanager_grpc.pb_mock.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	pb "github.com/ShatteredRealms/go-backend/pkg/pb"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// MockServerManagerServiceClient is a mock of ServerManagerServiceClient interface.
type MockServerManagerServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockServerManagerServiceClientMockRecorder
}

// MockServerManagerServiceClientMockRecorder is the mock recorder for MockServerManagerServiceClient.
type MockServerManagerServiceClientMockRecorder struct {
	mock *MockServerManagerServiceClient
}

// NewMockServerManagerServiceClient creates a new mock instance.
func NewMockServerManagerServiceClient(ctrl *gomock.Controller) *MockServerManagerServiceClient {
	mock := &MockServerManagerServiceClient{ctrl: ctrl}
	mock.recorder = &MockServerManagerServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServerManagerServiceClient) EXPECT() *MockServerManagerServiceClientMockRecorder {
	return m.recorder
}

// CreateDimension mocks base method.
func (m *MockServerManagerServiceClient) CreateDimension(ctx context.Context, in *pb.CreateDimensionRequest, opts ...grpc.CallOption) (*pb.Dimension, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateDimension", varargs...)
	ret0, _ := ret[0].(*pb.Dimension)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDimension indicates an expected call of CreateDimension.
func (mr *MockServerManagerServiceClientMockRecorder) CreateDimension(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDimension", reflect.TypeOf((*MockServerManagerServiceClient)(nil).CreateDimension), varargs...)
}

// CreateMap mocks base method.
func (m *MockServerManagerServiceClient) CreateMap(ctx context.Context, in *pb.CreateMapRequest, opts ...grpc.CallOption) (*pb.Map, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateMap", varargs...)
	ret0, _ := ret[0].(*pb.Map)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMap indicates an expected call of CreateMap.
func (mr *MockServerManagerServiceClientMockRecorder) CreateMap(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMap", reflect.TypeOf((*MockServerManagerServiceClient)(nil).CreateMap), varargs...)
}

// DeleteDimension mocks base method.
func (m *MockServerManagerServiceClient) DeleteDimension(ctx context.Context, in *pb.DimensionTarget, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteDimension", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteDimension indicates an expected call of DeleteDimension.
func (mr *MockServerManagerServiceClientMockRecorder) DeleteDimension(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDimension", reflect.TypeOf((*MockServerManagerServiceClient)(nil).DeleteDimension), varargs...)
}

// DeleteMap mocks base method.
func (m *MockServerManagerServiceClient) DeleteMap(ctx context.Context, in *pb.MapTarget, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteMap", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteMap indicates an expected call of DeleteMap.
func (mr *MockServerManagerServiceClientMockRecorder) DeleteMap(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMap", reflect.TypeOf((*MockServerManagerServiceClient)(nil).DeleteMap), varargs...)
}

// DuplicateDimension mocks base method.
func (m *MockServerManagerServiceClient) DuplicateDimension(ctx context.Context, in *pb.DuplicateDimensionRequest, opts ...grpc.CallOption) (*pb.Dimension, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DuplicateDimension", varargs...)
	ret0, _ := ret[0].(*pb.Dimension)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DuplicateDimension indicates an expected call of DuplicateDimension.
func (mr *MockServerManagerServiceClientMockRecorder) DuplicateDimension(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DuplicateDimension", reflect.TypeOf((*MockServerManagerServiceClient)(nil).DuplicateDimension), varargs...)
}

// EditDimension mocks base method.
func (m *MockServerManagerServiceClient) EditDimension(ctx context.Context, in *pb.EditDimensionRequest, opts ...grpc.CallOption) (*pb.Dimension, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "EditDimension", varargs...)
	ret0, _ := ret[0].(*pb.Dimension)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditDimension indicates an expected call of EditDimension.
func (mr *MockServerManagerServiceClientMockRecorder) EditDimension(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditDimension", reflect.TypeOf((*MockServerManagerServiceClient)(nil).EditDimension), varargs...)
}

// EditMap mocks base method.
func (m *MockServerManagerServiceClient) EditMap(ctx context.Context, in *pb.EditMapRequest, opts ...grpc.CallOption) (*pb.Map, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "EditMap", varargs...)
	ret0, _ := ret[0].(*pb.Map)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditMap indicates an expected call of EditMap.
func (mr *MockServerManagerServiceClientMockRecorder) EditMap(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditMap", reflect.TypeOf((*MockServerManagerServiceClient)(nil).EditMap), varargs...)
}

// GetAllDimension mocks base method.
func (m *MockServerManagerServiceClient) GetAllDimension(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb.Dimensions, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAllDimension", varargs...)
	ret0, _ := ret[0].(*pb.Dimensions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllDimension indicates an expected call of GetAllDimension.
func (mr *MockServerManagerServiceClientMockRecorder) GetAllDimension(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllDimension", reflect.TypeOf((*MockServerManagerServiceClient)(nil).GetAllDimension), varargs...)
}

// GetAllMaps mocks base method.
func (m *MockServerManagerServiceClient) GetAllMaps(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb.Maps, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAllMaps", varargs...)
	ret0, _ := ret[0].(*pb.Maps)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllMaps indicates an expected call of GetAllMaps.
func (mr *MockServerManagerServiceClientMockRecorder) GetAllMaps(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllMaps", reflect.TypeOf((*MockServerManagerServiceClient)(nil).GetAllMaps), varargs...)
}

// GetDimension mocks base method.
func (m *MockServerManagerServiceClient) GetDimension(ctx context.Context, in *pb.DimensionTarget, opts ...grpc.CallOption) (*pb.Dimension, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetDimension", varargs...)
	ret0, _ := ret[0].(*pb.Dimension)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDimension indicates an expected call of GetDimension.
func (mr *MockServerManagerServiceClientMockRecorder) GetDimension(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDimension", reflect.TypeOf((*MockServerManagerServiceClient)(nil).GetDimension), varargs...)
}

// GetMap mocks base method.
func (m *MockServerManagerServiceClient) GetMap(ctx context.Context, in *pb.MapTarget, opts ...grpc.CallOption) (*pb.Map, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetMap", varargs...)
	ret0, _ := ret[0].(*pb.Map)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMap indicates an expected call of GetMap.
func (mr *MockServerManagerServiceClientMockRecorder) GetMap(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMap", reflect.TypeOf((*MockServerManagerServiceClient)(nil).GetMap), varargs...)
}

// MockServerManagerServiceServer is a mock of ServerManagerServiceServer interface.
type MockServerManagerServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockServerManagerServiceServerMockRecorder
}

// MockServerManagerServiceServerMockRecorder is the mock recorder for MockServerManagerServiceServer.
type MockServerManagerServiceServerMockRecorder struct {
	mock *MockServerManagerServiceServer
}

// NewMockServerManagerServiceServer creates a new mock instance.
func NewMockServerManagerServiceServer(ctrl *gomock.Controller) *MockServerManagerServiceServer {
	mock := &MockServerManagerServiceServer{ctrl: ctrl}
	mock.recorder = &MockServerManagerServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServerManagerServiceServer) EXPECT() *MockServerManagerServiceServerMockRecorder {
	return m.recorder
}

// CreateDimension mocks base method.
func (m *MockServerManagerServiceServer) CreateDimension(arg0 context.Context, arg1 *pb.CreateDimensionRequest) (*pb.Dimension, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDimension", arg0, arg1)
	ret0, _ := ret[0].(*pb.Dimension)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDimension indicates an expected call of CreateDimension.
func (mr *MockServerManagerServiceServerMockRecorder) CreateDimension(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDimension", reflect.TypeOf((*MockServerManagerServiceServer)(nil).CreateDimension), arg0, arg1)
}

// CreateMap mocks base method.
func (m *MockServerManagerServiceServer) CreateMap(arg0 context.Context, arg1 *pb.CreateMapRequest) (*pb.Map, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMap", arg0, arg1)
	ret0, _ := ret[0].(*pb.Map)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMap indicates an expected call of CreateMap.
func (mr *MockServerManagerServiceServerMockRecorder) CreateMap(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMap", reflect.TypeOf((*MockServerManagerServiceServer)(nil).CreateMap), arg0, arg1)
}

// DeleteDimension mocks base method.
func (m *MockServerManagerServiceServer) DeleteDimension(arg0 context.Context, arg1 *pb.DimensionTarget) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteDimension", arg0, arg1)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteDimension indicates an expected call of DeleteDimension.
func (mr *MockServerManagerServiceServerMockRecorder) DeleteDimension(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDimension", reflect.TypeOf((*MockServerManagerServiceServer)(nil).DeleteDimension), arg0, arg1)
}

// DeleteMap mocks base method.
func (m *MockServerManagerServiceServer) DeleteMap(arg0 context.Context, arg1 *pb.MapTarget) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMap", arg0, arg1)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteMap indicates an expected call of DeleteMap.
func (mr *MockServerManagerServiceServerMockRecorder) DeleteMap(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMap", reflect.TypeOf((*MockServerManagerServiceServer)(nil).DeleteMap), arg0, arg1)
}

// DuplicateDimension mocks base method.
func (m *MockServerManagerServiceServer) DuplicateDimension(arg0 context.Context, arg1 *pb.DuplicateDimensionRequest) (*pb.Dimension, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DuplicateDimension", arg0, arg1)
	ret0, _ := ret[0].(*pb.Dimension)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DuplicateDimension indicates an expected call of DuplicateDimension.
func (mr *MockServerManagerServiceServerMockRecorder) DuplicateDimension(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DuplicateDimension", reflect.TypeOf((*MockServerManagerServiceServer)(nil).DuplicateDimension), arg0, arg1)
}

// EditDimension mocks base method.
func (m *MockServerManagerServiceServer) EditDimension(arg0 context.Context, arg1 *pb.EditDimensionRequest) (*pb.Dimension, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditDimension", arg0, arg1)
	ret0, _ := ret[0].(*pb.Dimension)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditDimension indicates an expected call of EditDimension.
func (mr *MockServerManagerServiceServerMockRecorder) EditDimension(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditDimension", reflect.TypeOf((*MockServerManagerServiceServer)(nil).EditDimension), arg0, arg1)
}

// EditMap mocks base method.
func (m *MockServerManagerServiceServer) EditMap(arg0 context.Context, arg1 *pb.EditMapRequest) (*pb.Map, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditMap", arg0, arg1)
	ret0, _ := ret[0].(*pb.Map)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditMap indicates an expected call of EditMap.
func (mr *MockServerManagerServiceServerMockRecorder) EditMap(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditMap", reflect.TypeOf((*MockServerManagerServiceServer)(nil).EditMap), arg0, arg1)
}

// GetAllDimension mocks base method.
func (m *MockServerManagerServiceServer) GetAllDimension(arg0 context.Context, arg1 *emptypb.Empty) (*pb.Dimensions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllDimension", arg0, arg1)
	ret0, _ := ret[0].(*pb.Dimensions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllDimension indicates an expected call of GetAllDimension.
func (mr *MockServerManagerServiceServerMockRecorder) GetAllDimension(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllDimension", reflect.TypeOf((*MockServerManagerServiceServer)(nil).GetAllDimension), arg0, arg1)
}

// GetAllMaps mocks base method.
func (m *MockServerManagerServiceServer) GetAllMaps(arg0 context.Context, arg1 *emptypb.Empty) (*pb.Maps, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllMaps", arg0, arg1)
	ret0, _ := ret[0].(*pb.Maps)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllMaps indicates an expected call of GetAllMaps.
func (mr *MockServerManagerServiceServerMockRecorder) GetAllMaps(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllMaps", reflect.TypeOf((*MockServerManagerServiceServer)(nil).GetAllMaps), arg0, arg1)
}

// GetDimension mocks base method.
func (m *MockServerManagerServiceServer) GetDimension(arg0 context.Context, arg1 *pb.DimensionTarget) (*pb.Dimension, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDimension", arg0, arg1)
	ret0, _ := ret[0].(*pb.Dimension)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDimension indicates an expected call of GetDimension.
func (mr *MockServerManagerServiceServerMockRecorder) GetDimension(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDimension", reflect.TypeOf((*MockServerManagerServiceServer)(nil).GetDimension), arg0, arg1)
}

// GetMap mocks base method.
func (m *MockServerManagerServiceServer) GetMap(arg0 context.Context, arg1 *pb.MapTarget) (*pb.Map, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMap", arg0, arg1)
	ret0, _ := ret[0].(*pb.Map)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMap indicates an expected call of GetMap.
func (mr *MockServerManagerServiceServerMockRecorder) GetMap(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMap", reflect.TypeOf((*MockServerManagerServiceServer)(nil).GetMap), arg0, arg1)
}

// mustEmbedUnimplementedServerManagerServiceServer mocks base method.
func (m *MockServerManagerServiceServer) mustEmbedUnimplementedServerManagerServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedServerManagerServiceServer")
}

// mustEmbedUnimplementedServerManagerServiceServer indicates an expected call of mustEmbedUnimplementedServerManagerServiceServer.
func (mr *MockServerManagerServiceServerMockRecorder) mustEmbedUnimplementedServerManagerServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedServerManagerServiceServer", reflect.TypeOf((*MockServerManagerServiceServer)(nil).mustEmbedUnimplementedServerManagerServiceServer))
}

// MockUnsafeServerManagerServiceServer is a mock of UnsafeServerManagerServiceServer interface.
type MockUnsafeServerManagerServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeServerManagerServiceServerMockRecorder
}

// MockUnsafeServerManagerServiceServerMockRecorder is the mock recorder for MockUnsafeServerManagerServiceServer.
type MockUnsafeServerManagerServiceServerMockRecorder struct {
	mock *MockUnsafeServerManagerServiceServer
}

// NewMockUnsafeServerManagerServiceServer creates a new mock instance.
func NewMockUnsafeServerManagerServiceServer(ctrl *gomock.Controller) *MockUnsafeServerManagerServiceServer {
	mock := &MockUnsafeServerManagerServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeServerManagerServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeServerManagerServiceServer) EXPECT() *MockUnsafeServerManagerServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedServerManagerServiceServer mocks base method.
func (m *MockUnsafeServerManagerServiceServer) mustEmbedUnimplementedServerManagerServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedServerManagerServiceServer")
}

// mustEmbedUnimplementedServerManagerServiceServer indicates an expected call of mustEmbedUnimplementedServerManagerServiceServer.
func (mr *MockUnsafeServerManagerServiceServerMockRecorder) mustEmbedUnimplementedServerManagerServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedServerManagerServiceServer", reflect.TypeOf((*MockUnsafeServerManagerServiceServer)(nil).mustEmbedUnimplementedServerManagerServiceServer))
}
