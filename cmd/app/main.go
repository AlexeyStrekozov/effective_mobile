package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/AlexeyStrekozov/effective_mobile/internal/core/config"
	core_logger "github.com/AlexeyStrekozov/effective_mobile/internal/core/logger"
	core_pgx_pool "github.com/AlexeyStrekozov/effective_mobile/internal/core/repositorty/postgres/pool/pgx"
	core_http_middleware "github.com/AlexeyStrekozov/effective_mobile/internal/core/transport/http/middleware"
	core_http_server "github.com/AlexeyStrekozov/effective_mobile/internal/core/transport/http/server"
	statistics_postgres_repository "github.com/AlexeyStrekozov/effective_mobile/internal/features/statistics/repository/postgres"
	statistics_service "github.com/AlexeyStrekozov/effective_mobile/internal/features/statistics/service"
	statistics_transport_http "github.com/AlexeyStrekozov/effective_mobile/internal/features/statistics/transport/http"
	subscriptions_postgres_repository "github.com/AlexeyStrekozov/effective_mobile/internal/features/subscriptions/repository/postgres"
	subscriptions_service "github.com/AlexeyStrekozov/effective_mobile/internal/features/subscriptions/service"
	subscriptions_transport_http "github.com/AlexeyStrekozov/effective_mobile/internal/features/subscriptions/transport/http"
	"go.uber.org/zap"

	_ "github.com/AlexeyStrekozov/effective_mobile/docs"
)

func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("application time zone", zap.Any("zone", time.Local))

	logger.Debug("Init postgres connection pool")
	pool, err := core_pgx_pool.NewPoll(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)

	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("init feature", zap.String("feature", "subscriptions"))
	subscriptionsRepository := subscriptions_postgres_repository.NewSubscriptionsRepository(pool)
	subscriptionsService := subscriptions_service.NewSubscriptionsService(subscriptionsRepository)
	subscriptionsTransportHTTP := subscriptions_transport_http.NewSubscriptionsHTTPHandler(subscriptionsService)

	logger.Debug("init feature", zap.String("feature", "statistics"))
	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPHandler(statisticsService)

	logger.Debug("initializing HTTP server")
	httpConfig := core_http_server.NewConfigMust()
	httpServer := core_http_server.NewHTTPServer(
		httpConfig,
		logger,
		core_http_middleware.CORS(httpConfig.AllowedOrigins),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	apiVersionRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouterV1.AddRoutes(subscriptionsTransportHTTP.Routes()...)
	apiVersionRouterV1.AddRoutes(statisticsTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouterV1)

	httpServer.RegisterSwagger()

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
