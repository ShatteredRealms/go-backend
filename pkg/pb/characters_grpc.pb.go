// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.15.8
// source: sro/characters/characters.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	CharactersService_GetGenders_FullMethodName              = "/sro.characters.CharactersService/GetGenders"
	CharactersService_GetRealms_FullMethodName               = "/sro.characters.CharactersService/GetRealms"
	CharactersService_GetCharacters_FullMethodName           = "/sro.characters.CharactersService/GetCharacters"
	CharactersService_GetCharacter_FullMethodName            = "/sro.characters.CharactersService/GetCharacter"
	CharactersService_CreateCharacter_FullMethodName         = "/sro.characters.CharactersService/CreateCharacter"
	CharactersService_DeleteCharacter_FullMethodName         = "/sro.characters.CharactersService/DeleteCharacter"
	CharactersService_GetAllCharactersForUser_FullMethodName = "/sro.characters.CharactersService/GetAllCharactersForUser"
	CharactersService_EditCharacter_FullMethodName           = "/sro.characters.CharactersService/EditCharacter"
	CharactersService_AddCharacterPlayTime_FullMethodName    = "/sro.characters.CharactersService/AddCharacterPlayTime"
)

// CharactersServiceClient is the client API for CharactersService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CharactersServiceClient interface {
	GetGenders(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Genders, error)
	GetRealms(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Realms, error)
	GetCharacters(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*CharactersResponse, error)
	GetCharacter(ctx context.Context, in *CharacterTarget, opts ...grpc.CallOption) (*CharacterResponse, error)
	CreateCharacter(ctx context.Context, in *CreateCharacterRequest, opts ...grpc.CallOption) (*CharacterResponse, error)
	DeleteCharacter(ctx context.Context, in *CharacterTarget, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetAllCharactersForUser(ctx context.Context, in *UserTarget, opts ...grpc.CallOption) (*CharactersResponse, error)
	EditCharacter(ctx context.Context, in *EditCharacterRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Adds the given amount of playtime to the character and returns the total
	// playtime
	AddCharacterPlayTime(ctx context.Context, in *AddPlayTimeRequest, opts ...grpc.CallOption) (*PlayTimeResponse, error)
}

type charactersServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCharactersServiceClient(cc grpc.ClientConnInterface) CharactersServiceClient {
	return &charactersServiceClient{cc}
}

func (c *charactersServiceClient) GetGenders(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Genders, error) {
	out := new(Genders)
	err := c.cc.Invoke(ctx, CharactersService_GetGenders_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *charactersServiceClient) GetRealms(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Realms, error) {
	out := new(Realms)
	err := c.cc.Invoke(ctx, CharactersService_GetRealms_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *charactersServiceClient) GetCharacters(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*CharactersResponse, error) {
	out := new(CharactersResponse)
	err := c.cc.Invoke(ctx, CharactersService_GetCharacters_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *charactersServiceClient) GetCharacter(ctx context.Context, in *CharacterTarget, opts ...grpc.CallOption) (*CharacterResponse, error) {
	out := new(CharacterResponse)
	err := c.cc.Invoke(ctx, CharactersService_GetCharacter_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *charactersServiceClient) CreateCharacter(ctx context.Context, in *CreateCharacterRequest, opts ...grpc.CallOption) (*CharacterResponse, error) {
	out := new(CharacterResponse)
	err := c.cc.Invoke(ctx, CharactersService_CreateCharacter_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *charactersServiceClient) DeleteCharacter(ctx context.Context, in *CharacterTarget, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, CharactersService_DeleteCharacter_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *charactersServiceClient) GetAllCharactersForUser(ctx context.Context, in *UserTarget, opts ...grpc.CallOption) (*CharactersResponse, error) {
	out := new(CharactersResponse)
	err := c.cc.Invoke(ctx, CharactersService_GetAllCharactersForUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *charactersServiceClient) EditCharacter(ctx context.Context, in *EditCharacterRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, CharactersService_EditCharacter_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *charactersServiceClient) AddCharacterPlayTime(ctx context.Context, in *AddPlayTimeRequest, opts ...grpc.CallOption) (*PlayTimeResponse, error) {
	out := new(PlayTimeResponse)
	err := c.cc.Invoke(ctx, CharactersService_AddCharacterPlayTime_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CharactersServiceServer is the server API for CharactersService service.
// All implementations must embed UnimplementedCharactersServiceServer
// for forward compatibility
type CharactersServiceServer interface {
	GetGenders(context.Context, *emptypb.Empty) (*Genders, error)
	GetRealms(context.Context, *emptypb.Empty) (*Realms, error)
	GetCharacters(context.Context, *emptypb.Empty) (*CharactersResponse, error)
	GetCharacter(context.Context, *CharacterTarget) (*CharacterResponse, error)
	CreateCharacter(context.Context, *CreateCharacterRequest) (*CharacterResponse, error)
	DeleteCharacter(context.Context, *CharacterTarget) (*emptypb.Empty, error)
	GetAllCharactersForUser(context.Context, *UserTarget) (*CharactersResponse, error)
	EditCharacter(context.Context, *EditCharacterRequest) (*emptypb.Empty, error)
	// Adds the given amount of playtime to the character and returns the total
	// playtime
	AddCharacterPlayTime(context.Context, *AddPlayTimeRequest) (*PlayTimeResponse, error)
	mustEmbedUnimplementedCharactersServiceServer()
}

// UnimplementedCharactersServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCharactersServiceServer struct {
}

func (UnimplementedCharactersServiceServer) GetGenders(context.Context, *emptypb.Empty) (*Genders, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGenders not implemented")
}
func (UnimplementedCharactersServiceServer) GetRealms(context.Context, *emptypb.Empty) (*Realms, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRealms not implemented")
}
func (UnimplementedCharactersServiceServer) GetCharacters(context.Context, *emptypb.Empty) (*CharactersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCharacters not implemented")
}
func (UnimplementedCharactersServiceServer) GetCharacter(context.Context, *CharacterTarget) (*CharacterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCharacter not implemented")
}
func (UnimplementedCharactersServiceServer) CreateCharacter(context.Context, *CreateCharacterRequest) (*CharacterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCharacter not implemented")
}
func (UnimplementedCharactersServiceServer) DeleteCharacter(context.Context, *CharacterTarget) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCharacter not implemented")
}
func (UnimplementedCharactersServiceServer) GetAllCharactersForUser(context.Context, *UserTarget) (*CharactersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllCharactersForUser not implemented")
}
func (UnimplementedCharactersServiceServer) EditCharacter(context.Context, *EditCharacterRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditCharacter not implemented")
}
func (UnimplementedCharactersServiceServer) AddCharacterPlayTime(context.Context, *AddPlayTimeRequest) (*PlayTimeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddCharacterPlayTime not implemented")
}
func (UnimplementedCharactersServiceServer) mustEmbedUnimplementedCharactersServiceServer() {}

// UnsafeCharactersServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CharactersServiceServer will
// result in compilation errors.
type UnsafeCharactersServiceServer interface {
	mustEmbedUnimplementedCharactersServiceServer()
}

func RegisterCharactersServiceServer(s grpc.ServiceRegistrar, srv CharactersServiceServer) {
	s.RegisterService(&CharactersService_ServiceDesc, srv)
}

func _CharactersService_GetGenders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CharactersServiceServer).GetGenders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CharactersService_GetGenders_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CharactersServiceServer).GetGenders(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _CharactersService_GetRealms_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CharactersServiceServer).GetRealms(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CharactersService_GetRealms_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CharactersServiceServer).GetRealms(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _CharactersService_GetCharacters_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CharactersServiceServer).GetCharacters(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CharactersService_GetCharacters_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CharactersServiceServer).GetCharacters(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _CharactersService_GetCharacter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CharacterTarget)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CharactersServiceServer).GetCharacter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CharactersService_GetCharacter_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CharactersServiceServer).GetCharacter(ctx, req.(*CharacterTarget))
	}
	return interceptor(ctx, in, info, handler)
}

func _CharactersService_CreateCharacter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCharacterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CharactersServiceServer).CreateCharacter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CharactersService_CreateCharacter_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CharactersServiceServer).CreateCharacter(ctx, req.(*CreateCharacterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CharactersService_DeleteCharacter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CharacterTarget)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CharactersServiceServer).DeleteCharacter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CharactersService_DeleteCharacter_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CharactersServiceServer).DeleteCharacter(ctx, req.(*CharacterTarget))
	}
	return interceptor(ctx, in, info, handler)
}

