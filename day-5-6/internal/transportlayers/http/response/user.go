package response

import "time"

type UserResponse struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
