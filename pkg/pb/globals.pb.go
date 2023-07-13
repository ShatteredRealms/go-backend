// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.15.8
// source: sro/globals.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type UserTarget struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Target:
	//
	//	*UserTarget_Id
	//	*UserTarget_Username
	Target isUserTarget_Target `protobuf_oneof:"target"`
}

func (x *UserTarget) Reset() {
	*x = UserTarget{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sro_globals_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserTarget) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserTarget) ProtoMessage() {}

func (x *UserTarget) ProtoReflect() protoreflect.Message {
	mi := &file_sro_globals_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserTarget.ProtoReflect.Descriptor instead.
func (*UserTarget) Descriptor() ([]byte, []int) {
	return file_sro_globals_proto_rawDescGZIP(), []int{0}
}

func (m *UserTarget) GetTarget() isUserTarget_Target {
	if m != nil {
		return m.Target
	}
	return nil
}

func (x *UserTarget) GetId() string {
	if x, ok := x.GetTarget().(*UserTarget_Id); ok {
		return x.Id
	}
	return ""
}

func (x *UserTarget) GetUsername() string {
	if x, ok := x.GetTarget().(*UserTarget_Username); ok {
		return x.Username
	}
	return ""
}

type isUserTarget_Target interface {
	isUserTarget_Target()
}

type UserTarget_Id struct {
	Id string `protobuf:"bytes,1,opt,name=id,proto3,oneof"`
}

type UserTarget_Username struct {
	Username string `protobuf:"bytes,2,opt,name=username,proto3,oneof"`
}

func (*UserTarget_Id) isUserTarget_Target() {}

func (*UserTarget_Username) isUserTarget_Target() {}

type Location struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	World string  `protobuf:"bytes,1,opt,name=world,proto3" json:"world,omitempty"`
	X     float32 `protobuf:"fixed32,2,opt,name=x,proto3" json:"x,omitempty"`
	Y     float32 `protobuf:"fixed32,3,opt,name=y,proto3" json:"y,omitempty"`
	Z     float32 `protobuf:"fixed32,4,opt,name=z,proto3" json:"z,omitempty"`
}

func (x *Location) Reset() {
	*x = Location{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sro_globals_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Location) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Location) ProtoMessage() {}

func (x *Location) ProtoReflect() protoreflect.Message {
	mi := &file_sro_globals_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Location.ProtoReflect.Descriptor instead.
func (*Location) Descriptor() ([]byte, []int) {
	return file_sro_globals_proto_rawDescGZIP(), []int{1}
}

func (x *Location) GetWorld() string {
	if x != nil {
		return x.World
	}
	return ""
}

func (x *Location) GetX() float32 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *Location) GetY() float32 {
	if x != nil {
		return x.Y
	}
	return 0
}

func (x *Location) GetZ() float32 {
	if x != nil {
		return x.Z
	}
	return 0
}

var File_sro_globals_proto protoreflect.FileDescriptor

var file_sro_globals_proto_rawDesc = []byte{
	0x0a, 0x11, 0x73, 0x72, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x03, 0x73, 0x72, 0x6f, 0x22, 0x46, 0x0a, 0x0a, 0x55, 0x73, 0x65, 0x72,
	0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12, 0x10, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x00, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1c, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x08, 0x75, 0x73,
	0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x08, 0x0a, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74,
	0x22, 0x4a, 0x0a, 0x08, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05,
	0x77, 0x6f, 0x72, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x77, 0x6f, 0x72,
	0x6c, 0x64, 0x12, 0x0c, 0x0a, 0x01, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x01, 0x78,
	0x12, 0x0c, 0x0a, 0x01, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52, 0x01, 0x79, 0x12, 0x0c,
	0x0a, 0x01, 0x7a, 0x18, 0x04, 0x20, 0x01, 0x28, 0x02, 0x52, 0x01, 0x7a, 0x42, 0x08, 0x5a, 0x06,
	0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_sro_globals_proto_rawDescOnce sync.Once
	file_sro_globals_proto_rawDescData = file_sro_globals_proto_rawDesc
)

func file_sro_globals_proto_rawDescGZIP() []byte {
	file_sro_globals_proto_rawDescOnce.Do(func() {
		file_sro_globals_proto_rawDescData = protoimpl.X.CompressGZIP(file_sro_globals_proto_rawDescData)
	})
	return file_sro_globals_proto_rawDescData
}

var file_sro_globals_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_sro_globals_proto_goTypes = []interface{}{
	(*UserTarget)(nil), // 0: sro.UserTarget
	(*Location)(nil),   // 1: sro.Location
}
var file_sro_globals_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_sro_globals_proto_init() }
func file_sro_globals_proto_init() {
	if File_sro_globals_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_sro_globals_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserTarget); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_sro_globals_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Location); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_sro_globals_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*UserTarget_Id)(nil),
		(*UserTarget_Username)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_sro_globals_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_sro_globals_proto_goTypes,
		DependencyIndexes: file_sro_globals_proto_depIdxs,
		MessageInfos:      file_sro_globals_proto_msgTypes,
	}.Build()
	File_sro_globals_proto = out.File
	file_sro_globals_proto_rawDesc = nil
	file_sro_globals_proto_goTypes = nil
	file_sro_globals_proto_depIdxs = nil
}
