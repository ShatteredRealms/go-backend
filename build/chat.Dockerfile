# Build application
FROM golang:1.19 AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY ./ ./
RUN go mod download
RUN go build -o /out/chat ./cmd/chat

# Run server
FROM alpine:3.15.0
WORKDIR /app
COPY --from=build /out/chat ./
EXPOSE 8180
CMD [ "./chat" ]