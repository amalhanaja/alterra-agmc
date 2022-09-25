package datasources_test

import (
	"alterra-agmc-day-7/internal/datasources"
	"alterra-agmc-day-7/internal/models"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateBook(t *testing.T) {
	testCases := []struct {
		name          string
		book          *models.Book
		expectedCount int
	}{
		{
			name: "Test Create book should add book to memory",
			book: &models.Book{
				Title:  "title",
				Writer: "writer",
				UserID: 12,
				Isbn:   "isbn",
			},
			expectedCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			ds := datasources.NewBookInMemoryDataSource()

			// Assert
			_, err := ds.Create(context.TODO(), tc.book)
			all, err := ds.FindAll(context.TODO())

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCount, len(all))
		})
	}
}

func TestFindAllBook(t *testing.T) {
	testCases := []struct {
		name     string
		books    []*models.Book
		expected []*models.Book
	}{
		{
			name:     "Test FindAll with empty books",
			books:    []*models.Book{},
			expected: []*models.Book{},
		},
		{
			name: "Test FindAll with empty books",
			books: []*models.Book{
				{
					Title:  "title",
					Writer: "writer",
					UserID: 12,
					Isbn:   "isbn",
				},
			},
			expected: []*models.Book{
				{
					Title:  "title",
					Writer: "writer",
					UserID: 12,
					Isbn:   "isbn",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ds := datasources.NewBookInMemoryDataSource()
			for _, b := range tc.books {
				_, _ = ds.Create(context.TODO(), b)
			}

			// Act
			result, _ := ds.FindAll(context.TODO())
			assert.Equal(t, len(tc.expected), len(result))
			for i, r := range result {
				expectedBook := tc.books[i]
				assert.Equal(t, expectedBook.Title, r.Title)
				assert.Equal(t, expectedBook.Writer, r.Writer)
				assert.Equal(t, expectedBook.Isbn, r.Isbn)
				assert.Equal(t, expectedBook.UserID, r.UserID)
			}
		})
	}
}

func TestDeleteByID(t *testing.T) {
	testCases := []struct {
		name          string
		books         []*models.Book
		id            uint
		expected      error
		expectedCount int
	}{
		{
			name:          "Test delete when id not found should return err record not found",
			books:         []*models.Book{},
			id:            9,
			expected:      datasources.ErrRecordNotFound{},
			expectedCount: 0,
		},
		{
			name: "Test delete when id found should return nil",
			books: []*models.Book{
				{
					Title:  "title",
					Writer: "writer",
					UserID: 12,
					Isbn:   "isbn",
				},
				{
					Title:  "title2",
					Writer: "writer2",
					UserID: 12,
					Isbn:   "isbn2",
				},
			},
			id:            1,
			expected:      nil,
			expectedCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ds := datasources.NewBookInMemoryDataSource()
			for _, b := range tc.books {
				_, _ = ds.Create(context.TODO(), b)
			}

			// Act
			result := ds.DeleteByID(context.TODO(), tc.id)
			all, _ := ds.FindAll(context.TODO())

			// Assert
			assert.Equal(t, tc.expected.Error(), result.Error())
			assert.Equal(t, tc.expectedCount, len(all))
		})
	}
}

func TestFindBookByID(t *testing.T) {
	testCases := []struct {
		name         string
		books        []*models.Book
		id           uint
		expectedErr  error
		expectedBook *models.Book
	}{
		{
			name:         "Test findby id when id not found should return err record not found",
			books:        []*models.Book{},
			id:           9,
			expectedErr:  datasources.ErrRecordNotFound{},
			expectedBook: nil,
		},
		{
			name: "Test find bt id when id found should return expected book",
			books: []*models.Book{
				{
					Title:  "title",
					Writer: "writer",
					UserID: 12,
					Isbn:   "isbn",
				},
			},
			id:          1,
			expectedErr: nil,
			expectedBook: &models.Book{
				Title:  "title",
				Writer: "writer",
				UserID: 12,
				Isbn:   "isbn",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ds := datasources.NewBookInMemoryDataSource()
			for _, b := range tc.books {
				_, _ = ds.Create(context.TODO(), b)
			}

			// Act
			result, err := ds.FindByID(context.TODO(), tc.id)

			// Assert
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			}
			if tc.expectedBook != nil {
				assert.Equal(t, tc.expectedBook.Title, result.Title)
				assert.Equal(t, tc.expectedBook.Writer, result.Writer)
				assert.Equal(t, tc.expectedBook.Isbn, result.Isbn)
				assert.Equal(t, tc.expectedBook.UserID, result.UserID)
			}
		})
	}
}

func TestUpdateeBook(t *testing.T) {
	testCases := []struct {
		name         string
		books        []*models.Book
		updateBook   *models.Book
		expectedErr  error
		expectedBook *models.Book
	}{
		{
			name:  "Test update when id not found should return err record not found",
			books: []*models.Book{},
			updateBook: &models.Book{
				Title:  "title",
				Writer: "writer",
				UserID: 12,
				Isbn:   "isbn",
			},
			expectedErr:  datasources.ErrRecordNotFound{},
			expectedBook: nil,
		},
		{
			name: "Test delete when id found should return nil",
			books: []*models.Book{
				{
					Title:  "title",
					Writer: "writer",
					UserID: 12,
					Isbn:   "isbn",
				},
				{
					Title:  "title2",
					Writer: "writer2",
					UserID: 15,
					Isbn:   "isbn2",
				},
			},
			updateBook: &models.Book{
				Title:  "title_new",
				Writer: "writer_new",
				ID:     2,
				Isbn:   "isbn_new",
			},
			expectedBook: &models.Book{
				Title:  "title_new",
				Writer: "writer_new",
				UserID: 15,
				Isbn:   "isbn_new",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ds := datasources.NewBookInMemoryDataSource()
			for _, b := range tc.books {
				_, _ = ds.Create(context.TODO(), b)
			}

			// Act
			result, err := ds.Update(context.TODO(), tc.updateBook)

			// Assert
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			}
			if tc.expectedBook != nil {
				assert.Equal(t, tc.expectedBook.Title, result.Title)
				assert.Equal(t, tc.expectedBook.Writer, result.Writer)
				assert.Equal(t, tc.expectedBook.Isbn, result.Isbn)
				assert.Equal(t, fmt.Sprintf("%d", tc.expectedBook.UserID), fmt.Sprintf("%d", result.UserID))
			}
		})
	}
}
