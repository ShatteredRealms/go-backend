package service

import "go.opentelemetry.io/otel"

var (
	gamebackendTracer = otel.Tracer("Inner-GamebackendService")
)

type gamebackendService struct {
	chatRepo repository.GamebackendRepository
}
