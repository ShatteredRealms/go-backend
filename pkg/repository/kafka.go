package repository

import (
	"fmt"
	"net"
	"strconv"

	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/segmentio/kafka-go"
)

var (
	currentConn    *kafka.Conn
	controllerConn *kafka.Conn
)

func ConnectKafka(address config.ServerAddress) (*kafka.Conn, error) {
	if currentConn != nil {
		_ = currentConn.Close()
	}
	if controllerConn != nil {
		_ = controllerConn.Close()
	}
	var err error

	currentConn, err = kafka.Dial("tcp", address.Address())
	if err != nil {
		return nil, fmt.Errorf("kafka connect: %v", err)
	}

	controller, err := currentConn.Controller()
	if err != nil {
		return nil, fmt.Errorf("controller: %v", err)
	}

	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return nil, fmt.Errorf("kafka controller connection: %v", err)
	}

	return controllerConn, nil
}
