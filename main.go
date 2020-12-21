package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vanilla/gin-crud/api/config"
	"github.com/vanilla/gin-crud/api/controllers"
	"github.com/vanilla/gin-crud/api/middleware"
	"github.com/vanilla/gin-crud/api/repository"
	"github.com/vanilla/gin-crud/api/services"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB                   = config.SetupDatabaseConnection()
	userRepository repository.UserRepository  = repository.NewUserRepository(db)
	bookRepository repository.BookRepository  = repository.NewBookRepository(db)
	jwtService     services.JWTService        = services.NewJWTService()
	userService    services.UserService       = services.NewUserService(userRepository)
	authService    services.AuthService       = services.NewAuthService(userRepository)
	bookService    services.BookService       = services.NewBookService(bookRepository)
	authController controllers.AuthController = controllers.NewAuthController(authService, jwtService)
	userController controllers.UserController = controllers.NewUserController(userService, jwtService)
	bookController controllers.BookController = controllers.NewBookController(bookService, jwtService)
)

func main() {
	defer config.CloseDatabaseConnection(db)

	r := gin.Default()

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	userRoutes := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/getprofile", userController.GetAllProfile)
		userRoutes.GET("/profile", userController.GetProfile)
		userRoutes.PUT("/profile", userController.UpdateProfile)
		userRoutes.DELETE("/profile/:id", userController.DeleteProfile)
	}

	bookRoutes := r.Group("api/books", middleware.AuthorizeJWT(jwtService))
	{
		bookRoutes.GET("/", bookController.GetAll)
		bookRoutes.POST("/", bookController.InsertBook)
		bookRoutes.GET("/:id", bookController.FindByID)
		bookRoutes.PUT("/:id", bookController.UpdateBook)
		bookRoutes.DELETE("/:id", bookController.DeleteBook)
	}

	r.Run()
}
