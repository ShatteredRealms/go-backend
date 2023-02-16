package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/helpers"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"sync"
	"time"
)

const (
	channelID = 1
	token     = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzYzNDU0MDQsImlhdCI6MTY3NjM0MTgwNCwiaXNzIjoic3JvLmNvbS9hY2NvdW50cy92MSIsIm5iZiI6MTY3NjM0MTgwNCwicHJlZmVycmVkX3VzZXJuYW1lIjoidW5yZWFsIiwic3ViIjoxfQ.L8P09rvezAXVjbPbsILKOCgHbd9IsItT4zjCz1xbcWxeY_hmZ0hePvdv0TYJnXjvPxSYW74yQuTeYnmvloUHOfYBayTAIIY9EAtVGZTXtbpT02MUAFr2tuUPH7sy0ykxvaXCfj5AWniBchuUHoeRLv3M8mmBDvkN-RlgU0817j6SdsUvrlSXHi_CmIUYvP5rrctZLPi_SclhWYNBGfeMrNl_JMxqUxEwwltvt02sf-rwghCJeHkWdESsgkuZNAEYdlmzwX2rbPFxG63anhCrqMEL3H9qNUidf_gE3deRR74zf_Olh0JhDbxK0voerEQW_Hz3fHQu1MSIu7tQPX_dKw"
)

var (
	md = metadata.New(
		map[string]string{
			"authorization": fmt.Sprintf(
				"Bearer %s", token,
			),
		},
	)

	c *int
	t *int

	count = 0
	mu    sync.Mutex
)

func main() {
	c = flag.Int("c", 100, "number of conncurrent connections")
	t = flag.Int("t", 1100, "number of ms between messages")

	for i := 0; i < *c; i++ {
		go connect()
	}

	for {
		fmt.Printf("\rMessages Sent: %d", count)
		time.Sleep(time.Second)
	}
}

func connect() {
	ctx := context.Background()

	authCtx := metadata.NewOutgoingContext(context.Background(), md)

	clientConn, err := grpc.Dial("chat.grpc.api.shatteredrealmsonline.com:8180", grpc.WithTransportCredentials(insecure.NewCredentials()))
	helpers.Check(ctx, err, "dial chat")

	chatClient := pb.NewChatServiceClient(clientConn)

	for err == nil {
		_, err = chatClient.SendChatMessage(authCtx,
			&pb.SendChatMessageRequest{
				ChannelId: channelID,
				Message:   "Test Message.",
			},
		)

		if err != nil {
			fmt.Printf("SEND ERROR: %v\n", err)
		} else {
		}

		mu.Lock()
		count++
		mu.Unlock()

		time.Sleep(time.Millisecond * time.Duration(*t))
	}
}
