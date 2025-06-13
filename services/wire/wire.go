//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"personal-project-core/config"
	"personal-project-core/context"
	"personal-project-core/services/database"
	"personal-project-core/services/log"
	"personal-project-core/services/server"
)

func InitializeServer() (*server.Server, error) {
	wire.Build(
		config.LoadConfig,
		database.NewPostgresConnection,
		log.NewZapLogger,
		context.NewServiceContext,
		server.NewServer,
	)

	return &server.Server{}, nil
}
