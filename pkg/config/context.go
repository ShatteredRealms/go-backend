package config

import (
	"context"
	"fmt"
	"time"

	"github.com/ShatteredRealms/go-backend/pkg/auth"
	"github.com/ShatteredRealms/go-backend/pkg/common"
	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/WilSimpson/gocloak/v13"
	"github.com/golang-jwt/jwt/v4"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

var (
	parser jwt.Parser
)

type ServerContext struct {
	GlobalConfig   *GlobalConfig
	KeycloakClient *gocloak.GoCloak
	Tracer         trace.Tracer
	RefSROServer   *SROServer

	jwt            *gocloak.JWT
	tokenExpiresAt time.Time

	characterClientWrapper     *characterClientWrapper
	chatClientWrapper          *chatClientWrapper
	servermanagerClientWrapper *servermanagerClientWrapper
}

type grpcClientManager struct {
	conn *grpc.ClientConn
	addr string
}

type characterClientWrapper struct {
	manager *grpcClientManager
	client  pb.CharacterServiceClient
}
type chatClientWrapper struct {
	manager *grpcClientManager
	client  pb.ChatServiceClient
}
type servermanagerClientWrapper struct {
	manager *grpcClientManager
	client  pb.ServerManagerServiceClient
}

// IsReady checks if the connection is ready to be used.
func (g *grpcClientManager) IsReady() bool {
	return g.conn != nil && g.conn.GetState() == connectivity.Ready
}

// Connect creates a new connection to the gRPC server.
func (g *grpcClientManager) Connect() error {
	if g.IsReady() {
		return nil
	}

	var err error
	g.conn, err = helpers.GrpcClientWithOtel(g.addr)
	return err
}

func NewServerContext(ctx context.Context, conf *GlobalConfig, tracer trace.Tracer, ref *SROServer) *ServerContext {
	server := &ServerContext{
		GlobalConfig:   conf,
		KeycloakClient: gocloak.NewClient(conf.Keycloak.BaseURL),
		Tracer:         tracer,
		RefSROServer:   ref,

		characterClientWrapper: &characterClientWrapper{
			manager: &grpcClientManager{
				addr: conf.Character.Remote.Address(),
			}},
		chatClientWrapper: &chatClientWrapper{
			manager: &grpcClientManager{
				addr: conf.Chat.Remote.Address(),
			}},
		servermanagerClientWrapper: &servermanagerClientWrapper{
			manager: &grpcClientManager{
				addr: conf.GameBackend.Remote.Address(),
			}},
	}
	server.KeycloakClient.RegisterMiddlewares(gocloak.OpenTelemetryMiddleware)
	log.Logger.Level = ref.LogLevel

	return server
}

func (srvCtx *ServerContext) GetChatClient() (pb.ChatServiceClient, error) {
	if srvCtx.chatClientWrapper.manager.IsReady() {
		return srvCtx.chatClientWrapper.client, nil
	}

	if err := srvCtx.chatClientWrapper.manager.Connect(); err != nil {
		return nil, fmt.Errorf("connecting to chat service: %w", err)
	}

	srvCtx.chatClientWrapper.client = pb.NewChatServiceClient(srvCtx.chatClientWrapper.manager.conn)
	return srvCtx.chatClientWrapper.client, nil
}

func (srvCtx *ServerContext) GetCharacterClient() (pb.CharacterServiceClient, error) {
	if srvCtx.characterClientWrapper.manager.IsReady() {
		return srvCtx.characterClientWrapper.client, nil
	}

	if err := srvCtx.characterClientWrapper.manager.Connect(); err != nil {
		return nil, fmt.Errorf("connecting to character service: %w", err)
	}

	srvCtx.characterClientWrapper.client = pb.NewCharacterServiceClient(srvCtx.characterClientWrapper.manager.conn)
	return srvCtx.characterClientWrapper.client, nil
}

func (srvCtx *ServerContext) GetServerManagerClient() (pb.ServerManagerServiceClient, error) {
	if srvCtx.servermanagerClientWrapper.manager.IsReady() {
		return srvCtx.servermanagerClientWrapper.client, nil
	}

	if err := srvCtx.servermanagerClientWrapper.manager.Connect(); err != nil {
		return nil, fmt.Errorf("connecting to server manager service: %w", err)
	}

	srvCtx.servermanagerClientWrapper.client = pb.NewServerManagerServiceClient(srvCtx.servermanagerClientWrapper.manager.conn)

	return srvCtx.servermanagerClientWrapper.client, nil
}

func (srvCtx *ServerContext) loginClient(ctx context.Context) (*gocloak.JWT, error) {
	var err error
	srvCtx.jwt, err = srvCtx.KeycloakClient.LoginClient(
		ctx,
		srvCtx.RefSROServer.Keycloak.ClientId,
		srvCtx.RefSROServer.Keycloak.ClientSecret,
		srvCtx.GlobalConfig.Keycloak.Realm,
	)
	if err != nil {
		return nil, fmt.Errorf("login keycloak: %v", err)
	}

	claims := &jwt.RegisteredClaims{}
	_, _, err = parser.ParseUnverified(srvCtx.jwt.AccessToken, claims)
	if err != nil {
		log.Logger.Errorf("parsing access token: %v", err)
		return srvCtx.jwt, nil
	}

	// Remove 5 seconds to ensure there are no race cases with expiration
	srvCtx.tokenExpiresAt = claims.ExpiresAt.Time.Add(-5 * time.Second)
	return srvCtx.jwt, nil
}

func (srvCtx *ServerContext) OutgoingClientAuth(ctx context.Context) (context.Context, error) {
	token, err := srvCtx.GetJWT(ctx)
	if err != nil {
		return ctx, err
	}

	return auth.AddOutgoingToken(
		ctx,
		token.AccessToken,
	), nil
}

func (srvCtx *ServerContext) GetJWT(ctx context.Context) (*gocloak.JWT, error) {
	if srvCtx.jwt != nil && time.Now().Before(srvCtx.tokenExpiresAt) {
		return srvCtx.jwt, nil
	}

	return srvCtx.loginClient(ctx)
}

func (srvCtx *ServerContext) GetCharacterNameFromTarget(
	ctx context.Context,
	target *pb.CharacterTarget,
) (string, error) {
	if target == nil {
		return "", fmt.Errorf("target cannot be nil")
	}

	targetCharacterName := ""
	switch t := target.Type.(type) {
	case *pb.CharacterTarget_Name:
		targetCharacterName = t.Name

	case *pb.CharacterTarget_Id:
		characterClient, err := srvCtx.GetCharacterClient()
		if err != nil {
			return "", err
		}

		targetChar, err := characterClient.GetCharacter(ctx, target)
		if err != nil {
			return "", err
		}
		targetCharacterName = targetChar.Name

	default:
		log.Logger.WithContext(ctx).Errorf("target type unknown: %+v", target)
		return "", common.ErrHandleRequest.Err()

	}

	return targetCharacterName, nil
}

func (srvCtx *ServerContext) GetCharacterIdFromTarget(
	ctx context.Context,
	target *pb.CharacterTarget,
) (uint, error) {
	if target == nil {
		return 0, fmt.Errorf("target cannot be nil")
	}

	targetCharacterId := uint(0)
	switch t := target.Type.(type) {
	case *pb.CharacterTarget_Name:
		characterClient, err := srvCtx.GetCharacterClient()
		if err != nil {
			return 0, err
		}

		targetChar, err := characterClient.GetCharacter(ctx, target)
		if err != nil {
			return 0, err
		}
		targetCharacterId = uint(targetChar.Id)

	case *pb.CharacterTarget_Id:
		targetCharacterId = uint(t.Id)

	default:
		log.Logger.WithContext(ctx).Errorf("target type unknown: %+v", target)
		return 0, common.ErrHandleRequest.Err()
	}

	return targetCharacterId, nil
}
