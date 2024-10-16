package repository

import (
	"codebase-api/internal/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	BaseRepository[domain.User]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepository: *NewBaseRepository[domain.User](db),
	}
}

func (r *UserRepository) GetUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepository) Searching(isActive *bool, search string) ([]domain.User, error) {
	var users []domain.User

	query := r.DB
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("username LIKE ? OR email LIKE ? OR phone LIKE ?", searchTerm, searchTerm, searchTerm)
	}

	err := query.Find(&users).Error
	return users, err
}
