// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: characters.proto

package pb

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PlayTimeMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CharacterId uint64 `protobuf:"varint,1,opt,name=character_id,json=characterId,proto3" json:"character_id,omitempty"`
	Time        uint64 `protobuf:"varint,2,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *PlayTimeMessage) Reset() {
	*x = PlayTimeMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_characters_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlayTimeMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlayTimeMessage) ProtoMessage() {}

func (x *PlayTimeMessage) ProtoReflect() protoreflect.Message {
	mi := &file_characters_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlayTimeMessage.ProtoReflect.Descriptor instead.
func (*PlayTimeMessage) Descriptor() ([]byte, []int) {
	return file_characters_proto_rawDescGZIP(), []int{0}
}

func (x *PlayTimeMessage) GetCharacterId() uint64 {
	if x != nil {
		return x.CharacterId
	}
	return 0
}

func (x *PlayTimeMessage) GetTime() uint64 {
	if x != nil {
		return x.Time
	}
	return 0
}

type DeleteCharacterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CharacterId uint64 `protobuf:"varint,1,opt,name=character_id,json=characterId,proto3" json:"character_id,omitempty"`
}

func (x *DeleteCharacterRequest) Reset() {
	*x = DeleteCharacterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_characters_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteCharacterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteCharacterRequest) ProtoMessage() {}

func (x *DeleteCharacterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_characters_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteCharacterRequest.ProtoReflect.Descriptor instead.
func (*DeleteCharacterRequest) Descriptor() ([]byte, []int) {
	return file_characters_proto_rawDescGZIP(), []int{1}
}

func (x *DeleteCharacterRequest) GetCharacterId() uint64 {
	if x != nil {
		return x.CharacterId
	}
	return 0
}

type CreateCharacterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId uint64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Name   string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Gender uint64 `protobuf:"varint,3,opt,name=gender,proto3" json:"gender,omitempty"`
	Realm  uint64 `protobuf:"varint,4,opt,name=realm,proto3" json:"realm,omitempty"`
}

func (x *CreateCharacterRequest) Reset() {
	*x = CreateCharacterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_characters_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateCharacterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateCharacterRequest) ProtoMessage() {}

func (x *CreateCharacterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_characters_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateCharacterRequest.ProtoReflect.Descriptor instead.
func (*CreateCharacterRequest) Descriptor() ([]byte, []int) {
	return file_characters_proto_rawDescGZIP(), []int{2}
}

func (x *CreateCharacterRequest) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *CreateCharacterRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateCharacterRequest) GetGender() uint64 {
	if x != nil {
		return x.Gender
	}
	return 0
}

func (x *CreateCharacterRequest) GetRealm() uint64 {
	if x != nil {
		return x.Realm
	}
	return 0
}

type UserTarget struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId uint64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *UserTarget) Reset() {
	*x = UserTarget{}
	if protoimpl.UnsafeEnabled {
		mi := &file_characters_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserTarget) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserTarget) ProtoMessage() {}

func (x *UserTarget) ProtoReflect() protoreflect.Message {
	mi := &file_characters_proto_msgTypes[3]
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
	return file_characters_proto_rawDescGZIP(), []int{3}
}

func (x *UserTarget) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type CharacterTarget struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CharacterId uint64 `protobuf:"varint,1,opt,name=character_id,json=characterId,proto3" json:"character_id,omitempty"`
}

func (x *CharacterTarget) Reset() {
	*x = CharacterTarget{}
	if protoimpl.UnsafeEnabled {
		mi := &file_characters_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CharacterTarget) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CharacterTarget) ProtoMessage() {}

func (x *CharacterTarget) ProtoReflect() protoreflect.Message {
	mi := &file_characters_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CharacterTarget.ProtoReflect.Descriptor instead.
func (*CharacterTarget) Descriptor() ([]byte, []int) {
	return file_characters_proto_rawDescGZIP(), []int{4}
}

func (x *CharacterTarget) GetCharacterId() uint64 {
	if x != nil {
		return x.CharacterId
	}
	return 0
}

type Character struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// The user account that owns the character
	Owner  *wrapperspb.UInt64Value `protobuf:"bytes,2,opt,name=owner,proto3" json:"owner,omitempty"`
	Name   *wrapperspb.StringValue `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Gender *wrapperspb.UInt64Value `protobuf:"bytes,4,opt,name=gender,proto3" json:"gender,omitempty"`
	Realm  *wrapperspb.UInt64Value `protobuf:"bytes,5,opt,name=realm,proto3" json:"realm,omitempty"`
	// Total play time in minutes
	PlayTime *wrapperspb.UInt64Value `protobuf:"bytes,7,opt,name=play_time,json=playTime,proto3" json:"play_time,omitempty"`
	Location *Location               `protobuf:"bytes,8,opt,name=location,proto3" json:"location,omitempty"`
}

