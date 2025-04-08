package repository

import (
	"GoCart/inventory/internal/domain"
	"errors"
	"sync"
)

// define the repository operations
type ProductRepository interface {
	Create(p *domain.Product) error
	GetByID(id string) (*domain.Product, error)
	Update(id string, p *domain.Product) error
	Delete(id string) error
	List() ([]domain.Product, error)
}

// an in-memory implementation
type inMemoryProductRepository struct {
	products map[string]*domain.Product
	mu       sync.RWMutex
}

// return a new in-memory repository
func NewInMemoryProductRepository() ProductRepository {
	return &inMemoryProductRepository{
		products: make(map[string]*domain.Product),
	}
}

func (repo *inMemoryProductRepository) Create(p *domain.Product) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, exists := repo.products[p.ID]; exists {
		return errors.New("product already exists")
	}
	repo.products[p.ID] = p
	return nil
}

func (repo *inMemoryProductRepository) GetByID(id string) (*domain.Product, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	if prod, exists := repo.products[id]; exists {
		return prod, nil
	}
	return nil, errors.New("product not found")
}

func (repo *inMemoryProductRepository) Update(id string, p *domain.Product) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, exists := repo.products[id]; !exists {
		return errors.New("product not found")
	}
	repo.products[id] = p
	return nil
}

func (repo *inMemoryProductRepository) Delete(id string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, exists := repo.products[id]; !exists {
		return errors.New("product not found")
	}
	delete(repo.products, id)
	return nil
}

func (repo *inMemoryProductRepository) List() ([]domain.Product, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	products := []domain.Product{}
	for _, p := range repo.products {
		products = append(products, *p)
	}
	return products, nil
}
