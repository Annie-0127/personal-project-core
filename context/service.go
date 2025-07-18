package context

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
	"personal-project-core/config"
	"personal-project-core/models"
	"personal-project-core/services/database"
	"time"
)

type ServiceContext struct {
	Config *config.Config
	DB     *gorm.DB
	Log    *zap.Logger

	UserRepo database.Repository[models.User]

	cancel   context.CancelFunc
	shutdown chan struct{}
}

func NewServiceContext(cfg *config.Config, log *zap.Logger, db *gorm.DB) *ServiceContext {
	userRepo := database.NewPostgresRepository[models.User](db)
	return &ServiceContext{
		Config:   cfg,
		DB:       db,
		Log:      log,
		UserRepo: userRepo,
	}
}

func (sc *ServiceContext) Shutdown() {
	sc.cancel()

	// gracefully shutdown the database connection
	sqlDB, err := sc.DB.DB()
	if err == nil {
		// wait for the database connection to be closed
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		sqlDB.Close()
	}

	close(sc.shutdown)
}

func (s *ServiceContext) Wait() {
	<-s.shutdown
	log.Println("Service context shutdown complete")
}
