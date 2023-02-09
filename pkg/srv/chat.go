package srv

import (
	"context"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/interceptor"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/model"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/pb"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type chatServiceServer struct {
	pb.UnimplementedChatServiceServer
	chatService service.ChatService
	jwtService  service.JWTService
}

func (s chatServiceServer) ConnectChannel(request *pb.ChannelIdMessage, server pb.ChatService_ConnectChannelServer) error {
	r := s.chatService.ChannelMessagesReader(server.Context(), uint(request.ChannelId))

	for {
		msg, err := r.ReadMessage(server.Context())
		if err != nil {
			_ = r.Close()
			return err
		}

		err = server.Send(&pb.ChatMessage{
			Username: string(msg.Key),
			Message:  string(msg.Value),
		})

		if err != nil {
			_ = r.Close()
			return err
		}
	}
}

func (s chatServiceServer) SendChatMessage(ctx context.Context, request *pb.SendChatMessageRequest) (*empty.Empty, error) {
	claims, err := interceptor.ExtractCtxClaims(ctx, s.jwtService)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, s.chatService.SendChannelMessage(ctx, claims["preferred_username"].(string), request.Message, uint(request.ChannelId))
}

func (s chatServiceServer) SendDirectMessage(ctx context.Context, request *pb.SendDirectMessageRequest) (*empty.Empty, error) {
	claims, err := interceptor.ExtractCtxClaims(ctx, s.jwtService)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, s.chatService.SendDirectMessage(ctx, claims["preferred_username"].(string), request.Message, request.Username)
}

func (s chatServiceServer) ConnectDirectMessage(_ *empty.Empty, srv pb.ChatService_ConnectDirectMessageServer) error {
	claims, err := interceptor.ExtractCtxClaims(srv.Context(), s.jwtService)
	if err != nil {
		return err
	}

	r := s.chatService.DirectMessagesReader(srv.Context(), claims["preferred_username"].(string))
	for {
		msg, err := r.ReadMessage(srv.Context())
		if err != nil {
			_ = r.Close()
			return err
		}

		err = srv.Send(&pb.ChatMessage{
			Message:  string(msg.Value),
			Username: string(msg.Key),
		})

		if err != nil {
			_ = r.Close()
			return err
		}
	}
}

func (s chatServiceServer) CreateChannel(ctx context.Context, msg *pb.CreateChannelMessage) (*empty.Empty, error) {
	_, err := s.chatService.CreateChannel(ctx, &model.ChatChannel{
		Name:   msg.Name,
		Public: msg.Public,
	})

	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s chatServiceServer) DeleteChannel(ctx context.Context, msg *pb.ChannelIdMessage) (*empty.Empty, error) {
	err := s.chatService.DeleteChannel(ctx, &model.ChatChannel{
		Model: gorm.Model{ID: uint(msg.ChannelId)},
	})

	return &emptypb.Empty{}, err
}

func (s chatServiceServer) AllChatChannels(ctx context.Context, _ *empty.Empty) (*pb.ChatChannels, error) {
	channels, err := s.chatService.AllChannels(ctx)
	if err != nil {
		return nil, err
	}

	resp := &pb.ChatChannels{
		Channels: make([]*pb.ChatChannel, len(channels)),
	}

	for i, c := range channels {
		resp.Channels[i] = c.ToPb()
	}

	return resp, nil
}

func (s chatServiceServer) GetChannel(ctx context.Context, msg *pb.ChannelIdMessage) (*pb.ChatChannel, error) {
	channel, err := s.chatService.GetChannel(ctx, uint(msg.ChannelId))
	if err != nil {
		return nil, err
	}

	return channel.ToPb(), nil
}

func (s chatServiceServer) EditChannel(ctx context.Context, msg *pb.UpdateChatChannelRequest) (*empty.Empty, error) {
	return &emptypb.Empty{}, s.chatService.UpdateChannel(ctx, msg)
}

func NewChatServiceServer(chatService service.ChatService, jwtService service.JWTService) pb.ChatServiceServer {
	return chatServiceServer{
		chatService: chatService,
		jwtService:  jwtService,
	}
}
