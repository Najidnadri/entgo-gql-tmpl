package base

import (
	appcore "chipin/app/core"
	"chipin/app/domain"
	app_plugin "chipin/app/plugin"
	"chipin/app/usecase"
	"chipin/config"
	"chipin/ent"
	"chipin/utils/logger"
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "chipin/ent/runtime"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
)

func InitAll() (*ent.Client, *domain.Repository, *app_plugin.AppPlugin, *usecase.Service, *appcore.AppCore, zerolog.Logger, error) {
	config.InitConfig()
	logger, err := logger.InitLogger(config.GetDebug())
	if err != nil {
		panic("err initializing zerolog logger: " + err.Error())
	}

	logger.Info().Msg("Starting chipin backend")
	logger.Info().Msg(fmt.Sprintf("debug mode: %v", config.GetDebug()))

	ctx := context.Background()

	// init database
	client, err := Open(config.GetDbUrl())
	if err != nil {
		logger.Fatal().Msg(fmt.Sprintf("Failed to connect to postgres: %v", err))
	}

	// init repositories & services
	repositories := InitRepositories(client, &logger)
	plugins := InitPlugins(ctx, repositories, client, &logger)
	services := InitServices(repositories, plugins, client, &logger)
	core := InitCores(repositories, services, plugins, &logger)

	return client, repositories, plugins, services, core, logger, nil
}

func Open(databaseUrl string) (*ent.Client, error) {
	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(50)
	db.SetConnMaxLifetime(time.Hour)

	// Create an ent.Driver from `db`.
	drv := entsql.OpenDB(dialect.Postgres, db)
	// drv2 := entcache.NewDriver(
	// 	drv,
	// 	entcache.TTL(time.Second),
	// 	entcache.Levels(entcache.NewLRU(128)),
	// )
	return ent.NewClient(ent.Driver(drv)), nil
}

func InitRepositories(client *ent.Client, logger *zerolog.Logger) *domain.Repository {

	return &domain.Repository{}
}

func InitServices(repos *domain.Repository, plugins *app_plugin.AppPlugin, client *ent.Client, logger *zerolog.Logger) *usecase.Service {

	return &usecase.Service{}
}

func InitCores(repos *domain.Repository, services *usecase.Service, plugins *app_plugin.AppPlugin, logger *zerolog.Logger) *appcore.AppCore {

	return &appcore.AppCore{}
}

func InitPlugins(ctx context.Context, repos *domain.Repository, client *ent.Client, logger *zerolog.Logger) *app_plugin.AppPlugin {

	return &app_plugin.AppPlugin{}
}
