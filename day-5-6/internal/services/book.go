package services

import (
	"alterra-agmc-day-5-6/internal/models"
	"alterra-agmc-day-5-6/internal/repositories"
	"context"
)

type BookService interface {
	FindAll(ctx context.Context) ([]*models.Book, error)
	FindByID(ctx context.Context, id uint) (*models.Book, error)
	Create(ctx context.Context, book *models.Book) (*models.Book, error)
	DeleteByID(ctx context.Context, id uint) error
	Update(ctx context.Context, book *models.Book) (*models.Book, error)
}

type bookServiceImpl struct {
	repo repositories.BookRepository
}

// Create implements BookService
func (s *bookServiceImpl) Create(ctx context.Context, book *models.Book) (*models.Book, error) {
	return s.repo.Create(ctx, book)
}

// DeleteByID implements BookService
func (s *bookServiceImpl) DeleteByID(ctx context.Context, id uint) error {
	return s.repo.DeleteByID(ctx, id)
}

// FindAll implements BookService
func (s *bookServiceImpl) FindAll(ctx context.Context) ([]*models.Book, error) {
	return s.repo.FindAll(ctx)
}

// FindByID implements BookService
func (s *bookServiceImpl) FindByID(ctx context.Context, id uint) (*models.Book, error) {
	return s.repo.FindByID(ctx, id)
}

// Update implements BookService
func (s *bookServiceImpl) Update(ctx context.Context, book *models.Book) (*models.Book, error) {
	return s.repo.Update(ctx, book)
}

func NewBookService(repo repositories.BookRepository) BookService {
	return &bookServiceImpl{repo: repo}
}
