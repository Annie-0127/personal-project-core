package database

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"math"
	"strings"
)

type PostgresRepository[T any] struct {
	db *gorm.DB
}

func NewPostgresRepository[T any](db *gorm.DB) Repository[T] {
	return &PostgresRepository[T]{
		db: db,
	}
}

func (r *PostgresRepository[T]) GetDB() *gorm.DB {
	return r.db
}

func (r *PostgresRepository[T]) WithTx(tx *gorm.DB) Repository[T] {
	return &PostgresRepository[T]{
		db: tx,
	}
}

func (r *PostgresRepository[T]) FindOne(ctx context.Context, conditions map[string]any, selectFields []string) (*T, error) {
	var entity T

	query := r.db.WithContext(ctx)
	if len(selectFields) > 0 {
		query = query.Select(selectFields)
	}
	if len(conditions) > 0 {
		query = applyCondition(query, conditions)
	}

	err := query.First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func applyCondition(query *gorm.DB, conditions map[string]any) *gorm.DB {
	if len(conditions) > 0 {
		for key, value := range conditions {
			cleanKey := key
			if cleanKey == "" || cleanKey == "OR" {
				continue
			}
			switch {
			case strings.HasSuffix(key, " IN"):
				cleanKey = strings.TrimSuffix(key, " IN")
				query = query.Where(cleanKey+" IN ?", value)
			case strings.HasSuffix(key, " NOT IN"):
				cleanKey = strings.TrimSuffix(key, " NOT IN")
				query = query.Where(cleanKey+" NOT IN ?", value)
			case strings.HasSuffix(key, " GREATER") || strings.HasSuffix(key, " >"):
				cleanKey = strings.TrimSuffix(key, " GREATER")
				cleanKey = strings.TrimSuffix(cleanKey, " >")
				query = query.Where(cleanKey+" > ?", value)
			case strings.HasSuffix(key, " LESS") || strings.HasSuffix(key, " <"):
				cleanKey = strings.TrimSuffix(key, " LESS")
				cleanKey = strings.TrimSuffix(cleanKey, " <")
				query = query.Where(cleanKey+" < ?", value)
			case strings.HasSuffix(key, " GREATER THAN") || strings.HasSuffix(key, " >="):
				cleanKey = strings.TrimSuffix(key, " GREATER THAN")
				cleanKey = strings.TrimSuffix(cleanKey, " >=")
				query = query.Where(cleanKey+" >= ?", value)
			case strings.HasSuffix(key, " LESS THAN") || strings.HasSuffix(key, " <="):
				cleanKey = strings.TrimSuffix(key, " LESS THAN")
				cleanKey = strings.TrimSuffix(cleanKey, " <=")
				query = query.Where(cleanKey+" <= ?", value)
			case strings.HasSuffix(key, " NOT EQUAL") || strings.HasSuffix(key, " !="):
				cleanKey = strings.TrimSuffix(key, " NOT EQUAL")
				cleanKey = strings.TrimSuffix(cleanKey, " !=")
				query = query.Where(cleanKey+" != ?", value)
			case strings.HasSuffix(key, " BETWEEN"):
				cleanKey = strings.TrimSuffix(key, " BETWEEN")
				if values, ok := value.([]interface{}); ok && len(values) == 2 {
					query = query.Where(cleanKey+" BETWEEN ? AND ?", values[0], values[1])
				}
			case strings.HasSuffix(key, "LIKE"):
				cleanKey = strings.TrimSuffix(key, "LIKE")
				query = query.Where("LOWER("+cleanKey+") LIKE ?", "%"+strings.ToLower(value.(string))+"%")
			default:
				query = query.Where(key+" = ?", value)
			}
		}
		if orCondition, ok := conditions["OR"].([]map[string]any); ok {
			orQuery := query.Session(&gorm.Session{NewDB: true})
			for _, cond := range orCondition {
				condition := query.Session(&gorm.Session{NewDB: true})
				condition = applyCondition(condition, cond)
				orQuery = orQuery.Or(condition)
			}
			query = query.Where(orQuery)
		}
	}
	return query
}

func (r *PostgresRepository[T]) Create(ctx context.Context, entity T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *PostgresRepository[T]) Paginate(ctx context.Context, conditions map[string]any, selectFields []string, limit, offset int, sort *string) Pagination[T] {
	var entities []T
	var totalCount int64

	query := r.db.WithContext(ctx).Model(new(T))

	if len(selectFields) > 0 {
		query.Select(selectFields)
	}

	if len(conditions) > 0 {
		query = applyCondition(query, conditions)
	}

	// Get total count before applying pagination
	countQuery := query.Session(&gorm.Session{})
	countQuery.Model(new(T)).Count(&totalCount)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if sort != nil {
		multiSort := strings.Split(*sort, ",")
		for _, sortField := range multiSort {
			sortField = strings.TrimSpace(sortField)
			isDesc := strings.HasPrefix(sortField, "-")
			if isDesc {
				field := sortField[1:]
				query = query.Order(field + " DESC")
				continue
			}
			query = query.Order(sortField)
		}
		query.Order("")
	} else {
		query = query.Order("created_at DESC")
	}
	query.Find(&entities)
	return Pagination[T]{
		Total: totalCount,
		Items: entities,
		Limit: limit,
		Page:  offset/limit + 1,
		Pages: int(math.Ceil(float64(totalCount) / float64(limit))),
	}
}

func (r *PostgresRepository[T]) First(ctx context.Context, id uuid.UUID) (*T, error) {
	var entity T

	condition := map[string]any{
		"id": id,
	}
	query := r.db.WithContext(ctx)

	query = applyCondition(query, condition)

	err := query.First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *PostgresRepository[T]) Update(ctx context.Context, id uuid.UUID, updates map[string]any) error {
	return r.db.WithContext(ctx).Model(new(T)).Where("id = ?", id).Updates(updates).Error
}

func (r *PostgresRepository[T]) DeleteWhere(ctx context.Context, conditions map[string]any) error {
	query := r.db.WithContext(ctx)
	query = applyCondition(query, conditions)
	return query.Delete(new(T)).Error
}

func (r *PostgresRepository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DeleteWhere(ctx, map[string]any{"id": id})
}
