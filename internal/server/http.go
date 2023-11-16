package server

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"yema.dev/internal/handler"
	"yema.dev/internal/middleware"
	"yema.dev/internal/model"
	"yema.dev/internal/repository"
	"yema.dev/pkg/jwt"
	"yema.dev/pkg/server/http"
)

func NewHTTPServer(
	log *zap.Logger,
	jwt *jwt.JWT,
	conf *http.Config,
	assetsHandler *handler.AssetsHandler,
	memberRepo repository.MemberRepository,
	userHandler *handler.UserHandler,
	memberHandler *handler.MemberHandler,
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
	//公共信息
	apiGroup.GET("/version", commonHandler.Version)
	apiGroup.GET("/system", commonHandler.System)

	//需要登陆
	authRouter := apiGroup.Group("", middleware.Auth(jwt, log))
	{
		// No route group has permission
		authRouter.POST("/logout", userHandler.Logout)
		authRouter.GET("/user_info", userHandler.Profile)

		superPermRouter := authRouter.Group("", middleware.Permission(memberRepo, model.RoleSuper))
		ownerPermRouter := authRouter.Group("", middleware.Permission(memberRepo, model.RoleOwner))
		masterPermRouter := authRouter.Group("", middleware.Permission(memberRepo, model.RoleMaster))

		//用户管理
		{
			superPermRouter.GET("/user", userHandler.List)
			superPermRouter.POST("/user", userHandler.Create)
			superPermRouter.DELETE("/user/:id", userHandler.Delete)
			superPermRouter.PUT("/user", userHandler.Update)
			superPermRouter.GET("/user/options", userHandler.Options)
		}

		//成员管理
		{
			ownerPermRouter.GET("/member", memberHandler.List)
			ownerPermRouter.POST("/member", memberHandler.Store)
			ownerPermRouter.DELETE("/member/:id", memberHandler.Delete)
		}

		//空间管理, super访问权限
		{
			superPermRouter.GET("/space", spaceHandler.List)
			superPermRouter.POST("/space", spaceHandler.Create)
			superPermRouter.DELETE("/space/:id", spaceHandler.Delete)
			superPermRouter.PUT("/space", spaceHandler.Update)
		}

		//服务器管理
		{
			ownerPermRouter.GET("/server", serverHandler.List)
			ownerPermRouter.POST("/server", serverHandler.Create)
			ownerPermRouter.DELETE("/server/:id", serverHandler.Delete)
			ownerPermRouter.PUT("/server", serverHandler.Update)
			//校验连接
			ownerPermRouter.POST("/server/:id/check", serverHandler.Check)
			//设置免登陆
			ownerPermRouter.POST("/server/set_authorized", serverHandler.SetAuthorized)
			//websocket 连接终端
			ownerPermRouter.GET("/server/:id/terminal", serverHandler.Terminal)
		}

		//环境管理
		{
			masterPermRouter.GET("/environment", environmentHandler.List)
			masterPermRouter.POST("/environment", environmentHandler.Create)
			masterPermRouter.DELETE("/environment/:id", environmentHandler.Delete)
			masterPermRouter.PUT("/environment", environmentHandler.Update)
			masterPermRouter.GET("/environment/options", environmentHandler.Options)
		}

		//项目管理
		{
			masterPermRouter.GET("/project", projectHandler.List)
			masterPermRouter.POST("/project", projectHandler.Create)
			masterPermRouter.DELETE("/project/:id", projectHandler.Delete)
			masterPermRouter.GET("/project/:id", projectHandler.Detail)
			masterPermRouter.PUT("/project", projectHandler.Update)
			masterPermRouter.GET("/project/options", projectHandler.Options)
			//项目检测 websocket
			masterPermRouter.GET("/project/:id/detection", projectHandler.Detection)
			masterPermRouter.GET("/project/:id/branches", projectHandler.Branches)
			masterPermRouter.GET("/project/:id/tags", projectHandler.Tags)
			masterPermRouter.GET("/project/:id/commits", projectHandler.Commits)
		}

		//部署管理
		{
			masterPermRouter.GET("/deploy", deployHandler.List)
			masterPermRouter.GET("/deploy/:id", deployHandler.Detail)
			masterPermRouter.POST("/deploy", deployHandler.Create)
			//审核
			masterPermRouter.POST("/deploy/:id/audit", deployHandler.Audit)
			//发布
			masterPermRouter.GET("/deploy/:id/release", deployHandler.Release)
			//发布
			masterPermRouter.GET("/deploy/:id/stop_release", deployHandler.Stop)
			//发布
			masterPermRouter.GET("/deploy/:id/rollback", deployHandler.Rollback)
			//websocket, 部署日志, 将整个部署过程日志输出
			//masterPermRouter.GET("/deploy/:id/console", deployHandler.Console)
		}
	}

	return s
}
