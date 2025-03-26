package grpc

import (
	"context"
	"log"
	payment "money-transfer-demo/pkg/apis/payment"
	"money-transfer-demo/pkg/payment/config"
	"money-transfer-demo/pkg/payment/memberacc"
	"net"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type Server struct {
	payment.UnimplementedPaymentServer
	dependencies *Dependencies
}

type Dependencies struct {
	MemberAccountSvc memberacc.Service
	Cfg              *config.Config
}

func NewServer(deps *Dependencies) *Server {
	return &Server{
		dependencies: deps,
	}
}

func (s *Server) Run(ctx context.Context) error {
	var err error
	listen, err := net.Listen("tcp", ":"+s.dependencies.Cfg.Server.GRPCPort)
	if err != nil {
		log.Printf("Failed to create listener for grpc endpoint: %s", err.Error())
		return err
	}

	opts := []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	}

	server := grpc.NewServer(opts...)

	payment.RegisterPaymentServer(server, s)

	go func() {
		<-ctx.Done()
		server.GracefulStop()
		log.Println("Shutting down RPC server")
	}()

	log.Printf("GRPC server started on port %s", s.dependencies.Cfg.Server.GRPCPort)
	if err := server.Serve(listen); err != nil {
		log.Printf("Starting GRPC server failed: %s", err.Error())
	}

	return err
}
