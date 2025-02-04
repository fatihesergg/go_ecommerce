package service

import (
	"strconv"
	"time"

	"github.com/fatihesergg/go_ecommerce/internal/model"
	"github.com/fatihesergg/go_ecommerce/internal/storage"
	"github.com/fatihesergg/go_ecommerce/internal/util"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
)

type CategoryService struct {
	Repository storage.CategoryRepository
}

type ProductService struct {
	Repository storage.ProductRepository
}

type UserService struct {
	Repository storage.UserRepository
}

type ReviewService struct {
	Repository storage.ReviewRepository
}

type OrderService struct {
	Repository storage.OrderRepository
}

type PaymentService struct {
	Repository storage.PaymentRepository
}

func NewPaymentService(repostiory storage.PaymentRepository) *PaymentService {
	return &PaymentService{Repository: repostiory}
}

func NewReviewService(repository storage.ReviewRepository) *ReviewService {
	return &ReviewService{Repository: repository}
}

func NewProductService(repository storage.ProductRepository) *ProductService {
	return &ProductService{Repository: repository}
}

func NewCategoryService(repository storage.CategoryRepository) *CategoryService {
	return &CategoryService{Repository: repository}
}

func NewUserService(repository storage.UserRepository) *UserService {
	return &UserService{Repository: repository}
}

func NewOrderService(repository storage.OrderRepository) *OrderService {
	return &OrderService{Repository: repository}
}

// Category Service
func (cs *CategoryService) Get(id string) (model.Category, error) {
	return cs.Repository.Get(id)
}

func (cs *CategoryService) GetAll() ([]model.Category, error) {
	return cs.Repository.GetAll()
}

func (cs *CategoryService) Create(name string) error {
	category := model.Category{Name: name}
	return cs.Repository.Create(category)
}

func (cs *CategoryService) Update(category model.Category) error {
	exist, err := cs.Get(strconv.Itoa(int(category.ID)))
	if err != nil {
		return err
	}
	exist.Name = category.Name
	return cs.Repository.Update(exist)
}

// Product Service
func (ps *ProductService) Get(id string) (model.Product, error) {
	return ps.Repository.Get(id)
}

func (ps *ProductService) GetAll() ([]model.Product, error) {
	return ps.Repository.GetAll()
}

func (ps *ProductService) Create(product model.Product) error {
	return ps.Repository.Create(product)
}

func (ps ProductService) Update(product model.Product) error {
	exist, err := ps.Get(strconv.Itoa(int(product.ID)))
	if err != nil {
		return err
	}

	exist.ImageURL = product.ImageURL
	exist.Name = product.Name
	exist.Price = product.Price
	exist.Stock = product.Stock
	exist.CategoryID = product.CategoryID
	return ps.Repository.Update(exist)
}

func (ps ProductService) Delete(id string) error {
	return ps.Repository.Delete(id)
}

// User Service

func (us *UserService) Get(id string) (model.User, error) {
	return us.Repository.Get(id)
}

func (us *UserService) GetByEmail(email string) (model.User, error) {
	return us.Repository.GetByEmail(email)
}

func (us *UserService) Create(user model.User) error {
	pass := user.Password
	hashedPassword, err := util.EncryptPassword(pass)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return us.Repository.Create(user)
}

func (us *UserService) Update(user model.User) error {
	return us.Repository.Update(user)
}

// Review Service

func (rs *ReviewService) Get(id string) (model.Review, error) {
	return rs.Repository.Get(id)
}

func (rs *ReviewService) Create(review model.Review) error {
	return rs.Repository.Create(review)
}

func (rs *ReviewService) Update(review model.Review) error {
	return rs.Repository.Update(review)
}

func (rs *ReviewService) Delete(id string) error {
	return rs.Repository.Delete(id)
}

func (os *OrderService) Get(id string) (model.Order, error) {
	return os.Repository.Get(id)
}

func (os *OrderService) GetAll() ([]model.Order, error) {
	return os.Repository.GetAll()
}

func (os *OrderService) Create(order model.Order) error {
	return os.Repository.Create(order)
}

func (os *OrderService) Update(order model.Order) error {
	return os.Repository.Update(order)
}

func (ps *PaymentService) SavePayment(model model.Payment) error {
	return ps.Repository.Create(model)
}

func (ps *PaymentService) ProcessPayment(payment model.Payment) error {
	payment.Status = model.PENDING
	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()

	params := stripe.PaymentIntentParams{
		Amount:                  stripe.Int64(int64(payment.Amount)),
		Currency:                stripe.String(string(stripe.CurrencyUSD)),
		PaymentMethod:           stripe.String("pm_card_visa"), // Test
		Confirm:                 stripe.Bool(true),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{AllowRedirects: stripe.String("never"), Enabled: stripe.Bool(true)},
	}
	result, err := paymentintent.New(&params)
	if err != nil {
		payment.Status = model.FAILED
		payment.UpdatedAt = time.Now()
		if err := ps.SavePayment(payment); err != nil {
			return err
		}

		return err
	}
	payment.Status = model.SUCCESS
	payment.UpdatedAt = time.Now()
	payment.TransactionId = result.ID
	if err := ps.SavePayment(payment); err != nil {
		return err
	}
	return nil
}
