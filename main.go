package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"todolist.go/db"
	"todolist.go/service"
)

const port = 8000

func main() {
	// initialize DB connection
	dsn := db.DefaultDSN(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	if err := db.Connect(dsn); err != nil {
		log.Fatal(err)
	}

	// initialize Gin engine
	engine := gin.Default()
	engine.LoadHTMLGlob("views/*.html")

	// prepare session
	store := cookie.NewStore([]byte("my-secret"))
	engine.Use(sessions.Sessions("user-session", store))

	// routing
	engine.Static("/assets", "./assets")
	engine.GET("/", service.Home)
	engine.GET("/list", service.LoginCheck, service.TaskList)

	taskGroup := engine.Group("/task")
	taskGroup.Use(service.LoginCheck)

	// Grouping /task/xxx
	{
		taskGroup.GET("/:id", service.ShowTask) // ":id" is a parameter
		// Create, Update, Delete
		engine.GET("/new", service.NewTaskForm)
		engine.POST("/new", service.NewTask)
		engine.GET("/:id/edit", service.EditTaskForm)
		engine.POST("/:id/edit", service.EditTask)
		engine.GET("/:id/delete", service.DeleteTask)
	}

	// user registration
	engine.GET("/user/new", service.NewUserForm)
	engine.POST("/user/new", service.RegisterUser)
	engine.GET("/user/change_password", service.ChangePasswordForm)
	engine.POST("/user/change_password", service.ChangePassword)
	// login
	engine.GET("/login", service.LoginForm)
	engine.POST("/login", service.Login)

	// start server
	engine.Run(fmt.Sprintf(":%d", port))
}
