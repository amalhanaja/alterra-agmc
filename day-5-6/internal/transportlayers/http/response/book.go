package response

import "time"

type BookResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Isbn      string    `json:"isbn"`
	Writer    string    `json:"writer"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