func (x *Character) Reset() {
	*x = Character{}
	if protoimpl.UnsafeEnabled {
		mi := &file_characters_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Character) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Character) ProtoMessage() {}

func (x *Character) ProtoReflect() protoreflect.Message {
	mi := &file_characters_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Character.ProtoReflect.Descriptor instead.
func (*Character) Descriptor() ([]byte, []int) {
	return file_characters_proto_rawDescGZIP(), []int{5}
}

func (x *Character) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Character) GetOwner() *wrapperspb.UInt64Value {
	if x != nil {
		return x.Owner
	}
	return nil
}

func (x *Character) GetName() *wrapperspb.StringValue {
	if x != nil {
		return x.Name
	}
	return nil
}

func (x *Character) GetGender() *wrapperspb.UInt64Value {
	if x != nil {
		return x.Gender
	}
	return nil
}

func (x *Character) GetRealm() *wrapperspb.UInt64Value {
	if x != nil {
		return x.Realm
	}
	return nil
}

func (x *Character) GetPlayTime() *wrapperspb.UInt64Value {
	if x != nil {
		return x.PlayTime
	}
	return nil
}

func (x *Character) GetLocation() *Location {
	if x != nil {
		return x.Location
	}
	return nil
}

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
		mi := &file_characters_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Location) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Location) ProtoMessage() {}

func (x *Location) ProtoReflect() protoreflect.Message {
	mi := &file_characters_proto_msgTypes[6]
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
	return file_characters_proto_rawDescGZIP(), []int{6}
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

type Characters struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Characters []*Character `protobuf:"bytes,1,rep,name=characters,proto3" json:"characters,omitempty"`
}

func (x *Characters) Reset() {
	*x = Characters{}
	if protoimpl.UnsafeEnabled {
		mi := &file_characters_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Characters) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Characters) ProtoMessage() {}

func (x *Characters) ProtoReflect() protoreflect.Message {
	mi := &file_characters_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Characters.ProtoReflect.Descriptor instead.
func (*Characters) Descriptor() ([]byte, []int) {
	return file_characters_proto_rawDescGZIP(), []int{7}
}

func (x *Characters) GetCharacters() []*Character {
	if x != nil {
		return x.Characters
	}
	return nil
}

type Gender struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Gender) Reset() {
	*x = Gender{}
	if protoimpl.UnsafeEnabled {
		mi := &file_characters_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Gender) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Gender) ProtoMessage() {}

func (x *Gender) ProtoReflect() protoreflect.Message {
	mi := &file_characters_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Gender.ProtoReflect.Descriptor instead.
func (*Gender) Descriptor() ([]byte, []int) {
	return file_characters_proto_rawDescGZIP(), []int{8}
}

func (x *Gender) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Gender) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type Realm struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Realm) Reset() {
	*x = Realm{}
	if protoimpl.UnsafeEnabled {
		mi := &file_characters_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Realm) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Realm) ProtoMessage() {}

