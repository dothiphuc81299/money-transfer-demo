package grpc

import (
	"context"
	payment "money-transfer-demo/pkg/apis/payment"
	"money-transfer-demo/pkg/payment/config"
	"money-transfer-demo/pkg/payment/memberacc"
	"net"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	payment.UnimplementedPaymentServer
	dependencies *Dependencies
	log          *zap.Logger
}

type Dependencies struct {
	MemberAccountSvc memberacc.Service
	Cfg              *config.Config
}

func NewServer(deps *Dependencies) *Server {
	return &Server{
		dependencies: deps,
		log:          zap.L().Named("grpc server"),
	}
}

func (s *Server) Run(ctx context.Context) error {
	var err error
	listen, err := net.Listen("tcp", ":"+s.dependencies.Cfg.Server.GRPCPort)
	if err != nil {
		s.log.Error("Starting GRPC server failed: ", zap.Error(err))
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
		s.log.Info("Shutting down RPC server")
	}()

	s.log.Info("Starting GRPC server...")
	if err := server.Serve(listen); err != nil {
		s.log.Error("Starting GRPC server failed: ", zap.Error(err))
	}

	return err
}
