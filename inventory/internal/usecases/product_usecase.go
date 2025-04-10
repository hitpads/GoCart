package usecases

import (
	"GoCart/inventory/internal/adapters/repository"
	"GoCart/inventory/internal/domain"
	"errors"
)

// business methods
type ProductUseCase interface {
	CreateProduct(p *domain.Product) error
	GetProductByID(id string) (*domain.Product, error)
	UpdateProduct(id string, p *domain.Product) error
	DeleteProduct(id string) error
	ListProducts() ([]domain.Product, error)
}

type productUseCase struct {
	repo repository.ProductRepository
}

// init product use case
func NewProductUseCase(repo repository.ProductRepository) ProductUseCase {
	return &productUseCase{repo: repo}
}

func (uc *productUseCase) CreateProduct(p *domain.Product) error {
	if p.Name == "" {
		return errors.New("product name is required")
	}
	return uc.repo.Create(p)
}

func (uc *productUseCase) GetProductByID(id string) (*domain.Product, error) {
	return uc.repo.GetByID(id)
}

func (uc *productUseCase) UpdateProduct(id string, p *domain.Product) error {
	return uc.repo.Update(id, p)
}

func (uc *productUseCase) DeleteProduct(id string) error {
	return uc.repo.Delete(id)
}

func (uc *productUseCase) ListProducts() ([]domain.Product, error) {
	return uc.repo.List()
}
