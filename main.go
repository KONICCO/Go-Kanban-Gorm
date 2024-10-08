package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/KONICCO/Go-Kanban-Gorm.git/handler/api"
	"github.com/KONICCO/Go-Kanban-Gorm.git/handler/middleware"
	"github.com/KONICCO/Go-Kanban-Gorm.git/handler/repository"
	"github.com/KONICCO/Go-Kanban-Gorm.git/handler/service"
	"github.com/KONICCO/Go-Kanban-Gorm.git/utils"
	"gorm.io/gorm"
)

func main() {
	// fmt.Println("halo")
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// } //pakai env
	// os.Setenv("DATABASE_URL", "postgres://postgres:12345@localhost:5432/kanban") // Hapus jika akan melakukan deploy ke fly.io & lokal
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		mux := http.NewServeMux()

		err := utils.ConnectDB()
		if err != nil {
			panic(err)
		}

		db := utils.GetDBConnection()

		mux = RunServer(db, mux)
		// mux = RunClient(mux, Resources)

		fmt.Println("Server is running on port 8080")
		err = http.ListenAndServe(":8080", mux)
		if err != nil {
			panic(err)
		}
	}()

	wg.Wait()
}

type APIHandler struct {
	UserAPIHandler     api.UserAPI
	TaskAPIHandler     api.TaskAPI
	CategoryAPIHandler api.CategoryAPI
}

func MuxRoute(mux *http.ServeMux, method string, path string, handler http.Handler, opt ...string) {
	if len(opt) > 0 {
		fmt.Printf("[%s]: %s %v \n", method, path, opt)
	} else {
		fmt.Printf("[%s]: %s \n", method, path)
	}

	mux.Handle(path, handler)
}
func RunServer(db *gorm.DB, mux *http.ServeMux) *http.ServeMux {
	// Conn to database
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	// CRUD
	userService := service.NewUserService(userRepo, categoryRepo)
	taskService := service.NewTaskService(taskRepo, categoryRepo)
	categoryService := service.NewCategoryService(categoryRepo, taskRepo)
	//REST API
	userAPIHandler := api.NewUserAPI(userService)
	taskAPIHandler := api.NewTaskAPI(taskService)
	categoryAPIHandler := api.NewCategoryAPI(categoryService)

	apiHandler := APIHandler{
		UserAPIHandler:     userAPIHandler,
		TaskAPIHandler:     taskAPIHandler,
		CategoryAPIHandler: categoryAPIHandler,
	}
	// MuxRoute(mux, "POST", "/api/v1/users/login", middleware.Post(http.HandlerFunc(apiHandler.UserAPIHandler.Login)))
	MuxRoute(mux, "POST", "/api/v1/users/login", middleware.Post(http.HandlerFunc(apiHandler.UserAPIHandler.Login)))
	MuxRoute(mux, "POST", "/api/v1/users/register", middleware.Post(http.HandlerFunc(apiHandler.UserAPIHandler.Register)))
	MuxRoute(mux, "POST", "/api/v1/users/logout", middleware.Post(http.HandlerFunc(apiHandler.UserAPIHandler.Logout)))
	MuxRoute(mux, "DELETE", "/api/v1/users/delete", middleware.Delete(http.HandlerFunc(apiHandler.UserAPIHandler.Delete)), "?user_id=")

	MuxRoute(mux, "GET", "/api/v1/tasks/get", middleware.Get(middleware.Auth(http.HandlerFunc(apiHandler.TaskAPIHandler.GetTask))), "?task_id=")
	MuxRoute(mux, "POST", "/api/v1/tasks/create", middleware.Post(middleware.Auth(http.HandlerFunc(apiHandler.TaskAPIHandler.CreateNewTask))))
	MuxRoute(mux, "PUT", "/api/v1/tasks/update", middleware.Put(middleware.Auth(http.HandlerFunc(apiHandler.TaskAPIHandler.UpdateTask))), "?task_id=")
	MuxRoute(mux, "PUT", "/api/v1/tasks/update/category", middleware.Put(middleware.Auth(http.HandlerFunc(apiHandler.TaskAPIHandler.UpdateTaskCategory))), "?task_id=")
	MuxRoute(mux, "DELETE", "/api/v1/tasks/delete", middleware.Delete(middleware.Auth(http.HandlerFunc(apiHandler.TaskAPIHandler.DeleteTask))), "?task_id=")

	MuxRoute(mux, "GET", "/api/v1/categories/get", middleware.Get(middleware.Auth(http.HandlerFunc(apiHandler.CategoryAPIHandler.GetCategory))))
	MuxRoute(mux, "GET", "/api/v1/categories/dashboard", middleware.Get(middleware.Auth(http.HandlerFunc(apiHandler.CategoryAPIHandler.GetCategoryWithTasks))))
	MuxRoute(mux, "POST", "/api/v1/categories/create", middleware.Post(middleware.Auth(http.HandlerFunc(apiHandler.CategoryAPIHandler.CreateNewCategory))))
	MuxRoute(mux, "DELETE", "/api/v1/categories/delete", middleware.Delete(middleware.Auth(http.HandlerFunc(apiHandler.CategoryAPIHandler.DeleteCategory))), "?category_id=")

	return mux
}
