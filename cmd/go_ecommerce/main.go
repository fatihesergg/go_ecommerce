package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/fatihesergg/go_ecommerce/internal/handler"
	"github.com/fatihesergg/go_ecommerce/internal/middleware"
	"github.com/fatihesergg/go_ecommerce/internal/model"
	"github.com/fatihesergg/go_ecommerce/internal/service"
	"github.com/fatihesergg/go_ecommerce/internal/storage"
	"github.com/fatihesergg/go_ecommerce/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title						go_ecommerce API
// @version					1.0
// @description				Basic e-commerce api writtin in go.
// @termsOfService				http://swagger.io/terms/
//
// @securityDefinitions.apikey	Bearer
// @in							header
// @name						Authorization
func main() {
	// Database
	dsn := "postgresql://localhost/go_ecommerce?user=fatih&password=test"
	address := ":3000"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	sugar := logger.Sugar()
	defer logger.Sync()
	middleware.LOGGER = sugar

	// JWT
	util.JWTSECRET = os.Getenv("JWT_SECRET")
	if util.JWTSECRET == "" {
		panic(errors.New("JWT_SECRET can't be empty"))
	}

	db.AutoMigrate(&model.Product{})
	db.AutoMigrate(&model.Category{})
	db.AutoMigrate(&model.Order{})
	db.AutoMigrate(&model.OrderItem{})
	db.AutoMigrate(&model.Payment{})
	db.AutoMigrate(&model.Review{})
	db.AutoMigrate(&model.User{})

	validate := validator.New(validator.WithRequiredStructEnabled())

	// Repositories
	categoryRepo := storage.NewCategoryRepository(db)
	productRepo := storage.NewProductRepository(db)
	userRepo := storage.NewUserRepository(db)
	reviewRepo := storage.NewReviewRepository(db)
	orderRepo := storage.NewOrderRepository(db)
	paymentRepo := storage.NewPaymentRepository(db)

	// Services
	categoryService := service.NewCategoryService(*categoryRepo)
	productService := service.NewProductService(*productRepo)
	userService := service.NewUserService(*userRepo)
	reviewService := service.NewReviewService(*reviewRepo)
	orderService := service.NewOrderService(*orderRepo)
	paymentService := service.NewPaymentService(*paymentRepo)

	// Handlers
	categoryHandler := handler.NewCategoryHandler(*categoryService, validate)
	producthandler := handler.NewProductHandler(*productService, *categoryService, validate)
	authHandler := handler.NewAuthHandler(*userService, validate)
	reviewHandler := handler.NewReviewHandler(*reviewService, *userService, *productService, validate)
	orderHandler := handler.NewOrderHandler(*orderService, *productService, *categoryService, *userService, validate)
	paymentHandler := handler.NewPaymentHandler(*paymentService, *orderService, validate)

	fs := http.FileServer(http.Dir("../../docs"))
	apiRouter := http.NewServeMux()

	apiRouter.Handle("/docs/", http.StripPrefix("/docs/", fs))
	// Category
	apiRouter.HandleFunc("GET /category", categoryHandler.GetAll)
	apiRouter.HandleFunc("GET /category/{id}", categoryHandler.Get)
	apiRouter.HandleFunc("POST /category", middleware.RequireLogin("admin", categoryHandler.Create))
	apiRouter.HandleFunc("PUT /category", middleware.RequireLogin("admin", categoryHandler.Update))

	// Product
	apiRouter.HandleFunc("GET /product", producthandler.GetAll)
	apiRouter.HandleFunc("GET /product/{id}", producthandler.Get)
	apiRouter.HandleFunc("POST /product", middleware.RequireLogin("admin", producthandler.Create))
	apiRouter.HandleFunc("PUT /product", middleware.RequireLogin("admin", producthandler.Update))
	apiRouter.HandleFunc("DELETE /product/{id}", middleware.RequireLogin("admin", producthandler.Delete))

	// Review
	apiRouter.HandleFunc("GET /review/{id}", reviewHandler.Get)
	apiRouter.HandleFunc("POST /review", middleware.RequireLogin("user", reviewHandler.Create))
	apiRouter.HandleFunc("PUT /review", middleware.RequireLogin("user", reviewHandler.Update))
	apiRouter.HandleFunc("DELETE /review/{id}", middleware.RequireLogin("user", reviewHandler.Delete))

	// Order
	apiRouter.HandleFunc("GET /order/{id}", middleware.RequireLogin("user", orderHandler.Get))
	apiRouter.HandleFunc("POST /order", middleware.RequireLogin("user", orderHandler.Create))

	// Payment
	apiRouter.HandleFunc("POST /payment/{id}", middleware.RequireLogin("user", paymentHandler.Create))

	// Auth
	apiRouter.HandleFunc("POST /login", authHandler.Login)
	apiRouter.HandleFunc("POST /register", authHandler.Register)

	// Swagger
	apiRouter.HandleFunc("/swagger/", httpSwagger.Handler(httpSwagger.URL("http://localhost:3000/docs/swagger.json")))

	if _, err := userService.GetByEmail("admin@example.com"); err != nil {
		userService.Create(model.User{
			Name:     "admin",
			LastName: "admin",
			UserName: "admin",
			Email:    "admin@example.com",
			Role:     "admin",
			Password: "1234",
		})
	}
	// Test
	http.ListenAndServe(address, middleware.LoggerMiddleware(apiRouter))
}
