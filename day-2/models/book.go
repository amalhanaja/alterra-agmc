package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title  string `json:"title"`
	Isbn   string `json:"isbn"`
	Writer string `json:"writer"`
}
