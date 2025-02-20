package models

import (
	"time"
	"watcharis/go-poc-protocal/pkg/response"
)

const (
	REDIS_USER_PROFILE = "USER_PROFILE|%s"
)

type ProifleRequest struct {
	FirstName string `json:"firstname" validate:"required"`
	LastName  string `json:"lastname" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
}

type ProifleResponse struct {
	response.CommonResponse
	Data  *ProfileDB              `json:"data,omitempty"`
	Error *response.ErrorResponse `json:"error,omitempty"`
}

type ProfileDB struct {
	ID        int       `json:"id" gorm:"column:id" redis:"id" validate:"required"`
	UUID      string    `json:"uuid" gorm:"column:uuid" redis:"uuid" validate:"required"`
	FirstName string    `json:"firstname" gorm:"column:firstname" redis:"firstname" validate:"required"`
	LastName  string    `json:"lastname" gorm:"column:lastname" redis:"lastname" validate:"required"`
	Email     string    `json:"email" gorm:"column:email" redis:"email" validate:"required"`
	Phone     string    `json:"phone" gorm:"column:phone" redis:"phone" validate:"required"`
	CreatedAt time.Time `json:"-" gorm:"column:created_at" redis:"-" validate:"required"`
	UpdatedAt time.Time `json:"-" gorm:"column:created_at" redis:"-" validate:"required"`
}

func (m *ProfileDB) TableName() string {
	return "profiles"
}
