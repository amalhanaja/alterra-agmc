package user

import (
	"alterra-agmc-day-5-6/internal/dto"
	"alterra-agmc-day-5-6/internal/models"
	"context"
)

type UserService interface {
	FindAll(ctx context.Context) ([]*models.User, error)
	FindByID(ctx context.Context, id uint) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, payload *dto.CreateUserPayload) (*models.User, error)
	Update(ctx context.Context, payload *dto.UpdateUserPayload) (*models.User, error)
	DeleteByID(ctx context.Context, id uint) error
}
