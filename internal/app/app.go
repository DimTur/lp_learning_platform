package app

import (
	"log/slog"

	grpcapp "github.com/DimTur/lp_learning_platform/internal/app/grpc"
	"github.com/DimTur/lp_learning_platform/internal/services/channel"
	postgresql "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/channels"
)

type App struct {
	GRPCSrv *grpcapp.Server
}

func NewApp(
	storage *postgresql.ChannelPostgresStorage,
	grpcAddr string,
	logger *slog.Logger,
) (*App, error) {
	lpGRPCHandlers := channel.New(
		logger,
		storage,
		storage,
		storage,
	)

	grpcServer, err := grpcapp.NewGRPCServer(
		grpcAddr,
		lpGRPCHandlers,
		logger,
	)
	if err != nil {
		return nil, err
	}

	return &App{
		GRPCSrv: grpcServer,
	}, nil
}
