//go:build wireinject
// +build wireinject

package wire

import (
	"context"
	"github.com/google/wire"
	"go.uber.org/zap"
	"yema.dev/internal/handler"
	"yema.dev/internal/repository"
	"yema.dev/internal/server"
	"yema.dev/internal/service"
	"yema.dev/pkg/app"
	"yema.dev/pkg/db"
	"yema.dev/pkg/helper/sid"
	"yema.dev/pkg/jwt"
	"yema.dev/pkg/repo"
	"yema.dev/pkg/server/http"
	"yema.dev/pkg/ssh"
)

var pkgSet = wire.NewSet(
	db.NewDB,
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
	repository.NewMemberRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
	service.NewServerService,
	service.NewEnvironmentService,
	service.NewSpaceService,
	service.NewProjectService,
	service.NewDeployService,
	service.NewMemberService,
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
	handler.NewMemberHandler,
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

func NewWire(context.Context, *zap.Logger, *handler.AssetsHandler, *db.Config, *ssh.Config, *repo.Config, *http.Config, *jwt.Config) (*app.App, func(), error) {
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
