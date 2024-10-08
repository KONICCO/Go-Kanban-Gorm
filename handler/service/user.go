package service

import (
	"context"
	"errors"
	"time"

	"github.com/KONICCO/Go-Kanban-Gorm.git/entity"
	"github.com/KONICCO/Go-Kanban-Gorm.git/handler/repository"
)

type UserService interface {
	Login(ctx context.Context, user *entity.User) (id int, err error)
	Register(ctx context.Context, user *entity.User) (entity.User, error)

	Delete(ctx context.Context, id int) error
}

type userService struct {
	userRepository repository.UserRepository
	categoryRepo   repository.CategoryRepository
}

func NewUserService(userRepo repository.UserRepository, categoryRepo repository.CategoryRepository) UserService {
	return &userService{userRepository: userRepo, categoryRepo: categoryRepo}
}

func (s *userService) Login(ctx context.Context, user *entity.User) (id int, err error) {
	//check email and password

	dbUser, err := s.userRepository.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return 0, err
	}

	if dbUser.Email == "" || dbUser.ID == 0 {
		return 0, errors.New("user not found")
	}

	if user.Password != dbUser.Password {
		return 0, errors.New("wrong email or password")
	}

	return dbUser.ID, nil
}

func (s *userService) Register(ctx context.Context, user *entity.User) (entity.User, error) {
	dbUser, err := s.userRepository.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return *user, err
	}

	if dbUser.Email != "" || dbUser.ID != 0 {
		return *user, errors.New("email already exists")
	}

	user.CreatedAt = time.Now()

	newUser, err := s.userRepository.CreateUser(ctx, *user)
	if err != nil {
		return *user, err
	}

	// create 4 category
	// Todo, In Progress, Done, Backlog
	categories := []entity.Category{
		{Type: "Todo", UserID: newUser.ID, CreatedAt: time.Now()},
		{Type: "In Progress", UserID: newUser.ID, CreatedAt: time.Now()},
		{Type: "Done", UserID: newUser.ID, CreatedAt: time.Now()},
		{Type: "Backlog", UserID: newUser.ID, CreatedAt: time.Now()},
	}

	err = s.categoryRepo.StoreManyCategory(ctx, categories)
	if err != nil {
		return *user, err
	}

	return newUser, nil
}

func (s *userService) Delete(ctx context.Context, id int) error {
	return s.userRepository.DeleteUser(ctx, id)
}
