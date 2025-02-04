package storage

import (
	"github.com/fatihesergg/go_ecommerce/internal/model"
	"gorm.io/gorm"
)

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{DB: db}
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{DB: db}
}

// Category Repository
type CategoryRepository struct {
	DB *gorm.DB
}

func (repo *CategoryRepository) Get(id string) (model.Category, error) {
	var result model.Category
	return result, repo.DB.Model(&model.Category{}).First(&result, "id = $1", id).Error
}

func (repo *CategoryRepository) GetAll() ([]model.Category, error) {
	var result []model.Category
	return result, repo.DB.Find(&result).Error
}

func (repo *CategoryRepository) Create(category model.Category) error {
	return repo.DB.Create(&category).Error
}

func (repo *CategoryRepository) Update(category model.Category) error {
	return repo.DB.Save(&category).Error
}

// Product Repository
type ProductRepository struct {
	DB *gorm.DB
}

func (repo *ProductRepository) Get(id string) (model.Product, error) {
	var result model.Product
	return result, repo.DB.Preload("Category").Where("id = $1", id).First(&result).Error
}

func (repo *ProductRepository) GetAll() ([]model.Product, error) {
	var result []model.Product
	return result, repo.DB.Find(&result).Error
}

func (repo *ProductRepository) Create(product model.Product) error {
	return repo.DB.Create(&product).Error
}

func (repo *ProductRepository) Update(product model.Product) error {
	return repo.DB.Save(&product).Error
}

func (repo *ProductRepository) Delete(id string) error {
	product, err := repo.Get(id)
	if err != nil {
		return err
	}
	return repo.DB.Delete(product).Error
}

// User Repository

type UserRepository struct {
	DB *gorm.DB
}

func (repo *UserRepository) Get(id string) (model.User, error) {
	var result model.User
	return result, repo.DB.Model(&model.User{}).First(&result, "id = $1", id).Error
}

func (repo *UserRepository) GetByEmail(email string) (model.User, error) {
	var result model.User
	return result, repo.DB.Model(&model.User{}).First(&result, "email = $1", email).Error
}

func (repo *UserRepository) Create(user model.User) error {
	return repo.DB.Create(&user).Error
}

func (repo *UserRepository) Update(user model.User) error {
	return repo.DB.Save(&user).Error
}

// Review Repository
type ReviewRepository struct {
	DB *gorm.DB
}

func (repo *ReviewRepository) Get(id string) (model.Review, error) {
	var result model.Review
	return result, repo.DB.Model(&model.Review{}).First(&result, "id = $1", id).Error
}

func (repo *ReviewRepository) Create(review model.Review) error {
	return repo.DB.Create(&review).Error
}

func (repo *ReviewRepository) Update(review model.Review) error {
	return repo.DB.Save(&review).Error
}

func (repo *ReviewRepository) Delete(id string) error {
	review, err := repo.Get(id)
	if err != nil {
		return err
	}
	return repo.DB.Delete(&review).Error
}

type OrderRepository struct {
	DB *gorm.DB
}

func (repo *OrderRepository) Get(id string) (model.Order, error) {
	var result model.Order
	return result, repo.DB.Preload("User").Preload("Products").Preload("Products.Product").First(&result, "id = $1", id).Error
}

func (repo *OrderRepository) GetAll() ([]model.Order, error) {
	var result []model.Order
	return result, repo.DB.Model(&model.Order{}).Find(&result).Error
}

func (repo *OrderRepository) Create(order model.Order) error {
	err := repo.DB.Create(&order).Error
	return err
}

func (repo *OrderRepository) Update(order model.Order) error {
	return repo.DB.Model(&model.Order{}).Save(order).Error
}

type PaymentRepository struct {
	DB *gorm.DB
}

func (repo *PaymentRepository) Create(payment model.Payment) error {
	return repo.DB.Model(&model.Payment{}).Create(&payment).Error
}

func (repo *PaymentRepository) Update(payment model.Payment) error {
	return repo.DB.Model(&model.Payment{}).Save(payment).Error
}
