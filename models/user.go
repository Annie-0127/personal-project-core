package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	BaseModel
	Name       string         `json:"name"`
	Username   string         `json:"user_name" gorm:"type:varchar(255);unique"`
	Email      string         `json:"email"`
	Password   string         `json:"-"`
	RoleIds    pq.StringArray `gorm:"type:text[]" json:"role_codes"`
	CompanyId  *uuid.UUID     `gorm:"type:uuid;index" json:"company_id"`
	StoreIds   pq.StringArray `gorm:"type:text[]" json:"store_ids"`
	SupplierId *uuid.UUID     `gorm:"type:uuid;index" json:"supplier_id"`
}
