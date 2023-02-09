package helpers

import (
    "context"
    log "github.com/sirupsen/logrus"
    easy "github.com/t-tomalak/logrus-easy-formatter"
    "github.com/uptrace/opentelemetry-go-extra/otellogrus"
    "google.golang.org/grpc"
    "os"
)

func UnaryLogRequest() grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        log.Info(info.FullMethod)
        return handler(ctx, req)
    }
}
func StreamLogRequest() grpc.StreamServerInterceptor {
    return func(
        srv interface{},
        stream grpc.ServerStream,
        info *grpc.StreamServerInfo,
        handler grpc.StreamHandler,
    ) error {
        log.Info(info.FullMethod)
        return handler(srv, stream)
    }
}

func SetupLogs() {
    log.AddHook(otellogrus.NewHook(
        otellogrus.WithLevels(
            log.PanicLevel,
            log.FatalLevel,
            log.ErrorLevel,
            log.WarnLevel)))

    log.SetOutput(os.Stdout)
    log.SetLevel(log.TraceLevel)
    log.SetFormatter(&easy.Formatter{
        TimestampFormat: "2006-01-02 15:04:05",
        LogFormat:       "%time% [%lvl%]: %msg%\n",
    })
}
