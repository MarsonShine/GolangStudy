package appservice

import (
	"fmt"
	"gormdemo/src/models"

	"gorm.io/gorm"
)

// 代表用户操作
type ProductService struct {
	db *gorm.DB
}

type ProductDto struct {
	Name         string `json:"name"`
	Age          uint   `json:"age"`
	Email        string `json:"email"`
	ProductName  string `json:"productName"`
	ProductPrice uint   `json:"productPrice"`
}

type ProductUpdated struct {
	ID           uint
	Name         string
	Age          uint
	Email        string
	ProductName  string
	ProductPrice uint
}

func NewProductService(gormDB *gorm.DB) ProductService {
	return ProductService{gormDB}
}

func (ps ProductService) CreateProduct(product *models.Product) error {
	r := ps.db.Create(product)
	if r.RowsAffected <= 0 && r.Error != nil {
		return r.Error
	}
	return nil
}

func (ps ProductService) GetProductDetail(id uint) ProductDto {
	var productDetail ProductDto
	// 联表查询
	ps.db.Table("users").Select("users.name,users.age,users.email,p.code ProductName,p.price ProductPrice").Where("p.id=?", id).Joins("left join products p on p.user_id = users.id").Scan(&productDetail)
	return productDetail
}

func (ps ProductService) UpdateProductAndUser(pu *ProductUpdated) (bool, error) {
	var p models.Product
	ps.db.First(&p, pu.ID)
	if p.IsEmpty() {
		return false, fmt.Errorf("用户产品不存在")
	}
	var u models.User
	ps.db.First(&u, p.UserID)
	// 事务：https://gorm.io/zh_CN/docs/transactions.html
	err := ps.db.Transaction(func(tx *gorm.DB) error {
		u.Name = pu.Name
		u.Email = &pu.Email
		u.Age = uint8(pu.Age)
		tx.Save(u)
		p.Code = pu.ProductName
		p.Price = pu.ProductPrice
		tx.Save(p)
		return nil
	})

	return err == nil, nil
}
