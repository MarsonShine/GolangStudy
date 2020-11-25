package productservice

import (
	"gosqldemo/src/domain"
	"gosqldemo/src/repository"
)

// 代表用户操作
type ProductService struct {
}

var r = &repository.ProductRepository{}

func NewProductService() *ProductService {
	r.New()
	return &ProductService{}
}

func (ps ProductService) CreateProduct(product *domain.Product) error {
	return r.CreateProduct(product)
}

func (ps ProductService) GetProductDetail(productID uint) domain.ProductDto {
	return r.GetProductDetail(productID)
}

func (ps ProductService) UpdateProductAndUser(productUpdated *domain.ProductUpdated) bool {
	return r.UpdateProductAndUser(productUpdated)
}
