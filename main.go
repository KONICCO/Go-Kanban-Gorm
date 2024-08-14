package main

import (
	"fmt"
	"net/http"
	"os"
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
	os.Setenv("DATABASE_URL", "postgres://postgres:12345@localhost:5432/kanban") // Hapus jika akan melakukan deploy ke fly.io
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

	MuxRoute(mux, "GET", "/api/v1/tasks/get", middleware.Get(middleware.Auth(http.HandlerFunc(apiHandler.TaskAPIHandler.GetTask))), "?task_id=")

	MuxRoute(mux, "GET", "/api/v1/categories/get", middleware.Get(middleware.Auth(http.HandlerFunc(apiHandler.CategoryAPIHandler.GetCategory))))

	return mux
}
