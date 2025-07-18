package database

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Pagination[T any] struct {
	Page  int
	Pages int
	Total int64
	Limit int
	Items []T
}

type Repository[T any] interface {
	GetDB() *gorm.DB
	WithTx(tx *gorm.DB) Repository[T]

	First(ctx context.Context, id uuid.UUID) (*T, error)
	FindOne(ctx context.Context, conditions map[string]any, selectFields []string) (*T, error)
	Create(ctx context.Context, entity T) error
	Paginate(ctx context.Context, conditions map[string]any, selectFields []string, limit, offset int, sort *string) Pagination[T]
	Update(ctx context.Context, id uuid.UUID, updates map[string]any) error
	Delete(ctx context.Context, id uuid.UUID) error
}
