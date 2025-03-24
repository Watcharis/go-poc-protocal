package db

import (
	"context"
	"watcharis/go-poc-protocal/restful_api/ratelimit/models"

	"gorm.io/gorm"
)

// mockgen -source=db/profiles.go -destination=db/mocks/profiles_mock.go -package=mocks
type ProfilesRepository interface {
	CreateUserProfile(ctx context.Context, data models.ProfileDB) (models.ProfileDB, error)
	GetUserProfile(ctx context.Context, uuid string) (models.ProfileDB, error)
}

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) ProfilesRepository {
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
