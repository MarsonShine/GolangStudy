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
