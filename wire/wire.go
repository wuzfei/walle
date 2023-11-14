//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"yema.dev/internal/handler"
	"yema.dev/internal/repository"
	"yema.dev/internal/server"
	"yema.dev/internal/service"
	"yema.dev/pkg/app"
	"yema.dev/pkg/helper/sid"
	"yema.dev/pkg/jwt"
	"yema.dev/pkg/repo"
	"yema.dev/pkg/server/http"
	"yema.dev/pkg/ssh"
)

var pkgSet = wire.NewSet(
	ssh.NewSSH,
	repo.NewRepos,
)

var repositorySet = wire.NewSet(
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewUserRepository,
	repository.NewEnvironmentRepository,
	repository.NewSpaceRepository,
	repository.NewServerRepository,
	repository.NewProjectRepository,
	repository.NewDeployRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
	service.NewServerService,
	service.NewEnvironmentService,
	service.NewSpaceService,
	service.NewProjectService,
	service.NewDeployService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
	handler.NewEnvironmentHandler,
	handler.NewSpaceHandler,
	handler.NewServerHandler,
	handler.NewProjectHandler,
	handler.NewDeployHandler,
	handler.NewCommonHandler,
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

func NewWire(*zap.Logger, *gorm.DB, *handler.AssetsHandler, *ssh.Config, *repo.Config, *http.Config, *jwt.Config) (*app.App, func(), error) {
	panic(wire.Build(
		pkgSet,
		repositorySet,
		serviceSet,
		handlerSet,
		serverSet,
		sid.NewSid,
		jwt.NewJwt,
		newApp,
	))
}