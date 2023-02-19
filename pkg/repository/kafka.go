package repository

import (
	"context"
	"net"
	"strconv"

	chat "github.com/ShatteredRealms/go-backend/cmd/chat/global"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/segmentio/kafka-go"
)

var (
	currentConn    *kafka.Conn
	controllerConn *kafka.Conn
)

func ConnectKafka(ctx context.Context, address config.ServerAddress) (*kafka.Conn, error) {
	if currentConn != nil {
		_ = currentConn.Close()
	}
	if controllerConn != nil {
		_ = controllerConn.Close()
	}
	var err error

	ctx, span := chat.Tracer.Start(ctx, "connect-kafka")
	currentConn, err = kafka.Dial("tcp", address.Address())
	if err != nil {
		return nil, err
	}

	controller, err := currentConn.Controller()
	if err != nil {
		return nil, err
	}

	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return nil, err
	}

	span.End()

	return controllerConn, nil
}
