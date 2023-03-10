package srv

import (
	"context"
	"database/sql"

	"github.com/ShatteredRealms/go-backend/pkg/interceptor"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/service"
	utilService "github.com/ShatteredRealms/go-backend/pkg/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gopkg.in/nullbio/null.v4"
)

// @TODO(wil): Change all errors to variables

type userServiceServer struct {
	pb.UnimplementedUserServiceServer
	userService       service.UserService
	permissionService service.PermissionService
	jwtService        utilService.JWTService
}

func NewUserServiceServer(
	u service.UserService,
	p service.PermissionService,
	j utilService.JWTService,
) *userServiceServer {
	return &userServiceServer{
		userService:       u,
		permissionService: p,
		jwtService:        j,
	}
}

func (s *userServiceServer) GetAll(
	ctx context.Context,
	message *emptypb.Empty,
) (*pb.GetAllUsersResponse, error) {
	users := s.userService.FindAll(ctx)
	resp := &pb.GetAllUsersResponse{
		Users: make([]*pb.UserMessage, len(users)),
	}

	for i, u := range users {
		resp.Users[i] = u.ToPb()
	}

	return resp, nil
}

func (s *userServiceServer) Get(
	ctx context.Context,
	message *pb.GetUserMessage,
) (*pb.GetUserResponse, error) {
	user := s.userService.FindByUsername(ctx, message.Username)
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return user.ToVerbosePb(s.permissionService.FindPermissionsForUsername(ctx, user.Username)), nil
}

func (s *userServiceServer) Edit(
	ctx context.Context,
	message *pb.EditUserDetailsRequest,
) (*emptypb.Empty, error) {
	can, err := interceptor.AuthorizedForTarget(ctx, s.jwtService, message.Username)
	if err != nil || !can {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}

	user := s.userService.FindByUsername(ctx, message.Username)
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	err = user.UpdateInfo(message)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = s.userService.Save(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *userServiceServer) Ban(ctx context.Context, message *pb.GetUserMessage) (*emptypb.Empty, error) {
	user := s.userService.FindByUsername(ctx, message.Username)
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	err := s.userService.Ban(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to ban user: %v", err.Error())
	}

	return &emptypb.Empty{}, nil
}
func (s *userServiceServer) UnBan(ctx context.Context, message *pb.GetUserMessage) (*emptypb.Empty, error) {
	user := s.userService.FindByUsername(ctx, message.Username)
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	err := s.userService.UnBan(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to unban user: %v", err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *userServiceServer) GetStatus(ctx context.Context, message *pb.GetUserMessage) (*pb.StatusResponse, error) {
	can, err := interceptor.AuthorizedForTarget(ctx, s.jwtService, message.Username)
	if err != nil || !can {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}

	user := s.userService.FindByUsername(ctx, message.Username)
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.StatusResponse{CharacterName: user.CurrentCharacterWrapper()}, nil
}

func (s *userServiceServer) SetStatus(ctx context.Context, message *pb.RequestSetStatus) (*emptypb.Empty, error) {
	can, err := interceptor.AuthorizedForTarget(ctx, s.jwtService, message.Username)
	if err != nil || !can {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}

	user := s.userService.FindByUsername(ctx, message.Username)
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	if message.CharacterName == nil {
		user.CurrentCharacter = null.String{}
	} else {
		user.CurrentCharacter = null.String{NullString: sql.NullString{String: message.CharacterName.Value}}
	}

	_, err = s.userService.Save(ctx, user)
	if err != nil {
		if user == nil {
			return nil, status.Error(codes.Internal, "unable to update user")
		}
	}

	return &emptypb.Empty{}, nil
}
func (s *userServiceServer) ChangePassword(ctx context.Context, message *pb.ChangePasswordRequest) (*emptypb.Empty, error) {
	can, err := interceptor.AuthorizedForTarget(ctx, s.jwtService, message.Username)
	if err != nil || !can {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}

	user := s.userService.FindByUsername(ctx, message.Username)
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	err = user.Login(message.CurrentPassword)
	if err != nil || !can {
		return nil, status.Error(codes.FailedPrecondition, "Incorrect password")
	}

	err = user.UpdatePassword(message.NewPassword)
	if err != nil || !can {
		return nil, status.Error(codes.FailedPrecondition, "Invalid new password")
	}

	_, err = s.userService.Save(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
