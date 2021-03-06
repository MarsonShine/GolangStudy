package models

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
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
	MemberNumber sql.NullString
	ActivedAt    sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt // 关于软删除：https://gorm.io/docs/delete.html#Soft-Delete
}

func (u User) IsEmpty() bool { return u == User{} || u.ID < 1 }

type Model struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
