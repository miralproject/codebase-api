package usecase

import (
	"codebase-api/internal/domain"
	"codebase-api/internal/repository"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepo *repository.UserRepository
}

func NewUserUseCase(userRepo *repository.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

func (u *UserUseCase) Register(user *domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return u.userRepo.Create(user)
}

func (u *UserUseCase) Login(username, password string) (*domain.User, error) {
	user, err := u.userRepo.GetUserByUsername(username)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (u *UserUseCase) FinAll() ([]domain.User, error) {
	users, err := u.userRepo.FindAll()
	if err != nil {
		return nil, errors.New("users not found")
	}
	return users, nil
}

func (u *UserUseCase) Searching(isActive *bool, search string) ([]domain.User, error) {
	users, err := u.userRepo.Searching(isActive, search)
	if err != nil {
		return nil, errors.New("users not found")
	}
	return users, nil
}

func (u *UserUseCase) FindById(id uint) (*domain.User, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (u *UserUseCase) Delete(id uint) error {
	var user domain.User
	err := u.userRepo.Delete(id, &user)
	if err != nil {
		return errors.New("failed to delete user")
	}
	return nil
}
