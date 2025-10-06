package catalog

import (
	"github.com/eya20/hiring_test/models"
)

// CatalogService defines the business logic interface for catalog operations
type CatalogService interface {
	GetProducts() ([]Product, error)
	GetProductsPaginated(offset, limit int) ([]Product, int64, error)
	GetProductsPaginatedWithFilters(offset, limit int, category string, priceLt *float64) ([]Product, int64, error)
	GetProductByCode(code string) (ProductDetails, error)
}

// catalogService implements the business logic for catalog operations
type catalogService struct {
	repo models.ProductsRepositoryInterface
}

// NewCatalogService creates a new catalog service
func NewCatalogService(repo models.ProductsRepositoryInterface) CatalogService {
	return &catalogService{
		repo: repo,
	}
}

// GetProducts retrieves all products and maps them to API format
func (s *catalogService) GetProducts() ([]Product, error) {
	dbProducts, err := s.repo.GetAllProducts()
	if err != nil {
		return nil, err
	}

	// Check if no products found
	if len(dbProducts) == 0 {
		return []Product{}, nil
	}

	// Map database products to API products
	products := make([]Product, len(dbProducts))
	for i, p := range dbProducts {
		products[i] = Product{
			Code:     p.Code,
			Price:    p.Price.InexactFloat64(),
			Category: p.Category.Name,
		}
	}

	return products, nil
}

// GetProductsPaginatedWithFilters retrieves products with pagination and filtering
func (s *catalogService) GetProductsPaginatedWithFilters(offset, limit int, category string, priceLt *float64) ([]Product, int64, error) {
	dbProducts, err := s.repo.GetProductsPaginatedWithFilters(offset, limit, category, priceLt)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.GetProductsCountWithFilters(category, priceLt)
	if err != nil {
		return nil, 0, err
	}

	// Map database products to API products
	products := make([]Product, len(dbProducts))
	for i, p := range dbProducts {
		products[i] = Product{
			Code:     p.Code,
			Price:    p.Price.InexactFloat64(),
			Category: p.Category.Name,
		}
	}

	return products, total, nil
}

// GetProductsPaginated retrieves products with pagination
func (s *catalogService) GetProductsPaginated(offset, limit int) ([]Product, int64, error) {
	dbProducts, err := s.repo.GetProductsPaginated(offset, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.GetProductsCount()
	if err != nil {
		return nil, 0, err
	}

	// Map database products to API products
	products := make([]Product, len(dbProducts))
	for i, p := range dbProducts {
		products[i] = Product{
			Code:     p.Code,
			Price:    p.Price.InexactFloat64(),
			Category: p.Category.Name,
		}
	}

	return products, total, nil
}

// GetProductByCode retrieves a product by its code with variants
func (s *catalogService) GetProductByCode(code string) (ProductDetails, error) {
	var dbProduct models.Product
	if err := s.repo.GetProductByCode(code, &dbProduct); err != nil {
		return ProductDetails{}, err
	}

	// Map variants with price inheritance logic
	variants := make([]Variant, len(dbProduct.Variants))
	for i, v := range dbProduct.Variants {
		price := dbProduct.Price.InexactFloat64() // Default to product price
		if !v.Price.IsZero() {
			price = v.Price.InexactFloat64() // Use variant price if set
		}

		variants[i] = Variant{
			Name:  v.Name,
			SKU:   v.SKU,
			Price: price,
		}
	}

	return ProductDetails{
		Code:     dbProduct.Code,
		Price:    dbProduct.Price.InexactFloat64(),
		Category: dbProduct.Category.Name,
		Variants: variants,
	}, nil
}
