package service

import (
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	"go.opentelemetry.io/otel"
)

var (
	gamebackendTracer = otel.Tracer("Inner-GamebackendService")
)

type gamebackendService struct {
	chatRepo repository.GamebackendRepository
}
