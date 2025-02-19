package main

import (
	"log"
	"os"
	"yoga-management/backend/internal/class"
	"yoga-management/backend/internal/handlers"
	"yoga-management/backend/internal/platform/auth"
	"yoga-management/backend/internal/platform/db"
	"yoga-management/backend/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Load Environment Variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Migrate Models of DB
	db.Database.AutoMigrate(&class.Class{})
	db.Database.AutoMigrate(&user.User{})

	// Router config
	router := gin.Default()

	// Add Access to static archives
	router.Static("/static", "./static")

	// Initialice Handlers
	classHandler := handlers.NewClassHandler(class.NewClassMysqlRepository(*db.Database))
	userHandler := handlers.NewUserHandler(user.NewUserMysqlRepository(*db.Database))

	// Public endpoints
	public := router.Group("/api")
	public.POST("/register", userHandler.RegisterUser)
	public.POST("/login", userHandler.LoginUser)

	// Registered endpoints
	registered := router.Group("/api", auth.AuthenticateMiddleware([]string{"user", "admin"}))
	registered.POST("/logout", handlers.Logout)
	registered.GET("/classes", classHandler.GetClasses)
	registered.GET("/classes/:id", classHandler.GetClassByID)

	// Protected endpoints
	protected := router.Group("/api", auth.AuthenticateMiddleware([]string{"admin"}))
	protected.POST("/classes", classHandler.CreateClass)
	protected.PUT("/classes/:id", classHandler.UpdateClass)
	protected.DELETE("/classes/:id", classHandler.DeleteClass)

	// Publics Views
	router.GET("/login", handlers.Login)
	router.GET("/register", handlers.Register)

	// Protected Views
	private := router.Group("/", auth.AuthenticateMiddleware([]string{"user", "admin"}))
	private.GET("/home", handlers.Home)

	//Listen and Serve APP
	port := os.Getenv("APP_PORT")
	router.Run(port)

}
