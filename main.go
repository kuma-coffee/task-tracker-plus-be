package main

import (
	"a21hc3NpZ25tZW50/db"
	"a21hc3NpZ25tZW50/handler/api"
	"a21hc3NpZ25tZW50/middleware"
	"a21hc3NpZ25tZW50/model"
	repo "a21hc3NpZ25tZW50/repository"
	"a21hc3NpZ25tZW50/service"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type APIHandler struct {
	UserAPIHandler     api.UserAPI
	CategoryAPIHandler api.CategoryAPI
	TaskAPIHandler     api.TaskAPI
}

func main() {
	gin.SetMode(gin.ReleaseMode) //release

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		router := gin.New()
		db := db.NewDB()
		router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("[%s] \"%s %s %s\"\n",
				param.TimeStamp.Format(time.RFC822),
				param.Method,
				param.Path,
				param.ErrorMessage,
			)
		}))
		router.Use(gin.Recovery())

		// dbCredential := model.Credential{
		// 	Host:         "localhost",
		// 	Username:     "postgres",
		// 	Password:     "postgres",
		// 	DatabaseName: "kampusmerdeka",
		// 	Port:         5432,
		// 	Schema:       "public",
		// }

		// conn, err := db.Connect(&dbCredential)
		// if err != nil {
		// 	panic(err)
		// }

		// os.Setenv("DATABASE_URL", "postgres://postgres:hiwOus48NkMMSSE@localhost:15432/postgres") // <- Gunakan ini untuk connect database di localhost
		os.Getenv("DATABASE_URL")
		clientUrl := os.Getenv("CLIENT_URL")

		conn, err := db.Connect()
		if err != nil {
			panic(err)
		}

		conn.AutoMigrate(&model.User{}, &model.Session{}, &model.Category{}, &model.Task{})

		config := cors.DefaultConfig()
		config.AllowOrigins = []string{clientUrl}
		config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Cookie"}
		config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
		config.AllowCredentials = true

		router.Use(cors.New(config))

		router = RunServer(conn, router)

		fmt.Println("Server is running on port 8080")
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		if err := router.Run(":" + port); err != nil {
			log.Panicf("error: %s", err)
		}
	}()

	wg.Wait()
}

func RunServer(db *gorm.DB, gin *gin.Engine) *gin.Engine {
	userRepo := repo.NewUserRepo(db)
	sessionRepo := repo.NewSessionsRepo(db)
	categoryRepo := repo.NewCategoryRepo(db)
	taskRepo := repo.NewTaskRepo(db)

	userService := service.NewUserService(userRepo, sessionRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	taskService := service.NewTaskService(taskRepo)

	userAPIHandler := api.NewUserAPI(userService)
	categoryAPIHandler := api.NewCategoryAPI(categoryService)
	taskAPIHandler := api.NewTaskAPI(taskService)

	apiHandler := APIHandler{
		UserAPIHandler:     userAPIHandler,
		CategoryAPIHandler: categoryAPIHandler,
		TaskAPIHandler:     taskAPIHandler,
	}

	version := gin.Group("/api/v1")
	{
		user := version.Group("/user")
		{
			user.POST("/login", apiHandler.UserAPIHandler.Login)
			user.POST("/register", apiHandler.UserAPIHandler.Register)

			user.Use(middleware.Auth())
			user.GET("/tasks", apiHandler.UserAPIHandler.GetUserTaskCategory)
			user.GET("/logout", apiHandler.UserAPIHandler.Logout)
		}

		task := version.Group("/task")
		{
			task.Use(middleware.Auth())
			task.POST("/add", apiHandler.TaskAPIHandler.AddTask)
			task.GET("/get/:id", apiHandler.TaskAPIHandler.GetTaskByID)
			task.PUT("/update/:id", apiHandler.TaskAPIHandler.UpdateTask)
			task.DELETE("/delete/:id", apiHandler.TaskAPIHandler.DeleteTask)
			task.GET("/list", apiHandler.TaskAPIHandler.GetTaskList)
			task.GET("/category/:id", apiHandler.TaskAPIHandler.GetTaskListByCategory)
		}

		category := version.Group("/category")
		{
			category.Use(middleware.Auth())
			category.POST("/add", apiHandler.CategoryAPIHandler.AddCategory)
			category.GET("/get/:id", apiHandler.CategoryAPIHandler.GetCategoryByID)
			category.PUT("/update/:id", apiHandler.CategoryAPIHandler.UpdateCategory)
			category.DELETE("/delete/:id", apiHandler.CategoryAPIHandler.DeleteCategory)
			category.GET("/list", apiHandler.CategoryAPIHandler.GetCategoryList)
		}
	}

	return gin
}
