package helpers

import (
	"context"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net"
	"net/http"
	"strings"
)

func GRPCHandlerFunc(grpcServer http.Handler, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType, x-grpc-web")

			if r.Method == "OPTIONS" {
				return
			}

			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

func StartServer(
	ctx context.Context,
	grpcServer http.Handler,
	gwmux http.Handler,
	address string,
) net.Listener {
	log.WithContext(ctx).Info("Starting server")
	listen, err := net.Listen("tcp", address)
	Check(ctx, err, "listen server")

	httpSrv := &http.Server{
		Addr:    address,
		Handler: GRPCHandlerFunc(grpcServer, gwmux),
	}

	go func() {
		Check(ctx, httpSrv.Serve(listen), "server stopped")
	}()

	return listen
}
