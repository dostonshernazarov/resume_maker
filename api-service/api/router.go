package api

import (
	// "net/http"
	"time"

	_ "github.com/dostonshernazarov/resume_maker/api-service/api/docs"
	v1 "github.com/dostonshernazarov/resume_maker/api-service/api/handlers/v1"

	"github.com/dostonshernazarov/resume_maker/api-service/api/middleware"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go.uber.org/zap"

	grpcClients "github.com/dostonshernazarov/resume_maker/api-service/internal/infrastructure/grpc_service_client"
	"github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/config"
	tokens "github.com/dostonshernazarov/resume_maker/api-service/internal/pkg/token"
	"github.com/dostonshernazarov/resume_maker/api-service/internal/usecase/app_version"
	"github.com/dostonshernazarov/resume_maker/api-service/internal/usecase/event"
	// "github.com/dostonshernazarov/resume_maker/api-service/internal/usecase/refresh_token"
)

type RouteOption struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	Service        grpcClients.ServiceClient
	JwtHandler     tokens.JwtHandler
	BrokerProducer event.BrokerProducer
	AppVersion     app_version.AppVersion
	Enforcer       *casbin.Enforcer
}

// NewRouter
// @title Welcome To CV Maker API
// @Description API for CV Maker
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewRoute(option RouteOption) *gin.Engine {

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	HandlerV1 := v1.New(&v1.HandlerV1Config{
		Config:         option.Config,
		Logger:         option.Logger,
		ContextTimeout: option.ContextTimeout,
		Service:        option.Service,
		JwtHandler:     option.JwtHandler,
		AppVersion:     option.AppVersion,
		BrokerProducer: option.BrokerProducer,
		Enforcer:       option.Enforcer,
	})

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowBrowserExtensions = true
	corsConfig.AllowMethods = []string{"*"}
	router.Use(cors.New(corsConfig))

	// router.Use(middleware.Tracing)
	router.Use(middleware.CheckCasbinPermission(option.Enforcer, *option.Config))

	router.Static("/media", "./media")
	api := router.Group("/v1")

	// USER METHODS
	api.POST("/users", HandlerV1.CreateUser)
	api.PUT("/users", HandlerV1.UpdateUser)
	api.DELETE("/users/:id", HandlerV1.DeleteUser)
	api.GET("/users/:id", HandlerV1.GetUser)
	api.GET("/users/list", HandlerV1.ListUsers)
	api.GET("/users/token", HandlerV1.GetByToken)

	// REGISTER METHODS
	api.POST("/users/register", HandlerV1.RegisterUser)
	api.GET("/users/verify", HandlerV1.Verification)
	api.POST("/users/login", HandlerV1.Login)
	api.GET("/users/set/:email", HandlerV1.ForgetPassword)
	api.GET("/users/code", HandlerV1.ForgetPasswordVerify)
	api.PUT("/users/password", HandlerV1.SetNewPassword)
	api.GET("/token/:refresh", HandlerV1.UpdateToken)

	// MEDIA
	api.POST("/media/user-photo", HandlerV1.UploadMedia)
	api.POST("/resume/resume-photo", HandlerV1.UploadResumePhoto)

	// RESUME
	api.POST("/resume/generate-resume", HandlerV1.GenerateResume)
	api.GET("/users/resume/list", HandlerV1.ListUserResume)
	api.GET("/resume/list", HandlerV1.ListResume)
	api.DELETE("/resumes/:id", HandlerV1.DeleteResume)

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return router
}
