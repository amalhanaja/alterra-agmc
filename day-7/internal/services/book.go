package services

import (
	"alterra-agmc-day-7/internal/models"
	"alterra-agmc-day-7/internal/repositories"
	"context"
	"errors"
)

type BookService interface {
	FindAll(ctx context.Context) ([]*models.Book, error)
	FindByID(ctx context.Context, id uint) (*models.Book, error)
	Create(ctx context.Context, book *models.Book) (*models.Book, error)
	DeleteByID(ctx context.Context, id uint, userId uint) error
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
func (s *bookServiceImpl) DeleteByID(ctx context.Context, id uint, userID uint) error {
	book, err := s.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if book.UserID != userID {
		return errors.New("access denied")
	}
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
	b, err := s.FindByID(ctx, book.ID)
	if err != nil {
		return nil, err
	}
	if b.UserID != book.UserID {
		return nil, errors.New("access denied")
	}
	return s.repo.Update(ctx, book)
}

func NewBookService(repo repositories.BookRepository) BookService {
	return &bookServiceImpl{repo: repo}
}
