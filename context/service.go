package context

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"personal-project-core/config"
)

type ServiceContext struct {
	Config *config.Config
	DB     *gorm.DB
	Log    *zap.Logger
}

func NewServiceContext(cfg *config.Config, log *zap.Logger, db *gorm.DB) *ServiceContext {
	return &ServiceContext{
		Config: cfg,
		DB:     db,
		Log:    log,
	}
}
