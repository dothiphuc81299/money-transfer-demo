package app

import (
	"context"
	"log"
	"money-transfer-demo/pkg/infra/storage/postgres"
	"money-transfer-demo/pkg/payment/bankacc/bankaccimpl"
	"money-transfer-demo/pkg/payment/config"
	"money-transfer-demo/pkg/payment/memberacc/memberaccimpl"
	"money-transfer-demo/pkg/payment/protocol/grpc"
	"money-transfer-demo/pkg/payment/protocol/rest"
	"net/http"

	"gorm.io/gorm"
)

type Server struct {
	Postgresdb *gorm.DB
	cfg        *config.Config
	RestServer *rest.Server
	GrpcServer *grpc.Server
}

const serviceName string = "payment"

func NewServer() (*Server, error) {
	cfg, err := config.FromEnv()
	if err != nil {
		return nil, err
	}

	postgresdb, err := postgres.New(cfg.Postgres.ConnectionString(), serviceName)
	if err != nil {
		return nil, err
	}

	memberAccStore := memberaccimpl.NewStore(postgresdb)
	memberAccSvc := memberaccimpl.NewService(memberAccStore)
	bankAccountStore := bankaccimpl.NewStore(postgresdb)
	bankAccSrv := bankaccimpl.NewService(bankAccountStore)

	grpcServer := grpc.NewServer(&grpc.Dependencies{
		MemberAccountSvc: memberAccSvc,
		Cfg:              cfg,
	})

	restServer := rest.NewServer(&rest.Dependencies{
		MemberAccSvc: memberAccSvc,
		BankAccSrv:   bankAccSrv,
		Cfg:          cfg,
	}, cfg)

	go func() {
		if err := grpcServer.Run(context.Background()); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	go func() {
		if err := restServer.Run(context.Background()); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	return &Server{
		Postgresdb: postgresdb.GetDB(),
		cfg:        cfg,
		RestServer: restServer,
	}, nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.RestServer != nil {
		return s.RestServer.Shutdown(ctx)
	}
	return nil
}
