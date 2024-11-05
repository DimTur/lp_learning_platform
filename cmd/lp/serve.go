package lp

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/DimTur/lp_learning_platform/internal/app"
	"github.com/DimTur/lp_learning_platform/internal/app/consumers"
	"github.com/DimTur/lp_learning_platform/internal/config"
	"github.com/DimTur/lp_learning_platform/internal/services/rabbitmq"
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
			var wg sync.WaitGroup

			cfg, err := config.Parse(configPath)
			if err != nil {
				return err
			}

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

			// Init RabbitMQ
			rmqUrl := fmt.Sprintf(
				"amqp://%s:%s@%s:%d/",
				cfg.RabbitMQ.UserName,
				cfg.RabbitMQ.Password,
				cfg.RabbitMQ.Host,
				cfg.RabbitMQ.Port,
			)
			rmq, err := rabbitmq.NewClient(rmqUrl)
			if err != nil {
				log.Error("failed init rabbit mq", slog.Any("err", err))
			}

			// Declare Share exchange
			if err := rmq.DeclareExchange(
				cfg.RabbitMQ.ShareExchange.Name,
				cfg.RabbitMQ.ShareExchange.Kind,
				cfg.RabbitMQ.ShareExchange.Durable,
				cfg.RabbitMQ.ShareExchange.AutoDeleted,
				cfg.RabbitMQ.ShareExchange.Internal,
				cfg.RabbitMQ.ShareExchange.NoWait,
				cfg.RabbitMQ.ShareExchange.Args.ToMap(),
			); err != nil {
				log.Error("failed to declare Share exchange", slog.Any("err", err))
			}

			// Declare Channel Queue
			if _, err := rmq.DeclareQueue(
				cfg.RabbitMQ.Channel.ChannelQueue.Name,
				cfg.RabbitMQ.Channel.ChannelQueue.Durable,
				cfg.RabbitMQ.Channel.ChannelQueue.AutoDeleted,
				cfg.RabbitMQ.Channel.ChannelQueue.Exclusive,
				cfg.RabbitMQ.Channel.ChannelQueue.NoWait,
				cfg.RabbitMQ.Channel.ChannelQueue.Args.ToMap(),
			); err != nil {
				log.Error("failed to declare Channel queue", slog.Any("err", err))
			}

			// Bind Channel queue to Share exchange
			if err := rmq.BindQueueToExchange(
				cfg.RabbitMQ.Channel.ChannelQueue.Name,
				cfg.RabbitMQ.ShareExchange.Name,
				cfg.RabbitMQ.Channel.ChannelRoutingKey,
			); err != nil {
				log.Error("failed to bind Channel queue", slog.Any("err", err))
			}

			// Declare Plan Queue
			if _, err := rmq.DeclareQueue(
				cfg.RabbitMQ.Plan.PlanQueue.Name,
				cfg.RabbitMQ.Plan.PlanQueue.Durable,
				cfg.RabbitMQ.Plan.PlanQueue.AutoDeleted,
				cfg.RabbitMQ.Plan.PlanQueue.Exclusive,
				cfg.RabbitMQ.Plan.PlanQueue.NoWait,
				cfg.RabbitMQ.Plan.PlanQueue.Args.ToMap(),
			); err != nil {
				log.Error("failed to declare Plan queue", slog.Any("err", err))
			}

			// Bind Plan queue to Share exchange
			if err := rmq.BindQueueToExchange(
				cfg.RabbitMQ.Plan.PlanQueue.Name,
				cfg.RabbitMQ.ShareExchange.Name,
				cfg.RabbitMQ.Plan.PlanRoutingKey,
			); err != nil {
				log.Error("failed to bind Plan queue", slog.Any("err", err))
			}

			application, err := app.NewApp(
				channelStorage,
				planStorage,
				lessonStorage,
				pageStorage,
				questionStorage,
				attemptStorage,
				rmq,
				rmq,
				cfg.GRPCServer.Address,
				log,
				validate,
			)
			if err != nil {
				return err
			}

			// Start sharing channels with learning groups consumer
			channelsConsumer := consumers.NewConsumeChannel(rmq, channelStorage, log)
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := channelsConsumer.Start(
					ctx,
					cfg.RabbitMQ.Channel.ChannelConsumer.Queue,
					cfg.RabbitMQ.Channel.ChannelConsumer.Consumer,
					cfg.RabbitMQ.Channel.ChannelConsumer.AutoAck,
					cfg.RabbitMQ.Channel.ChannelConsumer.Exclusive,
					cfg.RabbitMQ.Channel.ChannelConsumer.NoLocal,
					cfg.RabbitMQ.Channel.ChannelConsumer.NoWait,
					cfg.RabbitMQ.Channel.ChannelConsumer.ConsumerArgs.ToMap(),
				); err != nil {
					log.Error("failed to start share channels consumer", slog.Any("err", err))
				}
			}()

			// Start sharing plans with users consumer
			plansConsumer := consumers.NewConsumePlan(rmq, planStorage, log)
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := plansConsumer.Start(
					ctx,
					cfg.RabbitMQ.Plan.PlanConsumer.Queue,
					cfg.RabbitMQ.Plan.PlanConsumer.Consumer,
					cfg.RabbitMQ.Plan.PlanConsumer.AutoAck,
					cfg.RabbitMQ.Plan.PlanConsumer.Exclusive,
					cfg.RabbitMQ.Plan.PlanConsumer.NoLocal,
					cfg.RabbitMQ.Plan.PlanConsumer.NoWait,
					cfg.RabbitMQ.Plan.PlanConsumer.ConsumerArgs.ToMap(),
				); err != nil {
					log.Error("failed to start share plans consumer", slog.Any("err", err))
				}
			}()

			grpcCloser, err := application.GRPCSrv.Run()
			if err != nil {
				return err
			}

			log.Info("server listening:", slog.Any("port", cfg.GRPCServer.Address))
			<-ctx.Done()
			wg.Wait()

			rmq.Close()
			grpcCloser()

			return nil
		},
	}

	c.Flags().StringVar(&configPath, "config", "", "path to config")
	return c
}