func (x *Realm) ProtoReflect() protoreflect.Message {
	mi := &file_characters_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Realm.ProtoReflect.Descriptor instead.
func (*Realm) Descriptor() ([]byte, []int) {
	return file_characters_proto_rawDescGZIP(), []int{9}
}

func (x *Realm) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Realm) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type Genders struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Genders []*Gender `protobuf:"bytes,1,rep,name=genders,proto3" json:"genders,omitempty"`
}

func (x *Genders) Reset() {
	*x = Genders{}
	if protoimpl.UnsafeEnabled {
		mi := &file_characters_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Genders) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Genders) ProtoMessage() {}

func (x *Genders) ProtoReflect() protoreflect.Message {
	mi := &file_characters_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Genders.ProtoReflect.Descriptor instead.
func (*Genders) Descriptor() ([]byte, []int) {
	return file_characters_proto_rawDescGZIP(), []int{10}
}

func (x *Genders) GetGenders() []*Gender {
	if x != nil {
		return x.Genders
	}
	return nil
}

type Realms struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Realms []*Realm `protobuf:"bytes,1,rep,name=realms,proto3" json:"realms,omitempty"`
}

func (x *Realms) Reset() {
	*x = Realms{}
	if protoimpl.UnsafeEnabled {
		mi := &file_characters_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Realms) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Realms) ProtoMessage() {}

