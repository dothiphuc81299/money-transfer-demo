package app

import (
	"context"
	"log"
	"money-transfer-demo/pkg/identity/config"
	"money-transfer-demo/pkg/identity/protocol/rest"
	"money-transfer-demo/pkg/identity/user/userimpl"
	"money-transfer-demo/pkg/infra/storage/postgres"
	"net/http"

	"go.uber.org/zap"
	"gorm.io/gorm"

	memberimpl "money-transfer-demo/pkg/identity/member/memberimpl"
)

type Server struct {
	Postgresdb *gorm.DB
	cfg        *config.Config
	RestServer *rest.Server
	Log        *zap.Logger
}

const serviceName = "identity"

func NewServer() (*Server, error) {
	cfg, err := config.FromEnv()
	if err != nil {
		return nil, err
	}

	postgresdb, err := postgres.New(cfg.Postgres.ConnectionString(), serviceName)
	if err != nil {
		return nil, err
	}

	paymentClient, err := getPaymentClient(cfg)
	if err != nil {
		return nil, err
	}

	memberStore := memberimpl.NewStore(postgresdb)
	memberSvc := memberimpl.NewService(memberStore, paymentClient)

	userStore := userimpl.NewStore(postgresdb)
	userSvc := userimpl.NewService(userStore)

	restServer := rest.NewServer(&rest.Dependencies{
		Postgres:  postgresdb,
		MemberSvc: memberSvc,
		UserSvc:   userSvc,
	}, cfg)

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
