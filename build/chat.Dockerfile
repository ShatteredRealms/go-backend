# Build application
FROM golang:1.19 AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY ./ ./
RUN go mod download
ARG APP_VERSION=v0.0.1
RUN go build \
	-ldflags="-X 'github.com/ShatteredRealms/go-backend/pkg/config/default.Version=${APP_VERSION}'" \
	-o /out/chat ./cmd/chat

# Run server
FROM alpine:3.15.0
WORKDIR /app
COPY --from=build /out/chat ./
EXPOSE 8180
ENTRYPOINT [ "./chat" ]
