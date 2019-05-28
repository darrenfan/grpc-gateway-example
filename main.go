package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	proto "example/proto"
)

const (
	listenPort = ":50001"
)

type Example struct {
}

func (t *Example) Echo(ctx context.Context, request *proto.StringMessage) (*proto.StringMessage, error) {
	return &proto.StringMessage{Value: "hello"}, nil
}

func main() {
    var err error

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	proto.RegisterExampleServer(grpcServer, &Example{})

	mux := http.NewServeMux()
	gwmux := runtime.NewServeMux()

	ctx := context.Background()

	dopts := []grpc.DialOption{grpc.WithInsecure()}

	err = proto.RegisterExampleHandlerFromEndpoint(ctx, gwmux, "localhost"+listenPort, dopts)
	if err != nil {
		log.Fatalf("failed to register handler: %v", err)
	}
	mux.Handle("/", gwmux)

	log.Printf("start to listen server localhost%s", listenPort)
	http.ListenAndServe(listenPort, grpcHandlerFunc(grpcServer, mux))
}

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from cockroachdb.
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print(r.Header.Get("Content-Type"))
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}