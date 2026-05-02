package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CustomDate struct {
	time.Time
}

const DateFormat = "01-2001"

// ParseJSON парсит "07-2025" из запроса
func (cd *CustomDate) ParseJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}
	t, err := time.Parse(DateFormat, s)
	if err != nil {
		return err
	}
	cd.Time = t
	return nil
}

// FormatJSON отдает "01-2001" в ответе
func (cd CustomDate) FormatJSON() ([]byte, error) {
	if cd.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", cd.Format(DateFormat))), nil
}

// Для работы sqlx с БД
func (cd CustomDate) Value() (driver.Value, error) {
	if cd.IsZero() {
		return nil, nil
	}
	return cd.Time, nil
}

// Для работы sqlx с БД
func (cd *CustomDate) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("type assertion to time.Time failed")
	}
	cd.Time = t
	return nil
}

type Subscription struct {
	ID          int        `db:"id" json:"id"`
	ServiceName string     `db:"service_name" json:"service_name" binding:"required"`
	Price       int        `db:"price" json:"price" binding:"required"`
	UserID      uuid.UUID  `db:"user_id" json:"user_id" binding:"required"`
	StartDate   time.Time  `db:"start_date" json:"start_date" binding:"required"`
	EndDate     *time.Time `db:"end_date" json:"end_date,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
}

type TotalCostResponse struct {
	TotalCost int `json:"total_cost"`
}
