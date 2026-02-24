package models

import "time"

type Racer struct {
	ID          string    `db:"id" json:"id"`
	RacerNumber int       `db:"racer_number" json:"racer_number"`
	FullName    string    `db:"full_name" json:"full_name"`
	Birthday    *string   `db:"birthday" json:"birthday,omitempty"`
	Country     *string   `db:"country" json:"country,omitempty"`
	City        *string   `db:"city" json:"city,omitempty"`
	Rating      int       `db:"rating" json:"rating"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// TableName возвращает имя таблицы
func (Racer) TableName() string {
	return "racers"
}

// Validate валидация данных
func (r *Racer) Validate() error {
	if r.FullName == "" {
		return ErrFullNameRequired
	}
	return nil
}

var ErrFullNameRequired = "full name is required"