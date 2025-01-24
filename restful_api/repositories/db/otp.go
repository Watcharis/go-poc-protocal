package db

import (
	"context"
	"watcharis/go-poc-protocal/restful_api/models"
	"watcharis/go-poc-protocal/restful_api/repositories"

	"gorm.io/gorm"
)

type otpRepository struct {
	db *gorm.DB
}

func NewOtpRepository(db *gorm.DB) repositories.OtpRepository {
	return &otpRepository{
		db: db,
	}
}

func (r *otpRepository) CreateOtp(ctx context.Context, data models.OtpDB) (models.OtpDB, error) {
	if err := r.db.WithContext(ctx).Debug().Table("otp_user").Create(&data).Error; err != nil {
		return models.OtpDB{}, err
	}
	return data, nil
}

func (r *otpRepository) GetOtp(ctx context.Context, uuid string, otp string) (models.OtpDB, error) {
	var result models.OtpDB
	if err := r.db.WithContext(ctx).Debug().Table("otp_user").
		Where("uuid=? AND otp=?", uuid, otp).
		Find(&result).Error; err != nil {
		return models.OtpDB{}, err
	}
	return result, nil
}
