package main

import (
	appcore "chipin/app/core"
	"chipin/app/domain"
	app_plugin "chipin/app/plugin"
	"chipin/app/usecase"
	"chipin/cmd/base"
	"chipin/cmd/playground"
	"chipin/config"
	"chipin/ent"
	"chipin/resolvers"
	"fmt"
	"net/http"

	_ "chipin/ent/runtime"

	"entgo.io/contrib/entgql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
)

func main() {
	client, repos, plugins, services, core, logger, err := base.InitAll()
	if client != nil {
		defer client.Close()
	}
	if err != nil {
		logger.Fatal().Msg(fmt.Sprintf("Failed to init: %v", err))
	}

	// start server
	if err := startGraphqlServer(":"+config.GetPort(), client, &logger, services, repos, core, plugins); err != nil {
		logger.Fatal().Msgf("failed starting graphql server: %v", err)
	}
}

func startGraphqlServer(address string, client *ent.Client, logger *zerolog.Logger, services *usecase.Service, repos *domain.Repository, core *appcore.AppCore, plugin *app_plugin.AppPlugin) error {
	srv := handler.NewDefaultServer(resolvers.NewSchema(client, logger, *services, core, plugin, repos))
	srv.Use(entgql.Transactioner{TxOpener: client})
	srv.AddTransport(transport.MultipartForm{})

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Group(func(r chi.Router) {
		//r.Use(api.TransactionalMiddleware(client))
		//r.Post("/stripe-webhook", api.POSTstripeWebhook(*services, *plugin, logger))
		//r.Post("/easyparcel-webhook", api.POSTeasyParcelWebhook(core, services, repos, client, logger))
	})

	router.Group(func(r chi.Router) {
		//r.Use(graphql.AuthMiddleware(*services, repos.UserRepo, infrastructure.UserTypeCustomer, logger))
		r.Handle("/store/graphql", srv)
		r.Handle("/store/playground1", playground.Handler("GraphiQL", "/store/graphql"))
		r.Handle("/store/playground2", playground.AltairHandler("Altair", "/store/graphql"))
		r.Handle("/store/playground3", playground.ApolloHandler("Apollo", "/store/graphql"))
	})

	router.Group(func(r chi.Router) {
		//r.Use(graphql.AuthMiddleware(*services, repos.UserRepo, infrastructure.UserTypeSeller, logger))
		r.Handle("/seller/graphql", srv)
		r.Handle("/seller/playground1", playground.Handler("GraphiQL", "/seller/graphql"))
		r.Handle("/seller/playground2", playground.AltairHandler("Altair", "/seller/graphql"))
		r.Handle("/seller/playground3", playground.ApolloHandler("Apollo", "/seller/graphql"))
	})

	router.Group(func(r chi.Router) {
		//r.Use(graphql.AuthMiddleware(*services, repos.UserRepo, infrastructure.UserTypeAdmin, logger))
		r.Handle("/admin/graphql", srv)
		r.Handle("/admin/playground1", playground.Handler("GraphiQL", "/admin/graphql"))
		r.Handle("/admin/playground2", playground.AltairHandler("Altair", "/admin/graphql"))
		r.Handle("/admin/playground3", playground.ApolloHandler("Apollo", "/admin/graphql"))
	})
	if err := http.ListenAndServe(address, router); err != nil {
		return err
	}
	return nil
}
