package models

import (
	"gorm.io/gorm"
)

// CategoriesRepositoryInterface defines the contract for category repository operations
type CategoriesRepositoryInterface interface {
	GetAllCategories() ([]Category, error)
	GetCategoryByCode(code string) (Category, error)
	CreateCategory(category *Category) error
}

type CategoriesRepository struct {
	db *gorm.DB
}

func NewCategoriesRepository(db *gorm.DB) *CategoriesRepository {
	return &CategoriesRepository{
		db: db,
	}
}

func (r *CategoriesRepository) GetAllCategories() ([]Category, error) {
	var categories []Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoriesRepository) GetCategoryByCode(code string) (Category, error) {
	var category Category
	if err := r.db.Where("code = ?", code).First(&category).Error; err != nil {
		return Category{}, err
	}
	return category, nil
}

// CreateCategory creates a new category in the database
func (r *CategoriesRepository) CreateCategory(category *Category) error {
	if err := r.db.Create(category).Error; err != nil {
		return err
	}
	return nil
}