func _CharactersService_GetAllCharactersForUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserTarget)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CharactersServiceServer).GetAllCharactersForUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CharactersService_GetAllCharactersForUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CharactersServiceServer).GetAllCharactersForUser(ctx, req.(*UserTarget))
	}
	return interceptor(ctx, in, info, handler)
}

func _CharactersService_EditCharacter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditCharacterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CharactersServiceServer).EditCharacter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CharactersService_EditCharacter_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CharactersServiceServer).EditCharacter(ctx, req.(*EditCharacterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CharactersService_AddCharacterPlayTime_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddPlayTimeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CharactersServiceServer).AddCharacterPlayTime(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CharactersService_AddCharacterPlayTime_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CharactersServiceServer).AddCharacterPlayTime(ctx, req.(*AddPlayTimeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CharactersService_ServiceDesc is the grpc.ServiceDesc for CharactersService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CharactersService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sro.characters.CharactersService",
	HandlerType: (*CharactersServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetGenders",
			Handler:    _CharactersService_GetGenders_Handler,
		},
		{
			MethodName: "GetRealms",
			Handler:    _CharactersService_GetRealms_Handler,
		},
		{
			MethodName: "GetCharacters",
			Handler:    _CharactersService_GetCharacters_Handler,
		},
		{
			MethodName: "GetCharacter",
			Handler:    _CharactersService_GetCharacter_Handler,
		},
		{
			MethodName: "CreateCharacter",
			Handler:    _CharactersService_CreateCharacter_Handler,
		},
		{
			MethodName: "DeleteCharacter",
			Handler:    _CharactersService_DeleteCharacter_Handler,
		},
		{
			MethodName: "GetAllCharactersForUser",
			Handler:    _CharactersService_GetAllCharactersForUser_Handler,
		},
		{
			MethodName: "EditCharacter",
			Handler:    _CharactersService_EditCharacter_Handler,
		},
		{
			MethodName: "AddCharacterPlayTime",
			Handler:    _CharactersService_AddCharacterPlayTime_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sro/characters/characters.proto",
}
