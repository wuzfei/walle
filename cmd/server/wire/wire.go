//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
	"yema.dev/internal/handler"
	"yema.dev/internal/repository"
	"yema.dev/internal/server"
	"yema.dev/internal/service"
	"yema.dev/pkg/app"
	"yema.dev/pkg/helper/sid"
	"yema.dev/pkg/jwt"
	"yema.dev/pkg/log"
	"yema.dev/pkg/server/http"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRedis,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewUserRepository,
	repository.NewEnvironmentRepository,
	repository.NewSpaceRepository,
	repository.NewServerRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
	service.NewServerService,
	service.NewEnvironmentService,
	service.NewSpaceService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
	handler.NewEnvironmentHandler,
	handler.NewSpaceHandler,
	handler.NewServerHandler,
)

var serverSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewJob,
	server.NewTask,
)

// build App
func newApp(httpServer *http.Server, job *server.Job) *app.App {
	return app.NewApp(
		app.WithServer(httpServer, job),
		app.WithName("demo-server"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {

	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		serverSet,
		sid.NewSid,
		jwt.NewJwt,
		newApp,
	))
}
