package repositories

import (
	"alterra-agmc-day-7/internal/models"
	"context"
)

type BookRepository interface {
	FindAll(ctx context.Context) ([]*models.Book, error)
	FindByID(ctx context.Context, id uint) (*models.Book, error)
	Create(ctx context.Context, book *models.Book) (*models.Book, error)
	DeleteByID(ctx context.Context, id uint) error
	Update(ctx context.Context, book *models.Book) (*models.Book, error)
}
