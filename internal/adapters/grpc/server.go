package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/chyiyaqing/gmicro-proto/golang/user"
	"github.com/chyiyaqing/gmicro-user/config"
	"github.com/chyiyaqing/gmicro-user/internal/ports"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Adapter struct {
	api    ports.APIPort
	port   int
	server *grpc.Server
	user.UnimplementedUserServer
}

func NewAdaptor(api ports.APIPort, port int) *Adapter {
	return &Adapter{
		api:  api,
		port: port,
	}
}

func (a *Adapter) Run() {
	var err error
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatalf("failed to listen on port %d, error: %v", a.port, err)
	}
	var opts []grpc.ServerOption
	opts = append(opts,
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	grpcServer := grpc.NewServer(opts...)

	a.server = grpcServer

	user.RegisterUserServer(grpcServer, a)

	if config.GetEnv() == "development" {
		reflection.Register(grpcServer)
	}

	log.Printf("starting user service on port %d...", a.port)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve grpc on port %d", a.port)
	}
}

func (a *Adapter) Stop() {
	a.server.Stop()
}
