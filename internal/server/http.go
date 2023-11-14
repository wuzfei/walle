package server

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"yema.dev/internal/handler"
	"yema.dev/internal/middleware"
	"yema.dev/internal/model"
	"yema.dev/pkg/jwt"
	"yema.dev/pkg/server/http"
)

func NewHTTPServer(
	log *zap.Logger,
	jwt *jwt.JWT,
	conf *http.Config,
	assetsHandler *handler.AssetsHandler,
	userHandler *handler.UserHandler,
	spaceHandler *handler.SpaceHandler,
	serverHandler *handler.ServerHandler,
	environmentHandler *handler.EnvironmentHandler,
	projectHandler *handler.ProjectHandler,
	deployHandler *handler.DeployHandler,
	commonHandler *handler.CommonHandler,
) *http.Server {
	gin.SetMode(gin.DebugMode)
	s := http.NewServer(
		gin.Default(),
		log,
		conf,
	)
	//注册静态路由
	assetsHandler.Register(s.Engine)

	apiGroup := s.Group("/api",
		middleware.CORSMiddleware(),
		middleware.RequestLogMiddleware(log.Named("middleware")),
	)
	// No route group has permission
	apiGroup.POST("/login", userHandler.Login)
	apiGroup.POST("/refresh_token", userHandler.RefreshToken)

	//需要登陆
	authRouter := apiGroup.Group("", middleware.Auth(jwt, log))
	{
		// No route group has permission
		authRouter.POST("/logout", userHandler.Logout)
		authRouter.GET("/user_info", userHandler.Profile)

		superPermRouter := authRouter.Group("", middleware.Permission(userService, model.RoleSuper))
		ownerPermRouter := authRouter.Group("", middleware.Permission(userService, model.RoleOwner))
		masterPermRouter := authRouter.Group("", middleware.Permission(userService, model.RoleMaster))

		// Non-strict permission routing group
		noStrictAuthRouter := apiGroup.Group("/").Use(middleware.NoStrictAuth(jwt, logger))
		{
			noStrictAuthRouter.GET("/user", userHandler.GetProfile)
		}

		// Strict permission routing group
		strictAuthRouter := apiGroup.Group("/").Use(middleware.StrictAuth(jwt, logger))
		{
			strictAuthRouter.PUT("/user", userHandler.UpdateProfile)
		}
	}

	return s
}
