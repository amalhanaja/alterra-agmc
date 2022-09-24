package repositories

import (
	"alterra-agmc-day-5-6/internal/models"
	"context"
)

type UserRepository interface {
	FindAll(ctx context.Context) ([]*models.User, error)
	FindByID(ctx context.Context, id uint) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	DeleteByID(ctx context.Context, id uint) error
}
