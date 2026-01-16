package app

import (
	"context"
	"net/http"
	"sync"

	"money-transfer-demo/pkg/infra/storage/postgres"
	"money-transfer-demo/pkg/payment/bankacc/bankaccimpl"
	"money-transfer-demo/pkg/payment/config"
	"money-transfer-demo/pkg/payment/deposit/depositimpl"
	"money-transfer-demo/pkg/payment/memberacc/memberaccimpl"
	"money-transfer-demo/pkg/payment/memberpayacc/memberpayaccimpl"
	"money-transfer-demo/pkg/payment/protocol/grpc"
	"money-transfer-demo/pkg/payment/protocol/rest"
	"money-transfer-demo/pkg/payment/transfer/transferimpl"
	"money-transfer-demo/pkg/payment/withdrawal"
	"money-transfer-demo/pkg/payment/withdrawal/withdrawalimpl"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Server struct {
	Postgresdb *gorm.DB
	Cfg        *config.Config
	Log        *zap.Logger

	RestServer *rest.Server
	GrpcServer *grpc.Server

	cancel context.CancelFunc
	wg     sync.WaitGroup
}

const serviceName = "payment"

func NewServer() (*Server, error) {
	cfg, err := config.FromEnv()
	if err != nil {
		return nil, err
	}

	// DB
	pg, err := postgres.New(cfg.Postgres.ConnectionString(), serviceName)
	if err != nil {
		return nil, err
	}

	// Stores
	memberAccStore := memberaccimpl.NewStore(pg)
	bankAccountStore := bankaccimpl.NewStore(pg)
	depositStore := depositimpl.NewStore(pg)
	withdrawalStore := withdrawalimpl.NewStore(pg)
	memberPayAccStore := memberpayaccimpl.NewStore(pg)
	transferStore := transferimpl.NewStore(pg)

	// Services
	memberAccSvc := memberaccimpl.NewService(memberAccStore)
	bankAccSvc := bankaccimpl.NewService(bankAccountStore)
	memberPayAccSvc := memberpayaccimpl.NewService(memberPayAccStore, memberAccSvc)
	depositSvc := depositimpl.NewService(depositStore, memberAccSvc, bankAccSvc, cfg)
	withdrawalSvc := withdrawalimpl.NewService(withdrawalStore, memberAccSvc, bankAccSvc, memberPayAccStore, cfg)
	transferSvc := transferimpl.NewService(memberAccSvc, transferStore)

	// Servers
	grpcSrv := grpc.NewServer(&grpc.Dependencies{
		MemberAccountSvc: memberAccSvc,
		Cfg:              cfg,
	})

	restSrv := rest.NewServer(&rest.Dependencies{
		MemberAccSvc:    memberAccSvc,
		BankAccSrv:      bankAccSvc,
		DepositSrv:      depositSvc,
		WithdrawalSrv:   withdrawalSvc,
		MemberPayAccSrv: memberPayAccSvc,
		TransferSrv:     transferSvc,
		Cfg:             cfg,
		Postgres:        pg,
	}, cfg)

	s := &Server{
		Postgresdb: pg.GetDB(),
		Cfg:        cfg,
		RestServer: restSrv,
		GrpcServer: grpcSrv,
		Log:        zap.L().Named("apiserver"),
	}

	// Context for background tasks
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	// Run servers & background jobs
	s.wg.Add(3)
	go s.runGRPC(ctx)
	go s.runREST(ctx)
	go s.runWithdrawalWorker(ctx, withdrawalSvc)

	return s, nil
}

func (s *Server) runGRPC(ctx context.Context) {
	defer s.wg.Done()
	if err := s.GrpcServer.Run(ctx); err != nil && err != http.ErrServerClosed {
		s.Log.Error("gRPC server failed", zap.Error(err))
	}
}

func (s *Server) runREST(ctx context.Context) {
	defer s.wg.Done()
	if err := s.RestServer.Run(ctx); err != nil && err != http.ErrServerClosed {
		s.Log.Error("REST server failed", zap.Error(err))
	}
}

func (s *Server) runWithdrawalWorker(ctx context.Context, svc withdrawal.Service) {
	defer s.wg.Done()
	svc.Run(ctx)
}

// Shutdown gracefully
func (s *Server) Shutdown(ctx context.Context) error {
	// Cancel background tasks
	if s.cancel != nil {
		s.cancel()
	}

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		s.Log.Warn("shutdown timeout")
	}

	if s.RestServer != nil {
		_ = s.RestServer.Shutdown(ctx)
	}

	sqlDB, err := s.Postgresdb.DB()
	if err == nil {
		_ = sqlDB.Close()
	}

	s.Log.Info("server shutdown complete")
	return nil
}
