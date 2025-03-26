package app

import (
	"context"
	"log"
	"money-transfer-demo/pkg/identity/config"
	"money-transfer-demo/pkg/identity/protocol/rest"

	"gorm.io/gorm"

	"money-transfer-demo/pkg/identity/db"

	memberimpl "money-transfer-demo/pkg/identity/member/memberimpl"
)

type Server struct {
	Postgresdb *gorm.DB
	cfg        *config.Config
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

	err = restServer.Run(context.Background())
	if err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
		return nil, err
	}

	return &Server{
		Postgresdb: postgresdb.GetDB(),
		cfg:        cfg,
	}, nil

}
