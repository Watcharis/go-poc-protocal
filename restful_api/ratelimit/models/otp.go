package models

import (
	"time"
	"watcharis/go-poc-protocal/pkg/response"
)

const (
	REDIS_OTP           string = "OTP|%s"
	REDIS_RATELIMIT_OTP string = "RATELIMIT_OTP|%s"
)

const (
	RATELIMIT_OTP        int64         = 3
	RATELIMIT_OTP_EXPIRE time.Duration = 1 * time.Minute
	OTP_EXPIRE           time.Duration = 15 * time.Minute
)

type OtpRequest struct {
	UUID string `json:"uuid" validate:"required"`
	Otp  string `json:"otp" validate:"required"`
}

type OtpResponse struct {
	response.CommonResponse
	Error *response.ErrorResponse `json:"error,omitempty"`
}

type OtpDB struct {
	ID        int       `json:"id" gorm:"column:id"`
	UUID      string    `json:"uuid" gorm:"column:uuid"`
	Otp       string    `json:"otp" gorm:"column:otp"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (m *OtpDB) TableName() string {
	return "otp_user"
}
