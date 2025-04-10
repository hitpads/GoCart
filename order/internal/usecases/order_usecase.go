package usecases

import (
	"GoCart/order/internal/adapters/repository"
	"GoCart/order/internal/domain"
	"errors"
	"time"
)

// business methods
type OrderUseCase interface {
	CreateOrder(o *domain.Order) error
	GetOrderByID(id string) (*domain.Order, error)
	UpdateOrderStatus(id string, status domain.OrderStatus) error
	ListOrdersByUser(userID string) ([]domain.Order, error)
}

type orderUseCase struct {
	repo repository.OrderRepository
}

// create new OrderUseCase
func NewOrderUseCase(repo repository.OrderRepository) OrderUseCase {
	return &orderUseCase{repo: repo}
}

func (uc *orderUseCase) CreateOrder(o *domain.Order) error {
	if o.UserID == "" || len(o.Items) == 0 {
		return errors.New("invalid order: missing user_id or items")
	}
	o.Status = domain.StatusPending
	o.CreatedAt = time.Now()
	return uc.repo.Create(o)
}

func (uc *orderUseCase) GetOrderByID(id string) (*domain.Order, error) {
	return uc.repo.GetByID(id)
}

func (uc *orderUseCase) UpdateOrderStatus(id string, status domain.OrderStatus) error {
	return uc.repo.UpdateStatus(id, status)
}

func (uc *orderUseCase) ListOrdersByUser(userID string) ([]domain.Order, error) {
	return uc.repo.ListByUser(userID)
}
