package lp

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/DimTur/lp_learning_platform/internal/app"
	"github.com/DimTur/lp_learning_platform/internal/config"
	attstorage "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/attempts"
	channelstorage "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/channels"
	lessonstorage "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/lessons"
	pagestorage "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/pages"
	planstorage "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/plans"
	questionstorage "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/questions"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
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

			// storage, err := sqlite.New(cfg.Storage.SQLitePath)
			// if err != nil {
			// 	return err
			// }

			dsn := fmt.Sprintf(
				"postgres://%s:%s@%s:%d/%s?sslmode=disable",
				cfg.Storage.User,
				cfg.Storage.Password,
				cfg.Storage.Host,
				cfg.Storage.Port,
				cfg.Storage.DBName,
			)
			storagePool, err := pgxpool.New(ctx, dsn)
			if err != nil {
				return err
			}
			defer storagePool.Close()

			channelStorage := channelstorage.NewChannelStorage(storagePool)
			planStorage := planstorage.NewPlansStorage(storagePool)
			lessonStorage := lessonstorage.NewLessonsStorage(storagePool)
			pageStorage := pagestorage.NewPagesStorage(storagePool)
			questionStorage := questionstorage.NewQuestionsStorage(storagePool)
			attemptStorage := attstorage.NewAttemptsStorage(storagePool)

			validate := validator.New()

			application, err := app.NewApp(
				channelStorage,
				planStorage,
				lessonStorage,
				pageStorage,
				questionStorage,
				attemptStorage,
				cfg.GRPCServer.Address,
				log,
				validate,
			)
			if err != nil {
				return err
			}

			grpcCloser, err := application.GRPCSrv.Run()
			if err != nil {
				return err
			}

			log.Info("server listening:", slog.Any("port", cfg.GRPCServer.Address))
			<-ctx.Done()

			// if err := storagePool.Close(); err != nil {
			// 	log.Error("storage.Close", slog.Any("err", err))
			// }

			grpcCloser()

			return nil
		},
	}

	c.Flags().StringVar(&configPath, "config", "", "path to config")
	return c
}
