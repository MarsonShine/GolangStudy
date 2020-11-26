package repository

import (
	"fmt"
	"gosqldemo/src/domain"
)

type ProductRepository struct {
}

func (repository *ProductRepository) New() {
	db = singletonInstance()
}

// 创建用户
func (repository *ProductRepository) CreateProduct(p *domain.Product) error {
	_, err := db.Exec("insert into products (user_id,code,price,created_at) values (?, ?, ?, NOW())", p.UserID, p.Code, p.Price)
	if err != nil {
		return fmt.Errorf("用户添加失败: %s\n", err.Error())
	}
	return nil
}

func (repository *ProductRepository) GetProductDetail(productID uint) domain.ProductDto {
	// 字段重命名要全部小写，不能大写，否则 rows.StructScan 无法映射
	sql := `select users.name,users.age,users.email,p.code productname,p.price productprice 
	from users 
	left join products p on p.user_id = users.id 
	where p.id=?`
	productDetail := domain.ProductDto{}
	rows, err := db.Queryx(sql, productID)
	defer rows.Close()
	if err != nil {
		return productDetail
	}
	if rows.Next() {
		err = rows.StructScan(&productDetail)
		if err != nil {
			return productDetail
		}
	}
	return productDetail
}

func (repository *ProductRepository) UpdateProductAndUser(pu *domain.ProductUpdated) bool {
	var product domain.Product
	db.Get(&product, "select * from products where id=?", pu.ID)
	if product.IsEmpty() {
		return false
	}
	var user domain.User
	db.Get(&user, "select id,name,email,age,birthday,member_number,actived_at,created_at,updated_at,deleted_at from users where id = ?", product.UserID)
	if user.IsEmpty() {
		return false
	}
	tx := db.MustBegin()
	tx.MustExec(tx.Rebind("update users set name=?,email=?,age=?,updated_at=NOW() where id=?"), pu.Name, pu.Email, pu.Age, product.UserID)
	tx.MustExec(tx.Rebind("update products set code=?,price=?,updated_at=NOW() where id=?"), pu.ProductName, pu.ProductPrice, pu.ID)
	err := tx.Commit()
	return err == nil
}
