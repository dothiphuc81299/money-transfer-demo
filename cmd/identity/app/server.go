package app

import (
	"context"
	"log"
	"money-transfer-demo/pkg/identity/config"
	"money-transfer-demo/pkg/identity/protocol/rest"
	"net/http"

	"gorm.io/gorm"

	"money-transfer-demo/pkg/identity/db"

	memberimpl "money-transfer-demo/pkg/identity/member/memberimpl"
)

type Server struct {
	Postgresdb *gorm.DB
	cfg        *config.Config
	RestServer *rest.Server
}

func NewServer() (*Server, error) {
	cfg, err := config.FromEnv()
	if err != nil {
		return nil, err
	}

	postgresdb, err := db.New(cfg.Postgres.ConnectionString())
	if err != nil {
		return nil, err
	}

	memberStore := memberimpl.NewStore(postgresdb)
	memberSvc := memberimpl.NewService(memberStore)

	restServer := rest.NewServer(&rest.Dependencies{
		MemberSvc: memberSvc,
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