package datasources

import (
	gormModels "alterra-agmc-day-5-6/internal/datasources/models"
	"alterra-agmc-day-5-6/internal/models"
	"alterra-agmc-day-5-6/internal/repositories"
	"context"

	"gorm.io/gorm"
)

type UserGormDataSource struct {
	db *gorm.DB
}

// Create implements repositories.UserRepository
func (ds *UserGormDataSource) Create(ctx context.Context, user *models.User) (*models.User, error) {
	userData := gormModels.UserGormModel{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
	userData.CreatedAt = user.CreatedAt
	userData.UpdatedAt = user.UpdatedAt
	if err := ds.db.Create(&userData).Error; err != nil {
		return nil, err
	}
	return &models.User{
		ID:        userData.ID,
		Password:  userData.Password,
		Email:     userData.Email,
		Name:      userData.Name,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
	}, nil
}

// DeleteByID implements repositories.UserRepository
func (ds *UserGormDataSource) DeleteByID(ctx context.Context, id uint) error {
	return ds.db.Delete(&gormModels.UserGormModel{}, id).Error
}

// FindAll implements repositories.UserRepository
func (ds *UserGormDataSource) FindAll(ctx context.Context) ([]*models.User, error) {
	var userData []gormModels.UserGormModel
	var users []*models.User
	if err := ds.db.Find(&userData).Error; err != nil {
		return users, err
	}
	for _, ud := range userData {
		users = append(
			users,
			&models.User{
				ID:        ud.ID,
				Password:  ud.Password,
				Email:     ud.Email,
				Name:      ud.Name,
				CreatedAt: ud.CreatedAt,
				UpdatedAt: ud.UpdatedAt,
			},
		)
	}
	return users, nil
}

// FindByEmail implements repositories.UserRepository
func (ds *UserGormDataSource) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	ud := &gormModels.UserGormModel{}
	if err := ds.db.Where("email = ?", email).First(ud).Error; err != nil {
		return nil, err
	}
	return &models.User{
		ID:        ud.ID,
		Password:  ud.Password,
		Email:     ud.Email,
		Name:      ud.Name,
		CreatedAt: ud.CreatedAt,
		UpdatedAt: ud.UpdatedAt,
	}, nil
}

// FindByID implements repositories.UserRepository
func (ds *UserGormDataSource) FindByID(ctx context.Context, id uint) (*models.User, error) {
	ud := &gormModels.UserGormModel{}
	if err := ds.db.First(ud, id).Error; err != nil {
		return nil, err
	}
	return &models.User{
		ID:        ud.ID,
		Password:  ud.Password,
		Email:     ud.Email,
		Name:      ud.Name,
		CreatedAt: ud.CreatedAt,
		UpdatedAt: ud.UpdatedAt,
	}, nil
}

// Update implements repositories.UserRepository
func (ds *UserGormDataSource) Update(ctx context.Context, user *models.User) (*models.User, error) {
	ud := &gormModels.UserGormModel{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
	ud.ID = user.ID
	ud.UpdatedAt = user.UpdatedAt
	if err := ds.db.Model(ud).Updates(ud).Error; err != nil {
		return nil, err
	}
	return &models.User{
		ID:        ud.ID,
		Password:  ud.Password,
		Email:     ud.Email,
		Name:      ud.Name,
		CreatedAt: ud.CreatedAt,
		UpdatedAt: ud.UpdatedAt,
	}, nil
}

func NewUserGormDataSource(db *gorm.DB) repositories.UserRepository {
	return &UserGormDataSource{db: db}
}
