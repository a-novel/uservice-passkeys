package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"buf.build/gen/go/a-novel/proto/grpc/go/passkeys/v1/passkeysv1grpc"

	"github.com/a-novel/golib/database"
	anovelgrpc "github.com/a-novel/golib/grpc"
	"github.com/a-novel/golib/loggers"
	"github.com/a-novel/golib/loggers/adapters"
	"github.com/a-novel/golib/loggers/formatters"

	"github.com/a-novel/uservice-passkeys/config"
	"github.com/a-novel/uservice-passkeys/migrations"
	"github.com/a-novel/uservice-passkeys/pkg/dao"
	"github.com/a-novel/uservice-passkeys/pkg/handlers"
	"github.com/a-novel/uservice-passkeys/pkg/services"
)

var rpcServices = []grpc.ServiceDesc{
	healthpb.Health_ServiceDesc,
	passkeysv1grpc.CreateService_ServiceDesc,
	passkeysv1grpc.DeleteService_ServiceDesc,
	passkeysv1grpc.GetService_ServiceDesc,
	passkeysv1grpc.UpdateService_ServiceDesc,
}

func getDepsCheck(database *bun.DB) *anovelgrpc.DepsCheck {
	return &anovelgrpc.DepsCheck{
		Dependencies: anovelgrpc.DepCheckCallbacks{
			"postgres": database.Ping,
		},
		Services: anovelgrpc.DepCheckServices{
			"create": {"postgres"},
			"delete": {"postgres"},
			"get":    {"postgres"},
			"update": {"postgres"},
		},
	}
}

func main() {
	logger := config.Logger.Formatter

	loader := formatters.NewLoader(
		fmt.Sprintf("Acquiring database connection at %s...", config.App.Postgres.DSN),
		spinner.Meter,
	)
	logger.Log(loader, loggers.LogLevelInfo)

	postgresDB, closePostgresDB, err := database.OpenDB(config.App.Postgres.DSN)
	if err != nil {
		logger.Log(formatters.NewError(err, "open database conn"), loggers.LogLevelFatal)
	}
	defer closePostgresDB()

	logger.Log(
		loader.SetDescription("Database connection successfully acquired.").SetCompleted(),
		loggers.LogLevelInfo,
	)

	if err := database.Migrate(postgresDB, migrations.SQLMigrations, logger); err != nil {
		logger.Log(formatters.NewError(err, "migrate database"), loggers.LogLevelFatal)
	}

	loader = formatters.NewLoader("Setup services...", spinner.Meter)
	logger.Log(loader, loggers.LogLevelInfo)

	grpcReporter := adapters.NewGRPC(logger)

	createPasskeyDAO := dao.NewCreatePasskey(postgresDB)
	deletePasskeyDAO := dao.NewDeletePasskey(postgresDB)
	getPasskeyDAO := dao.NewGetPasskey(postgresDB)
	updatePasskeyDAO := dao.NewUpdatePasskey(postgresDB)

	createPasskeyService := services.NewCreatePasskey(createPasskeyDAO)
	deletePasskeyService := services.NewDeletePasskey(deletePasskeyDAO)
	getPasskeyService := services.NewGetPasskey(getPasskeyDAO)
	updatePasskeyService := services.NewUpdatePasskey(updatePasskeyDAO)

	createPasskeyHandler := handlers.NewCreatePasskey(createPasskeyService, grpcReporter)
	deletePasskeyHandler := handlers.NewDeletePasskey(deletePasskeyService, grpcReporter)
	getPasskeyHandler := handlers.NewGetPasskey(getPasskeyService, grpcReporter)
	updatePasskeyHandler := handlers.NewUpdatePasskey(updatePasskeyService, grpcReporter)

	logger.Log(loader.SetDescription("Services successfully setup.").SetCompleted(), loggers.LogLevelInfo)

	listener, server, err := anovelgrpc.StartServer(config.App.Server.Port)
	if err != nil {
		logger.Log(formatters.NewError(err, "start server"), loggers.LogLevelFatal)
	}
	defer anovelgrpc.CloseServer(listener, server)

	reflection.Register(server)
	healthpb.RegisterHealthServer(server, anovelgrpc.NewHealthServer(getDepsCheck(postgresDB), time.Minute))
	passkeysv1grpc.RegisterCreateServiceServer(server, createPasskeyHandler)
	passkeysv1grpc.RegisterDeleteServiceServer(server, deletePasskeyHandler)
	passkeysv1grpc.RegisterGetServiceServer(server, getPasskeyHandler)
	passkeysv1grpc.RegisterUpdateServiceServer(server, updatePasskeyHandler)

	report := formatters.NewDiscoverGRPC(rpcServices, config.App.Server.Port)
	logger.Log(report, loggers.LogLevelInfo)

	if err := server.Serve(listener); err != nil {
		logger.Log(formatters.NewError(err, "serve"), loggers.LogLevelFatal)
	}
}
