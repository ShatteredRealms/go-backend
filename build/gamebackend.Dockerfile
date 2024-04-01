# Build application
FROM golang:1.22 AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ARG APP_VERSION=v0.0.1
RUN go build \
	-ldflags="-X 'github.com/ShatteredRealms/go-backend/pkg/config/default.Version=${APP_VERSION}'" \
	-o /out/gamebackend ./cmd/gamebackend

# Run server
FROM alpine:3.15.0
WORKDIR /app
COPY --from=build /out/gamebackend ./
EXPOSE 8082
ENTRYPOINT [ "./gamebackend" ]
