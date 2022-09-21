package book

import (
	"alterra-agmc-day-5-6/internal/dto"
	"alterra-agmc-day-5-6/internal/models"
	"context"
)

type BookService interface {
	FindAll(ctx context.Context) ([]*models.Book, error)
	FindByID(ctx context.Context, id uint) (*models.Book, error)
	Create(ctx context.Context, payload *dto.CreateBookPayload) (*models.Book, error)
	DeleteByID(ctx context.Context, id uint) error
	Update(ctx context.Context, id uint, payload *dto.UpdateBookPayload) (*models.Book, error)
}
