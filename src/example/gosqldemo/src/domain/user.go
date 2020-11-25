package domain

import (
	"database/sql"
	"fmt"
	"time"
)

type JSONTime time.Time

const (
	DateFormat = "2006-01-02"
	TimeFormat = "2006-01-02 15:04:05"
)

func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(TimeFormat))
	return []byte(stamp), nil
}

func (t *JSONTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+TimeFormat+`"`, string(data), time.Local)
	*t = JSONTime(now)
	return
}

type User struct {
	ID           uint
	Name         string
	Email        *string
	Age          uint8
	Birthday     *time.Time
	MemberNumber sql.NullString `db:"member_number"`
	ActivedAt    sql.NullTime   `db:"actived_at"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    sql.NullString `db:"updated_at"`
	DeletedAt    sql.NullString `db:"deleted_at"`
}

func (u User) IsEmpty() bool {
	return User{} == u
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

type Product struct {
	ID        uint
	Code      string
	Price     uint
	UserID    uint           `db:"user_id"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt sql.NullString `db:"updated_at"`
	DeletedAt sql.NullString `db:"deleted_at"`
}

func (p Product) IsEmpty() bool { return p == Product{} || p.ID < 1 }
