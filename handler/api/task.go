package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/KONICCO/Go-Kanban-Gorm.git/entity"
	"github.com/KONICCO/Go-Kanban-Gorm.git/handler/service"
)

type TaskAPI interface {
	GetTask(w http.ResponseWriter, r *http.Request)
	CreateNewTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
	UpdateTaskCategory(w http.ResponseWriter, r *http.Request)
}

type taskAPI struct {
	taskService service.TaskService
}

func NewTaskAPI(taskService service.TaskService) *taskAPI {
	return &taskAPI{taskService}
}

func (t *taskAPI) GetTask(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here
	userId := r.Context().Value("id")
	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}
	tasID := r.URL.Query().Get("task_id")
	if tasID == "" {
		GETtas, err := t.taskService.GetTasks(r.Context(), idLogin)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GETtas)
		return
	} else {
		xx, _ := strconv.Atoi(tasID)
		GETtas, err := t.taskService.GetTaskByID(r.Context(), xx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GETtas)
		return
	}

}

func (t *taskAPI) CreateNewTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}
	if task.Title == "" || task.Description == "" || strconv.Itoa(task.CategoryID) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}
	userId := r.Context().Value("id")
	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}
	TAS := entity.Task{

		Title:       task.Title,
		Description: task.Description,
		CategoryID:  task.CategoryID,
		UserID:      idLogin,
	}
	GETtas, err := t.taskService.StoreTask(r.Context(), &TAS)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": idLogin,
		"task_id": GETtas.ID,
		"message": "success create new task",
	})
	// TODO: answer here
}

func (t *taskAPI) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here
	userId := r.Context().Value("id")
	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}
	taskID := r.URL.Query().Get("task_id")
	xx, _ := strconv.Atoi(taskID)
	DELtas := t.taskService.DeleteTask(r.Context(), xx)
	if DELtas != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": idLogin,
		"task_id": xx,
		"message": "success delete task",
	})
}

func (t *taskAPI) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}
	userId := r.Context().Value("id")
	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}
	TAS := &entity.Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		CategoryID:  task.CategoryID,
	}
	UPtas, err := t.taskService.UpdateTask(r.Context(), TAS)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": UPtas.UserID,
		"task_id": idLogin,
		"message": "success update task",
	})
	// TODO: answer here
}

func (t *taskAPI) UpdateTaskCategory(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskCategoryRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	var updateTask = entity.Task{
		ID:         task.ID,
		CategoryID: task.CategoryID,
		UserID:     int(idLogin),
	}

	_, err = t.taskService.UpdateTask(r.Context(), &updateTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userId,
		"task_id": task.ID,
		"message": "success update task category",
	})
}
