package repository

import (
	"context"

	"github.com/KONICCO/Go-Kanban-Gorm.git/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
	UpdateUser(ctx context.Context, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	var user entity.User
	tx := r.db.WithContext(ctx).Model(&user).Where("id = ?", id).Scan(&user)
	if tx.Error != nil {
		return entity.User{}, tx.Error
	}
	return user, nil // TODO: replace this
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	tx := r.db.WithContext(ctx).Model(&user).Where("email = ?", email).Scan(&user)
	if tx.Error != nil {
		return entity.User{}, tx.Error
	}
	return user, nil // TODO: replace this
	// return entity.User{}, nil // TODO: replace this
}

func (r *userRepository) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	var users entity.User

	tx := r.db.WithContext(ctx).Create(&user).Scan(&users)
	if tx.Error != nil {
		return entity.User{}, tx.Error
	}
	return users, nil // TODO: replace this
}

func (r *userRepository) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	var users entity.User
	tx := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", user.ID).Updates(entity.User{
		ID:       user.ID,
		Fullname: user.Fullname,
		Email:    user.Email,
		Password: user.Password,
	}).Scan(&user)
	if tx.Error != nil {
		return entity.User{}, tx.Error
	}
	return users, nil // TODO: replace this
}

func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	tx := r.db.WithContext(ctx).Delete(&entity.User{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	return nil // TODO: replace this
}
