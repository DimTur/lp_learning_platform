package lp

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/DimTur/lp_learning_platform/internal/app"
	"github.com/DimTur/lp_learning_platform/internal/config"
	sqlite "github.com/DimTur/lp_learning_platform/internal/services/storage/sqlite/channel"
	"github.com/spf13/cobra"
)

func NewServeCmd() *cobra.Command {
	var configPath string

	c := &cobra.Command{
		Use:     "serve",
		Aliases: []string{"s"},
		Short:   "Start gRPS LP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

			ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
			defer cancel()

			cfg, err := config.Parse(configPath)
			if err != nil {
				return err
			}

			storage, err := sqlite.New(cfg.Storage.SQLitePath)
			if err != nil {
				return err
			}

			application, err := app.NewApp(storage, cfg.GRPCServer.Address, log)
			if err != nil {
				return err
			}

			grpcCloser, err := application.GRPCSrv.Run()
			if err != nil {
				return err
			}

			log.Info("server listening:", slog.Any("port", cfg.GRPCServer.Address))
			<-ctx.Done()

			if err := storage.Close(); err != nil {
				log.Error("storage.Close", slog.Any("err", err))
			}

			grpcCloser()

			return nil
		},
	}

	c.Flags().StringVar(&configPath, "config", "", "path to config")
	return c
}
