package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	Id        uuid.UUID       `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ClientId  uuid.UUID       `json:"client_id" gorm:"type:uuid;index;not null"`
	IsActive  bool            `json:"is_active" gorm:"type:bool;not null;default:false"`
	CreatedAt time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedBy uuid.UUID       `json:"created_by" gorm:"type:uuid;"`
	UpdatedBy uuid.UUID       `json:"updated_by" gorm:"type:uuid;"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
