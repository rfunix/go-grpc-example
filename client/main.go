package main

import (
	"context"
	"github.com/rfunix/grpc_api/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

type Authentication struct {
	Login    string
	Password string
}

func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"login":    a.Login,
		"password": a.Password,
	}, nil
}

func (a *Authentication) RequireTransportSecurity() bool {
	return true
}

func main() {
	var conn *grpc.ClientConn
	creds, err := credentials.NewClientTLSFromFile("cert/server.crt", "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}

	auth := Authentication{
		Login:    "rfunix",
		Password: "xpto123",
	}

	conn, err = grpc.Dial(
		"localhost:7777",
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(&auth))

	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	defer conn.Close()

	c := api.NewPingClient(conn)

	response, err := c.SayHello(context.Background(), &api.PingMessage{Greeting: "foo"})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}

	log.Printf("Response from server: %s", response.Greeting)
}
