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
	"os"
	"os/signal"
	"personal-project-core/api/auth"
	"personal-project-core/api/user"
	svCtx "personal-project-core/context"
	"personal-project-core/middleware"
	"personal-project-core/models"
	"syscall"
	"time"
)

type Server struct {
	router     *gin.Engine
	sc         *svCtx.ServiceContext
	httpServer *http.Server
}

func NewServer(sv *svCtx.ServiceContext) *Server {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	return &Server{
		router: router,
		sc:     sv,
	}
}

func (s *Server) SetupRoutes() {
	s.router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to Personal Project API")
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
		authServer := auth.NewAuthServer(s.sc)
		authHandler := auth.NewStrictHandler(authServer, nil)
		auth.RegisterHandlersWithOptions(apiGroupV1, authHandler, auth.GinServerOptions{})

		userServer := user.NewUserServer(s.sc)
		userHandler := user.NewStrictHandler(userServer, nil)
		user.RegisterHandlersWithOptions(apiGroupV1, userHandler, user.GinServerOptions{})
	}
}

func (s *Server) Run() error {
	cfg := s.sc.Config

	s.httpServer = &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: s.router,
	}
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Fatal("Server forced to shutdown:", err)
		}

		s.sc.Shutdown()
	}()

	log.Printf("Server is running on %s", cfg.Server.Port)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	// service context wait for shutdown
	s.sc.Wait()

	return nil
}

func (s *Server) AutoMigrate() {
	err := s.sc.DB.AutoMigrate(
		&models.User{})

	if err != nil {
		log.Fatalf("error auto migrating models: %v", err)
	}
}
