package repository

import (
	"context"

	"github.com/KONICCO/Go-Kanban-Gorm.git/entity"
	"gorm.io/gorm"
)

type TaskRepository interface {
	GetTasks(ctx context.Context, id int) ([]entity.Task, error)
	StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error)
	GetTaskByID(ctx context.Context, id int) (entity.Task, error)
	GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error)
	UpdateTask(ctx context.Context, task *entity.Task) error
	DeleteTask(ctx context.Context, id int) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db}
}

func (r *taskRepository) GetTasks(ctx context.Context, id int) ([]entity.Task, error) {
	var task []entity.Task
	tx := r.db.WithContext(ctx).Model(&task).Where("user_id = ?", id).Scan(&task)
	if tx.Error != nil {
		return []entity.Task{}, tx.Error
	}
	return task, nil // TODO: replace this
}

func (r *taskRepository) StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error) {
	var task_ entity.Task
	tx := r.db.WithContext(ctx).Create(&task).Scan(&task_)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return task_.ID, nil
}

func (r *taskRepository) GetTaskByID(ctx context.Context, id int) (entity.Task, error) {
	var task entity.Task
	tx := r.db.WithContext(ctx).Model(&task).Where("id = ?", id).Scan(&task)
	if tx.Error != nil {
		return entity.Task{}, tx.Error
	}
	return task, nil // TODO: replace this
}

func (r *taskRepository) GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error) {
	task_ := []entity.Task{}
	tx := r.db.WithContext(ctx).Model(&task_).Where("category_id = ?", catId).Scan(&task_)
	if tx.Error != nil {
		return []entity.Task{}, tx.Error
	}
	return task_, nil // TODO: replace this
}

func (r *taskRepository) UpdateTask(ctx context.Context, task *entity.Task) error {
	var tas_ entity.Task
	tx := r.db.WithContext(ctx).Model(&entity.Task{}).Where("id = ?", task.ID).Updates(
		entity.Task{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			CategoryID:  task.CategoryID,
			UserID:      task.UserID,
		},
	).Scan(&tas_)
	if tx.Error != nil {
		return tx.Error
	}
	return nil // TODO: replace this
}

func (r *taskRepository) DeleteTask(ctx context.Context, id int) error {
	tx := r.db.WithContext(ctx).Delete(&entity.Task{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	return nil // TODO: replace this
}
