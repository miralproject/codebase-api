package repository

import (
	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	DB *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{DB: db}
}

func (r *BaseRepository[T]) Create(entity *T) error {
	return r.DB.Create(entity).Error
}

func (r *BaseRepository[T]) FindByID(id uint) (*T, error) {
	var entity T
	err := r.DB.First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T]) FindAll() ([]T, error) {
	var entities []T
	err := r.DB.Find(&entities).Error
	return entities, err
}

func (r *BaseRepository[T]) Update(entity *T) error {
	return r.DB.Save(entity).Error
}

func (r *BaseRepository[T]) Delete(id uint, entity *T) error {
	return r.DB.Delete(entity, id).Error
}
