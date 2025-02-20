package db

import (
	"context"
	"watcharis/go-poc-protocal/restful_api/ratelimit/models"
	"watcharis/go-poc-protocal/restful_api/ratelimit/repositories"

	"gorm.io/gorm"
)

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) repositories.ProfilesRepository {
	return &profileRepository{
		db: db,
	}
}

func (r *profileRepository) CreateUserProfile(ctx context.Context, data models.ProfileDB) (models.ProfileDB, error) {
	if err := r.db.WithContext(ctx).Debug().Table("profiles").Create(&data).Error; err != nil {
		return models.ProfileDB{}, err
	}
	return data, nil
}

func (r *profileRepository) GetUserProfile(ctx context.Context, uuid string) (models.ProfileDB, error) {
	var profile models.ProfileDB
	if err := r.db.WithContext(ctx).Debug().Table("profiles").Where("uuid = ?", uuid).First(&profile).Error; err != nil {
		return models.ProfileDB{}, err
	}
	return profile, nil
}
