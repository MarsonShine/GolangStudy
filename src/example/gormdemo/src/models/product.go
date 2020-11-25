package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code   string
	Price  uint
	UserID uint
}

func (p Product) IsEmpty() bool { return p == Product{} || p.ID < 1 }
