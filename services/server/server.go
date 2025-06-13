package server

import (
	"context"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	oMiddleware "github.com/oapi-codegen/gin-middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	svCtx "personal-project-core/context"
	"personal-project-core/middleware"
)

type Server struct {
	router     *gin.Engine
	sc         *svCtx.ServiceContext
	httpServer *http.Server
}

func NewServer(sv *svCtx.ServiceContext) *Server {
	router := gin.Default()
	return &Server{
		router: router,
		sc:     sv,
	}
}

func (s *Server) SetupRoutes() {
	s.router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to TRY API")
	})

	s.router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// serve swagger ui
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger.yaml")))
	s.router.StaticFile("/swagger.yaml", "bundled.yaml")

	// Load API specification directly from api.yaml instead of using GetSwagger()
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	swagger, err := loader.LoadFromFile("bundled.yaml")
	if err != nil {
		log.Fatalf("error loading swagger spec from file: %v", err)
	}

	// Validate the swagger document
	err = swagger.Validate(context.Background())
	if err != nil {
		log.Fatalf("error validating swagger spec: %v", err)
	}

	// register api group with swagger validator
	authMiddlewareFactory := middleware.AuthMiddleware(s.sc)
	apiPrefix := "/api/v1"
	apiGroupV1 := s.router.Group(
		apiPrefix,
		oMiddleware.OapiRequestValidatorWithOptions(swagger, &oMiddleware.Options{
			ErrorHandler: func(c *gin.Context, err string, statusCode int) {
				c.JSON(statusCode, gin.H{
					"message": "Validation failed",
					"error":   err,
				})
			},
			SilenceServersWarning: true,
			Options: openapi3filter.Options{
				AuthenticationFunc: authMiddlewareFactory,
			},
		}))

	apiGroupV1.Use()
	{

	}
}

func (s *Server) Run() {
	cfg := s.sc.Config

	s.httpServer = &http.Server{
		Addr:    cfg.Server.Port,
		Handler: s.router,
	}

}
