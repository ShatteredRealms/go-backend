# Build application
FROM golang:1.19 AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY ./ ./
RUN go mod download
RUN go build -o /out/accounts ./cmd/accounts

# Run server
FROM alpine:3.15.0
WORKDIR /app
COPY --from=build /out/accounts ./
EXPOSE 8080
CMD [ "./accounts" ]
