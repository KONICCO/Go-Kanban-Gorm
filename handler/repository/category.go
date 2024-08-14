package repository

import (
	"context"

	"github.com/KONICCO/Go-Kanban-Gorm.git/entity"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error)
	StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error)
	StoreManyCategory(ctx context.Context, categories []entity.Category) error
	GetCategoryByID(ctx context.Context, id int) (entity.Category, error)
	UpdateCategory(ctx context.Context, category *entity.Category) error
	DeleteCategory(ctx context.Context, id int) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error) {
	var cate []entity.Category
	tx := r.db.WithContext(ctx).Model(&entity.Category{}).Where("user_id = ?", id).Scan(&cate)
	if tx.Error != nil {
		return []entity.Category{}, tx.Error
	}
	return cate, nil // TODO: replace this
}

func (r *categoryRepository) StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error) {
	var cate_ entity.Category
	tx := r.db.WithContext(ctx).Create(&category).Scan(&cate_)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return cate_.ID, nil
}

func (r *categoryRepository) StoreManyCategory(ctx context.Context, categories []entity.Category) error {
	var cate []entity.Category
	tx := r.db.WithContext(ctx).Create(&categories).Scan(&cate)
	if tx.Error != nil {
		return tx.Error
	}
	return nil // TODO: replace this
}

func (r *categoryRepository) GetCategoryByID(ctx context.Context, id int) (entity.Category, error) {
	var cate entity.Category
	tx := r.db.WithContext(ctx).Model(&cate).Where("id = ?", id).Scan(&cate)
	if tx.Error != nil {
		return entity.Category{}, tx.Error
	}
	return cate, nil // TODO: replace this
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, category *entity.Category) error {
	var cate entity.Category
	tx := r.db.WithContext(ctx).Model(&entity.Category{}).Updates(
		entity.Category{
			ID:     cate.ID,
			Type:   cate.Type,
			UserID: cate.ID,
		},
	).Scan(&cate)
	if tx.Error != nil {
		return tx.Error
	}
	return nil // TODO: replace this
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id int) error {
	tx := r.db.WithContext(ctx).Delete(&entity.Category{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
