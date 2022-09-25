package models

import "time"

type BookMongoModel struct {
	ID        uint      `bson:"_id"`
	Title     string    `bson:"title"`
	Isbn      string    `bson:"isbn"`
	Writer    string    `bson:"writer"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	UserID    uint      `bson:"user_id"`
}
