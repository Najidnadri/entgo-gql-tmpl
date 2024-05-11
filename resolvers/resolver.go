package resolvers

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	appcore "chipin/app/core"
	"chipin/app/domain"
	app_plugin "chipin/app/plugin"
	"chipin/app/usecase"
	"chipin/ent"
	generated "chipin/resolvers/generated"

	"github.com/99designs/gqlgen/graphql"
	"github.com/rs/zerolog"
)

// Resolver is the resolver root.
type Resolver struct {
	client   *ent.Client
	logger   *zerolog.Logger
	services usecase.Service
	core     *appcore.AppCore
	plugins  *app_plugin.AppPlugin
	repos    *domain.Repository
}

// NewSchema creates a graphql executable schema.
func NewSchema(client *ent.Client, logger *zerolog.Logger, services usecase.Service, core *appcore.AppCore, plugins *app_plugin.AppPlugin, repos *domain.Repository) graphql.ExecutableSchema {
	return generated.NewExecutableSchema(generated.Config{
		Resolvers: &Resolver{
			client:   client,
			logger:   logger,
			services: services,
			core:     core,
			plugins:  plugins,
			repos:    repos,
		},
	})
}
