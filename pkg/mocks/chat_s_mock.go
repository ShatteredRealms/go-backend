// Code generated by MockGen. DO NOT EDIT.
// Source: /home/wil/sro/git/go-backend/pkg/service/chat_s.go
//
// Generated by this command:
//
//	mockgen -package=mocks -source=/home/wil/sro/git/go-backend/pkg/service/chat_s.go -destination=/home/wil/sro/git/go-backend/pkg/mocks/chat_s_mock.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	model "github.com/ShatteredRealms/go-backend/pkg/model"
	pb "github.com/ShatteredRealms/go-backend/pkg/pb"
	kafka "github.com/segmentio/kafka-go"
	gomock "go.uber.org/mock/gomock"
)

// MockChatService is a mock of ChatService interface.
type MockChatService struct {
	ctrl     *gomock.Controller
	recorder *MockChatServiceMockRecorder
}

// MockChatServiceMockRecorder is the mock recorder for MockChatService.
type MockChatServiceMockRecorder struct {
	mock *MockChatService
}

// NewMockChatService creates a new mock instance.
func NewMockChatService(ctrl *gomock.Controller) *MockChatService {
	mock := &MockChatService{ctrl: ctrl}
	mock.recorder = &MockChatServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChatService) EXPECT() *MockChatServiceMockRecorder {
	return m.recorder
}

// AllChannels mocks base method.
func (m *MockChatService) AllChannels(ctx context.Context) (model.ChatChannels, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllChannels", ctx)
	ret0, _ := ret[0].(model.ChatChannels)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllChannels indicates an expected call of AllChannels.
func (mr *MockChatServiceMockRecorder) AllChannels(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllChannels", reflect.TypeOf((*MockChatService)(nil).AllChannels), ctx)
}

// AuthorizedChannelsForCharacter mocks base method.
func (m *MockChatService) AuthorizedChannelsForCharacter(ctx context.Context, characterId uint) (model.ChatChannels, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthorizedChannelsForCharacter", ctx, characterId)
	ret0, _ := ret[0].(model.ChatChannels)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AuthorizedChannelsForCharacter indicates an expected call of AuthorizedChannelsForCharacter.
func (mr *MockChatServiceMockRecorder) AuthorizedChannelsForCharacter(ctx, characterId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthorizedChannelsForCharacter", reflect.TypeOf((*MockChatService)(nil).AuthorizedChannelsForCharacter), ctx, characterId)
}

// ChangeAuthorizationForCharacter mocks base method.
func (m *MockChatService) ChangeAuthorizationForCharacter(ctx context.Context, characterId uint, channelIds []uint, addAuth bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeAuthorizationForCharacter", ctx, characterId, channelIds, addAuth)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeAuthorizationForCharacter indicates an expected call of ChangeAuthorizationForCharacter.
func (mr *MockChatServiceMockRecorder) ChangeAuthorizationForCharacter(ctx, characterId, channelIds, addAuth any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeAuthorizationForCharacter", reflect.TypeOf((*MockChatService)(nil).ChangeAuthorizationForCharacter), ctx, characterId, channelIds, addAuth)
}

// ChannelMessagesReader mocks base method.
func (m *MockChatService) ChannelMessagesReader(ctx context.Context, channelId uint) *kafka.Reader {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChannelMessagesReader", ctx, channelId)
	ret0, _ := ret[0].(*kafka.Reader)
	return ret0
}

// ChannelMessagesReader indicates an expected call of ChannelMessagesReader.
func (mr *MockChatServiceMockRecorder) ChannelMessagesReader(ctx, channelId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChannelMessagesReader", reflect.TypeOf((*MockChatService)(nil).ChannelMessagesReader), ctx, channelId)
}

// CreateChannel mocks base method.
func (m *MockChatService) CreateChannel(ctx context.Context, channel *model.ChatChannel) (*model.ChatChannel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateChannel", ctx, channel)
	ret0, _ := ret[0].(*model.ChatChannel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateChannel indicates an expected call of CreateChannel.
func (mr *MockChatServiceMockRecorder) CreateChannel(ctx, channel any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateChannel", reflect.TypeOf((*MockChatService)(nil).CreateChannel), ctx, channel)
}

// DeleteChannel mocks base method.
func (m *MockChatService) DeleteChannel(ctx context.Context, channel *model.ChatChannel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteChannel", ctx, channel)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteChannel indicates an expected call of DeleteChannel.
func (mr *MockChatServiceMockRecorder) DeleteChannel(ctx, channel any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteChannel", reflect.TypeOf((*MockChatService)(nil).DeleteChannel), ctx, channel)
}

// DirectMessagesReader mocks base method.
func (m *MockChatService) DirectMessagesReader(ctx context.Context, username string) *kafka.Reader {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DirectMessagesReader", ctx, username)
	ret0, _ := ret[0].(*kafka.Reader)
	return ret0
}

// DirectMessagesReader indicates an expected call of DirectMessagesReader.
func (mr *MockChatServiceMockRecorder) DirectMessagesReader(ctx, username any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DirectMessagesReader", reflect.TypeOf((*MockChatService)(nil).DirectMessagesReader), ctx, username)
}

// GetChannel mocks base method.
func (m *MockChatService) GetChannel(ctx context.Context, id uint) (*model.ChatChannel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChannel", ctx, id)
	ret0, _ := ret[0].(*model.ChatChannel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChannel indicates an expected call of GetChannel.
func (mr *MockChatServiceMockRecorder) GetChannel(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChannel", reflect.TypeOf((*MockChatService)(nil).GetChannel), ctx, id)
}

// SendChannelMessage mocks base method.
func (m *MockChatService) SendChannelMessage(ctx context.Context, username, message string, channelId uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendChannelMessage", ctx, username, message, channelId)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendChannelMessage indicates an expected call of SendChannelMessage.
func (mr *MockChatServiceMockRecorder) SendChannelMessage(ctx, username, message, channelId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendChannelMessage", reflect.TypeOf((*MockChatService)(nil).SendChannelMessage), ctx, username, message, channelId)
}

// SendDirectMessage mocks base method.
func (m *MockChatService) SendDirectMessage(ctx context.Context, senderCharacter, message, targetCharacter string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendDirectMessage", ctx, senderCharacter, message, targetCharacter)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendDirectMessage indicates an expected call of SendDirectMessage.
func (mr *MockChatServiceMockRecorder) SendDirectMessage(ctx, senderCharacter, message, targetCharacter any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendDirectMessage", reflect.TypeOf((*MockChatService)(nil).SendDirectMessage), ctx, senderCharacter, message, targetCharacter)
}

// UpdateChannel mocks base method.
func (m *MockChatService) UpdateChannel(ctx context.Context, pb *pb.UpdateChatChannelRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateChannel", ctx, pb)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateChannel indicates an expected call of UpdateChannel.
func (mr *MockChatServiceMockRecorder) UpdateChannel(ctx, pb any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateChannel", reflect.TypeOf((*MockChatService)(nil).UpdateChannel), ctx, pb)
}