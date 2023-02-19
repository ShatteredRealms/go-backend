package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/service"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	channelID = 1
)

func main() {
	ctx := context.Background()

	jwtService, err := service.NewJWTService("./test/auth")
	helpers.Check(ctx, err, "jwt service")

	claims := jwt.MapClaims{
		"sub":                1,
		"preferred_username": "username",
		//"given_name":  u.FirstName,
		//"family_name": u.LastName,
		//"email":       u.Email,
	}

	token, err := jwtService.Create(ctx, time.Hour, "testing.fake", claims)
	helpers.Check(ctx, err, "auth token")
	//token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzQ2Njk4NzEsImlhdCI6MTY3NDY2NjI3MSwiaXNzIjoic3JvLmNvbS9hY2NvdW50cy92MSIsIm5iZiI6MTY3NDY2NjI3MSwicHJlZmVycmVkX3VzZXJuYW1lIjoidW5yZWFsIiwic3ViIjoxfQ.GBhAozhYVWVpwk6abObUZIhZuWLh74RKPQsvo6Mtdvlq-6RPOI39Y_LfjqGTagv7REvnhDPpbo7kkHykWEVz9A2T66cjEMts3ITlSa6tMUklvguqX8vxhRyzLa5eb1Ba7K7xeIYDI6DoCXxkiwJ3XEhImtYLKj-ZhsdZX4mBHtnGHQzgF7mMTSagG12WQkclPekcdVO6xLDPR5RXeE2Fq_vNDARXYkkSg7PF2rR9Qvn99zJK_MxQZtg7MIQxmNd_1K-P22BhmjX40eDx6ZMIPVkVvVASOA2Lq5FMWStNk7fqkssKvw2JA6aXtwZp3cRt3zdXe5H5wpVqSD1-iI2MRw"
	token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzYzNDU0MDQsImlhdCI6MTY3NjM0MTgwNCwiaXNzIjoic3JvLmNvbS9hY2NvdW50cy92MSIsIm5iZiI6MTY3NjM0MTgwNCwicHJlZmVycmVkX3VzZXJuYW1lIjoidW5yZWFsIiwic3ViIjoxfQ.L8P09rvezAXVjbPbsILKOCgHbd9IsItT4zjCz1xbcWxeY_hmZ0hePvdv0TYJnXjvPxSYW74yQuTeYnmvloUHOfYBayTAIIY9EAtVGZTXtbpT02MUAFr2tuUPH7sy0ykxvaXCfj5AWniBchuUHoeRLv3M8mmBDvkN-RlgU0817j6SdsUvrlSXHi_CmIUYvP5rrctZLPi_SclhWYNBGfeMrNl_JMxqUxEwwltvt02sf-rwghCJeHkWdESsgkuZNAEYdlmzwX2rbPFxG63anhCrqMEL3H9qNUidf_gE3deRR74zf_Olh0JhDbxK0voerEQW_Hz3fHQu1MSIu7tQPX_dKw"
	md := metadata.New(
		map[string]string{
			"authorization": fmt.Sprintf(
				"Bearer %s", token,
			),
			//"host":  "api.shatteredrealmsonline.com",
			//":path": "chat",
		},
	)
	authCtx := metadata.NewOutgoingContext(context.Background(), md)

	fmt.Println("Connecting...")

	clientConn, err := grpc.Dial("chat.grpc.api.shatteredrealmsonline.com:8180", grpc.WithTransportCredentials(insecure.NewCredentials()))
	//clientConn, err := grpc.Dial("shatteredrealmsonline.com:443", grpc.WithTransportCredentials(insecure.NewCredentials()))
	helpers.Check(ctx, err, "dial chat")

	chatClient := pb.NewChatServiceClient(clientConn)
	stream, err := chatClient.ConnectChannel(authCtx, &pb.ChannelIdMessage{ChannelId: channelID})
	helpers.Check(ctx, err, "connect channel")

	fmt.Println("Connected")

	go func() {
		var msg *pb.ChatMessage

		for err == nil {
			msg, err = stream.Recv()

			if err == io.EOF {
				fmt.Println("Chat connection closed.")
				os.Exit(0)
				return
			}

			if err == nil && msg != nil {
				fmt.Printf("Recieved Message:: %s: %s\n", msg.Username, msg.Message)
			}
		}
		fmt.Printf("\nERROR: %v\n", err)
		fmt.Println("Disconnected")
		os.Exit(1)
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		msg, _ := reader.ReadString('\n')
		msg = strings.ReplaceAll(msg, "\n", "")
		msg = strings.ReplaceAll(msg, "\r", "")
		_, err = chatClient.SendChatMessage(authCtx,
			&pb.SendChatMessageRequest{
				ChannelId: channelID,
				Message:   msg,
			},
		)

		if err != nil {
			fmt.Printf("SEND ERROR: %v\n", err)
		}
	}
}
