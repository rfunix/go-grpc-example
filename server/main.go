package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rfunix/grpc_api/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
	"net/http"
	"strings"
)

type contextKey int

const (
	clientIDKey contextKey = iota
)

func credMatcher(headerName string) (mdName string, ok bool) {
	if headerName == "Login" || headerName == "Password" {
		return headerName, true
	}

	return "", false
}

func authenticateClient(ctx context.Context, s *api.Server) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		clientLogin := strings.Join(md["login"], "")
		clientPassword := strings.Join(md["password"], "")

		if clientLogin != "rfunix" {
			return "", fmt.Errorf("unknown user %s", clientLogin)
		}

		if clientPassword != "xpto1234" {
			return "", fmt.Errorf("bad password %s", clientPassword)
		}

		log.Printf("authenticated client: %s", clientLogin)

		return "42", nil

	}

	return "", fmt.Errorf("missing credentials")

}

func unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	s, ok := info.Server.(*api.Server)
	if !ok {
		return nil, fmt.Errorf("unable to cast server")
	}

	clientID, err := authenticateClient(ctx, s)
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, clientIDKey, clientID)
	return handler(ctx, req)
}

func startGRPCServer(address, certFile, keyFile string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := api.Server{}

	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		log.Fatalf("could not load TLS keys: %s", err)
	}

	opts := []grpc.ServerOption{grpc.Creds(creds),
		grpc.UnaryInterceptor(unaryInterceptor)}

	grpcServer := grpc.NewServer(opts...)

	api.RegisterPingServer(grpcServer, &s)

	log.Printf("starting HTTP/2 gRPC server on %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

	return nil
}

func startRESTServer(address, grpcAddress, certFile string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(credMatcher))
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		return fmt.Errorf("could not load TLS certificate: %s", err)
	}

	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	err = api.RegisterPingHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
	if err != nil {
		return fmt.Errorf("could not register service Ping: %s", err)
	}

	log.Printf("starting HTTP/1.1 REST server on %s", address)

	http.ListenAndServe(address, mux)

	return nil
}

func main() {
	grpcAddress := fmt.Sprintf("%s:%d", "localhost", 7777)
	restAddress := fmt.Sprintf("%s:%d", "localhost", 7778)
	certFile := "cert/server.crt"
	keyFile := "cert/server.key"

	go func() {
		err := startGRPCServer(grpcAddress, certFile, keyFile)
		if err != nil {
			log.Fatalf("failed to start gRPC server: %s", err)
		}
	}()

	go func() {
		err := startRESTServer(restAddress, grpcAddress, certFile)
		if err != nil {
			log.Fatalf("failed to start gRPC server: %s", err)
		}
	}()
	log.Printf("Entering infinite loop")
	select {}
}
