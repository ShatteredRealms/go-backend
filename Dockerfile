# Build application
FROM golang:1.17 AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY ./ ./
RUN go mod download
RUN go build -o /out/accounts ./cmd/accounts/main.go

# Run server
FROM alpine:3.15.0
WORKDIR /app
COPY --from=build /out/accounts ./
EXPOSE 8080
CMD [ "./accounts", "-port=8080", "-mode=release" ]
