package repository

import (
	"GoCart/order/internal/domain"
	"errors"
	"sync"
)

// define operations for order persistence
type OrderRepository interface {
	Create(o *domain.Order) error
	GetByID(id string) (*domain.Order, error)
	UpdateStatus(id string, status domain.OrderStatus) error
	ListByUser(userID string) ([]domain.Order, error)
}

type inMemoryOrderRepository struct {
	orders map[string]*domain.Order
	mu     sync.RWMutex
}

// return new in-memory order repository
func NewInMemoryOrderRepository() OrderRepository {
	return &inMemoryOrderRepository{
		orders: make(map[string]*domain.Order),
	}
}

func (repo *inMemoryOrderRepository) Create(o *domain.Order) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, exists := repo.orders[o.ID]; exists {
		return errors.New("order already exists")
	}
	repo.orders[o.ID] = o
	return nil
}

func (repo *inMemoryOrderRepository) GetByID(id string) (*domain.Order, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	if order, exists := repo.orders[id]; exists {
		return order, nil
	}
	return nil, errors.New("order not found")
}

func (repo *inMemoryOrderRepository) UpdateStatus(id string, status domain.OrderStatus) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if order, exists := repo.orders[id]; exists {
		order.Status = status
		return nil
	}
	return errors.New("order not found")
}

func (repo *inMemoryOrderRepository) ListByUser(userID string) ([]domain.Order, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	var orders []domain.Order
	for _, o := range repo.orders {
		if o.UserID == userID {
			orders = append(orders, *o)
		}
	}
	return orders, nil
}
