package datasources

import (
	dsModels "alterra-agmc-day-7/internal/datasources/models"
	"alterra-agmc-day-7/internal/models"
	"alterra-agmc-day-7/internal/repositories"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type bookMongoDataSource struct {
	db *mongo.Database
}

// Create implements repositories.BookRepository
func (ds *bookMongoDataSource) Create(ctx context.Context, book *models.Book) (*models.Book, error) {
	all, err := ds.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	utcNow := time.Now().UTC()
	book.ID = uint((len(all) + 1))
	book.CreatedAt = utcNow
	book.UpdatedAt = utcNow
	mongoModel := &dsModels.BookMongoModel{
		ID:        book.ID,
		Title:     book.Title,
		Isbn:      book.Isbn,
		Writer:    book.Writer,
		CreatedAt: book.CreatedAt,
		UpdatedAt: book.UpdatedAt,
		UserID:    book.UserID,
	}
	res, err := ds.collections().InsertOne(ctx, mongoModel)
	if err != nil {
		return nil, err
	}
	book.ID = uint(res.InsertedID.(int64))
	return book, nil
}

// DeleteByID implements repositories.BookRepository
func (ds *bookMongoDataSource) DeleteByID(ctx context.Context, id uint) error {
	_, err := ds.collections().DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// FindAll implements repositories.BookRepository
func (ds *bookMongoDataSource) FindAll(ctx context.Context) ([]*models.Book, error) {
	results := make([]*models.Book, 0)
	cur, err := ds.collections().Find(ctx, bson.M{})
	if err != nil {
		return results, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		mongoModel := &dsModels.BookMongoModel{}
		err := cur.Decode(mongoModel)
		if err != nil {
			return results, err
		}
		book := &models.Book{
			ID:        mongoModel.ID,
			Title:     mongoModel.Title,
			Isbn:      mongoModel.Isbn,
			Writer:    mongoModel.Writer,
			CreatedAt: mongoModel.CreatedAt,
			UpdatedAt: mongoModel.UpdatedAt,
			UserID:    mongoModel.UserID,
		}
		results = append(results, book)
	}
	return results, nil
}

// FindByID implements repositories.BookRepository
func (ds *bookMongoDataSource) FindByID(ctx context.Context, id uint) (*models.Book, error) {
	res := ds.collections().FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		return nil, res.Err()
	}
	mongoModel := &dsModels.BookMongoModel{}
	err := res.Decode(mongoModel)
	if err != nil {
		return nil, err
	}
	book := &models.Book{
		ID:        mongoModel.ID,
		Title:     mongoModel.Title,
		Isbn:      mongoModel.Isbn,
		Writer:    mongoModel.Writer,
		CreatedAt: mongoModel.CreatedAt,
		UpdatedAt: mongoModel.UpdatedAt,
		UserID:    mongoModel.UserID,
	}
	return book, nil
}

// Update implements repositories.BookRepository
func (ds *bookMongoDataSource) Update(ctx context.Context, book *models.Book) (*models.Book, error) {
	book.UpdatedAt = time.Now().UTC()
	mongoModel := &dsModels.BookMongoModel{
		Title:     book.Title,
		Isbn:      book.Isbn,
		Writer:    book.Writer,
		CreatedAt: book.CreatedAt,
		UpdatedAt: book.UpdatedAt,
		UserID:    book.UserID,
	}
	_, err := ds.collections().UpdateByID(ctx, book.ID, mongoModel)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (ds *bookMongoDataSource) collections() *mongo.Collection {
	return ds.db.Collection("books")
}

func NewBookMongoDataSource(db *mongo.Database) repositories.BookRepository {
	return &bookMongoDataSource{db: db}
}
