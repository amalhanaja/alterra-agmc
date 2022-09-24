package datasources

import (
	"alterra-agmc-day-5-6/internal/models"
	"alterra-agmc-day-5-6/internal/repositories"
	"context"
	"time"
)

type BookInMemoryDataSource struct {
	books []*models.Book
}

// Create implements repositories.BookRepository
func (ds *BookInMemoryDataSource) Create(ctx context.Context, book *models.Book) (*models.Book, error) {
	book.CreatedAt = time.Now().UTC()
	book.UpdatedAt = time.Now().UTC()
	book.ID = uint(len(ds.books) + 1)
	ds.books = append(ds.books, book)
	return book, nil
}

// DeleteByID implements repositories.BookRepository
func (ds *BookInMemoryDataSource) DeleteByID(ctx context.Context, id uint) error {
	for i, book := range ds.books {
		if (book.ID) == id {
			ds.books = append(ds.books[:i], ds.books[i+1:]...)
			return nil
		}
	}
	return new(ErrRecordNotFound)
}

// FindAll implements repositories.BookRepository
func (ds *BookInMemoryDataSource) FindAll(ctx context.Context) ([]*models.Book, error) {
	return ds.books, nil
}

// FindByID implements repositories.BookRepository
func (ds *BookInMemoryDataSource) FindByID(ctx context.Context, id uint) (*models.Book, error) {
	for _, book := range ds.books {
		if (book.ID) == id {
			return book, nil
		}
	}
	return nil, new(ErrRecordNotFound)
}

// Update implements repositories.BookRepository
func (ds *BookInMemoryDataSource) Update(ctx context.Context, book *models.Book) (*models.Book, error) {
	for i, b := range ds.books {
		if b.ID == book.ID {
			book.UpdatedAt = time.Now().UTC()
			book.UserID = b.UserID
			ds.books[i] = book
			return book, nil
		}
	}
	return nil, new(ErrRecordNotFound)
}

func NewBookInMemoryDataSource() repositories.BookRepository {
	return &BookInMemoryDataSource{}
}