func (x *Realms) ProtoReflect() protoreflect.Message {
	mi := &file_characters_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Realms.ProtoReflect.Descriptor instead.
func (*Realms) Descriptor() ([]byte, []int) {
	return file_characters_proto_rawDescGZIP(), []int{11}
}

func (x *Realms) GetRealms() []*Realm {
	if x != nil {
		return x.Realms
	}
	return nil
}

var File_characters_proto protoreflect.FileDescriptor

var file_characters_proto_rawDesc = []byte{
	0x0a, 0x10, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0e, 0x73, 0x72, 0x6f, 0x2e, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65,
	0x72, 0x73, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61,
	0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77,
	0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x48, 0x0a,
	0x0f, 0x50, 0x6c, 0x61, 0x79, 0x54, 0x69, 0x6d, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x21, 0x0a, 0x0c, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65,
	0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x22, 0x3b, 0x0a, 0x16, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74,
	0x65, 0x72, 0x49, 0x64, 0x22, 0x73, 0x0a, 0x16, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x68,
	0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17,
	0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x67,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x67, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x05, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x22, 0x25, 0x0a, 0x0a, 0x55, 0x73, 0x65,
	0x72, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x22, 0x34, 0x0a, 0x0f, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x54, 0x61, 0x72,
	0x67, 0x65, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x63, 0x68, 0x61, 0x72, 0x61,
	0x63, 0x74, 0x65, 0x72, 0x49, 0x64, 0x22, 0xdc, 0x02, 0x0a, 0x09, 0x43, 0x68, 0x61, 0x72, 0x61,
	0x63, 0x74, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x32, 0x0a, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x49, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x52, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x30, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x34, 0x0a, 0x06, 0x67, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x49, 0x6e,
	0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x12, 0x32, 0x0a, 0x05, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x55, 0x49, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x72,
	0x65, 0x61, 0x6c, 0x6d, 0x12, 0x39, 0x0a, 0x09, 0x70, 0x6c, 0x61, 0x79, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x49, 0x6e, 0x74, 0x36, 0x34,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x54, 0x69, 0x6d, 0x65, 0x12,
	0x34, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x18, 0x2e, 0x73, 0x72, 0x6f, 0x2e, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65,
	0x72, 0x73, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x6c, 0x6f, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x4a, 0x0a, 0x08, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x14, 0x0a, 0x05, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x12, 0x0c, 0x0a, 0x01, 0x78, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x02, 0x52, 0x01, 0x78, 0x12, 0x0c, 0x0a, 0x01, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02,
	0x52, 0x01, 0x79, 0x12, 0x0c, 0x0a, 0x01, 0x7a, 0x18, 0x04, 0x20, 0x01, 0x28, 0x02, 0x52, 0x01,
	0x7a, 0x22, 0x47, 0x0a, 0x0a, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x12,
	0x39, 0x0a, 0x0a, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x73, 0x72, 0x6f, 0x2e, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63,
	0x74, 0x65, 0x72, 0x73, 0x2e, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x52, 0x0a,
	0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x22, 0x2c, 0x0a, 0x06, 0x47, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x2b, 0x0a, 0x05, 0x52, 0x65, 0x61, 0x6c,
	0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x3b, 0x0a, 0x07, 0x47, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x73,
	0x12, 0x30, 0x0a, 0x07, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x16, 0x2e, 0x73, 0x72, 0x6f, 0x2e, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65,
	0x72, 0x73, 0x2e, 0x47, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x52, 0x07, 0x67, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x73, 0x22, 0x37, 0x0a, 0x06, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x73, 0x12, 0x2d, 0x0a, 0x06,
	0x72, 0x65, 0x61, 0x6c, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x73,
	0x72, 0x6f, 0x2e, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x52, 0x65,
	0x61, 0x6c, 0x6d, 0x52, 0x06, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x73, 0x32, 0xe5, 0x07, 0x0a, 0x11,
	0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x55, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x47, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x17, 0x2e, 0x73, 0x72, 0x6f,
	0x2e, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x47, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x73, 0x22, 0x13, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0d, 0x12, 0x0b, 0x2f, 0x76, 0x31,
	0x2f, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x73, 0x12, 0x52, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x41,
	0x6c, 0x6c, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x1a, 0x16, 0x2e, 0x73, 0x72, 0x6f, 0x2e, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72,
	0x73, 0x2e, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x73, 0x22, 0x12, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0c,
	0x12, 0x0a, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x73, 0x12, 0x5e, 0x0a, 0x10,
	0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1a, 0x2e, 0x73, 0x72, 0x6f, 0x2e, 0x63,
	0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63,
	0x74, 0x65, 0x72, 0x73, 0x22, 0x16, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10, 0x12, 0x0e, 0x2f, 0x76,
	0x31, 0x2f, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x12, 0x79, 0x0a, 0x17,
	0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73,
	0x46, 0x6f, 0x72, 0x55, 0x73, 0x65, 0x72, 0x12, 0x1a, 0x2e, 0x73, 0x72, 0x6f, 0x2e, 0x63, 0x68,
	0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x54, 0x61, 0x72,
	0x67, 0x65, 0x74, 0x1a, 0x1a, 0x2e, 0x73, 0x72, 0x6f, 0x2e, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63,
	0x74, 0x65, 0x72, 0x73, 0x2e, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x22,
	0x26, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x20, 0x12, 0x1e, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x2f, 0x7b, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x7d, 0x2f, 0x63, 0x68, 0x61,
	0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x12, 0x71, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x43, 0x68,
	0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x12, 0x1f, 0x2e, 0x73, 0x72, 0x6f, 0x2e, 0x63, 0x68,
	0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74,
	0x65, 0x72, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x1a, 0x19, 0x2e, 0x73, 0x72, 0x6f, 0x2e, 0x63,
	0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63,
	0x74, 0x65, 0x72, 0x22, 0x25, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1f, 0x12, 0x1d, 0x2f, 0x76, 0x31,
	0x2f, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x2f, 0x7b, 0x63, 0x68, 0x61,
	0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x7d, 0x12, 0x7f, 0x0a, 0x0f, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x12, 0x26, 0x2e,
	0x73, 0x72, 0x6f, 0x2e, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x73, 0x72, 0x6f, 0x2e, 0x63, 0x68, 0x61, 0x72,
	0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72,
	0x22, 0x29, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x23, 0x22, 0x1e, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73,
	0x65, 0x72, 0x73, 0x2f, 0x7b, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x7d, 0x2f, 0x63, 0x68,
	0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x3a, 0x01, 0x2a, 0x12, 0x61, 0x0a, 0x0f, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x12, 0x19,
	0x2e, 0x73, 0x72, 0x6f, 0x2e, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x2e,
	0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x1b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15, 0x2a, 0x13, 0x2f, 0x76, 0x31, 0x2f, 0x63,
	0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x65,
	0x0a, 0x0d, 0x45, 0x64, 0x69, 0x74, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x12,
	0x19, 0x2e, 0x73, 0x72, 0x6f, 0x2e, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73,
	0x2e, 0x43, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x1a, 0x19, 0x2e, 0x73, 0x72, 0x6f,
	0x2e, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x43, 0x68, 0x61, 0x72,
	0x61, 0x63, 0x74, 0x65, 0x72, 0x22, 0x1e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x1a, 0x13, 0x2f,
	0x76, 0x31, 0x2f, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x2f, 0x7b, 0x69,
	0x64, 0x7d, 0x3a, 0x01, 0x2a, 0x12, 0x8b, 0x01, 0x0a, 0x14, 0x41, 0x64, 0x64, 0x43, 0x68, 0x61,
	0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x50, 0x6c, 0x61, 0x79, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1f,
	0x2e, 0x73, 0x72, 0x6f, 0x2e, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x2e,
	0x50, 0x6c, 0x61, 0x79, 0x54, 0x69, 0x6d, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a,
	0x1f, 0x2e, 0x73, 0x72, 0x6f, 0x2e, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73,
	0x2e, 0x50, 0x6c, 0x61, 0x79, 0x54, 0x69, 0x6d, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x22, 0x31, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x2b, 0x1a, 0x26, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x68,
	0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x73, 0x2f, 0x7b, 0x63, 0x68, 0x61, 0x72, 0x61, 0x63,
	0x74, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x7d, 0x2f, 0x70, 0x6c, 0x61, 0x79, 0x74, 0x69, 0x6d, 0x65,
	0x3a, 0x01, 0x2a, 0x42, 0x08, 0x5a, 0x06, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_characters_proto_rawDescOnce sync.Once
	file_characters_proto_rawDescData = file_characters_proto_rawDesc
)

func file_characters_proto_rawDescGZIP() []byte {
	file_characters_proto_rawDescOnce.Do(func() {
		file_characters_proto_rawDescData = protoimpl.X.CompressGZIP(file_characters_proto_rawDescData)
	})
	return file_characters_proto_rawDescData
}

var file_characters_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_characters_proto_goTypes = []interface{}{
	(*PlayTimeMessage)(nil),        // 0: sro.characters.PlayTimeMessage
	(*DeleteCharacterRequest)(nil), // 1: sro.characters.DeleteCharacterRequest
	(*CreateCharacterRequest)(nil), // 2: sro.characters.CreateCharacterRequest
	(*UserTarget)(nil),             // 3: sro.characters.UserTarget
	(*CharacterTarget)(nil),        // 4: sro.characters.CharacterTarget
	(*Character)(nil),              // 5: sro.characters.Character
	(*Location)(nil),               // 6: sro.characters.Location
	(*Characters)(nil),             // 7: sro.characters.Characters
	(*Gender)(nil),                 // 8: sro.characters.Gender
	(*Realm)(nil),                  // 9: sro.characters.Realm
	(*Genders)(nil),                // 10: sro.characters.Genders
	(*Realms)(nil),                 // 11: sro.characters.Realms
	(*wrapperspb.UInt64Value)(nil), // 12: google.protobuf.UInt64Value
	(*wrapperspb.StringValue)(nil), // 13: google.protobuf.StringValue
	(*emptypb.Empty)(nil),          // 14: google.protobuf.Empty
}
var file_characters_proto_depIdxs = []int32{
	12, // 0: sro.characters.Character.owner:type_name -> google.protobuf.UInt64Value
	13, // 1: sro.characters.Character.name:type_name -> google.protobuf.StringValue
	12, // 2: sro.characters.Character.gender:type_name -> google.protobuf.UInt64Value
	12, // 3: sro.characters.Character.realm:type_name -> google.protobuf.UInt64Value
	12, // 4: sro.characters.Character.play_time:type_name -> google.protobuf.UInt64Value
	6,  // 5: sro.characters.Character.location:type_name -> sro.characters.Location
	5,  // 6: sro.characters.Characters.characters:type_name -> sro.characters.Character
	8,  // 7: sro.characters.Genders.genders:type_name -> sro.characters.Gender
	9,  // 8: sro.characters.Realms.realms:type_name -> sro.characters.Realm
	14, // 9: sro.characters.CharactersService.GetAllGenders:input_type -> google.protobuf.Empty
	14, // 10: sro.characters.CharactersService.GetAllRealms:input_type -> google.protobuf.Empty
	14, // 11: sro.characters.CharactersService.GetAllCharacters:input_type -> google.protobuf.Empty
	3,  // 12: sro.characters.CharactersService.GetAllCharactersForUser:input_type -> sro.characters.UserTarget
	4,  // 13: sro.characters.CharactersService.GetCharacter:input_type -> sro.characters.CharacterTarget
	2,  // 14: sro.characters.CharactersService.CreateCharacter:input_type -> sro.characters.CreateCharacterRequest
	5,  // 15: sro.characters.CharactersService.DeleteCharacter:input_type -> sro.characters.Character
	5,  // 16: sro.characters.CharactersService.EditCharacter:input_type -> sro.characters.Character
	0,  // 17: sro.characters.CharactersService.AddCharacterPlayTime:input_type -> sro.characters.PlayTimeMessage
	10, // 18: sro.characters.CharactersService.GetAllGenders:output_type -> sro.characters.Genders
	11, // 19: sro.characters.CharactersService.GetAllRealms:output_type -> sro.characters.Realms
	7,  // 20: sro.characters.CharactersService.GetAllCharacters:output_type -> sro.characters.Characters
	7,  // 21: sro.characters.CharactersService.GetAllCharactersForUser:output_type -> sro.characters.Characters
	5,  // 22: sro.characters.CharactersService.GetCharacter:output_type -> sro.characters.Character
	5,  // 23: sro.characters.CharactersService.CreateCharacter:output_type -> sro.characters.Character
	14, // 24: sro.characters.CharactersService.DeleteCharacter:output_type -> google.protobuf.Empty
	5,  // 25: sro.characters.CharactersService.EditCharacter:output_type -> sro.characters.Character
	0,  // 26: sro.characters.CharactersService.AddCharacterPlayTime:output_type -> sro.characters.PlayTimeMessage
	18, // [18:27] is the sub-list for method output_type
	9,  // [9:18] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_characters_proto_init() }
func file_characters_proto_init() {
	if File_characters_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_characters_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlayTimeMessage); i {
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
		file_characters_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteCharacterRequest); i {
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
		file_characters_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateCharacterRequest); i {
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
		file_characters_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
		file_characters_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CharacterTarget); i {
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
		file_characters_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Character); i {
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
		file_characters_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
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
		file_characters_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Characters); i {
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
		file_characters_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Gender); i {
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
		file_characters_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Realm); i {
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
		file_characters_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Genders); i {
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
		file_characters_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Realms); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_characters_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_characters_proto_goTypes,
		DependencyIndexes: file_characters_proto_depIdxs,
		MessageInfos:      file_characters_proto_msgTypes,
	}.Build()
	File_characters_proto = out.File
	file_characters_proto_rawDesc = nil
	file_characters_proto_goTypes = nil
	file_characters_proto_depIdxs = nil
}
