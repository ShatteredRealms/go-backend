package helpers_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"

	"github.com/ShatteredRealms/go-backend/pkg/helpers"
)

var _ = Describe("Grpc helpers", func() {

	var (
		grpcServer = &fakeHttpHandler{}
		httpServer = &fakeHttpHandler{}
		resp       = fakeResponseWriter{}
		req        = &http.Request{}

		httpHandler http.Handler
	)

	log.StandardLogger().ExitFunc = func(int) { wasFatal = true }

	BeforeEach(func() {
		grpcServer.HandledRequest = false
		httpServer.HandledRequest = false
		req = &http.Request{
			ProtoMajor: 2,
			Header:     make(http.Header),
		}
		req.Header.Set("Content-Type", "application/grpc; otherdata=this")
	})

	Describe("GRPCHandlerFunc", func() {
		BeforeEach(func() {
			httpHandler = helpers.GRPCHandlerFunc(grpcServer, httpServer)
		})

		It("should forward grpc requests to the grpcServer", func() {
			httpHandler.ServeHTTP(resp, req)
			Expect(grpcServer.HandledRequest).To(BeTrue(), "should handle as gRPC request")
			Expect(httpServer.HandledRequest).To(BeFalse(), "should not handle as http request")
		})

		It("should handle http when not ProtoMajor 2", func() {
			req.ProtoMajor = 1
			httpHandler.ServeHTTP(resp, req)
			Expect(grpcServer.HandledRequest).To(BeFalse(), "should not handle as gRPC request")
			Expect(httpServer.HandledRequest).To(BeTrue(), "should handle as http request")
		})

		It("should handle http when content type is not gRPC", func() {
			req.Header.Set("Content-Type", "text/html")
			httpHandler.ServeHTTP(resp, req)
			Expect(grpcServer.HandledRequest).To(BeFalse(), "should not handle as gRPC request")
			Expect(httpServer.HandledRequest).To(BeTrue(), "should handle as http request")
		})

		It("should not serve http OPTIONS", func() {
			req.Header.Set("Content-Type", "text/html")
			req.Method = "OPTIONS"
			httpHandler.ServeHTTP(resp, req)
			Expect(grpcServer.HandledRequest).To(BeFalse(), "should not handle as gRPC request")
			Expect(httpServer.HandledRequest).To(BeFalse(), "should handle as http request")
		})
	})

	// @TODO: Find way to test without race conditions
	Describe("StartServer", func() {
		It("should start a server", func() {
			ctx := context.Background()
			listener := helpers.StartServer(ctx, grpcServer, httpServer, "127.0.0.1:9999")
			Expect(listener.Addr().String()).To(Equal("127.0.0.1:9999"))
			listener.Close()
		})
	})
})

type fakeHttpHandler struct {
	HandledRequest bool
}

func (handler *fakeHttpHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	handler.HandledRequest = true
}

type fakeResponseWriter struct {
	header http.Header
}

func (w fakeResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}

	return w.header
}

func (w fakeResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (w fakeResponseWriter) WriteHeader(statusCode int) {

}
