package dto

type CreateBookPayload struct {
	Title  string `json:"title" validate:"required"`
	Isbn   string `json:"isbn" validate:"required"`
	Writer string `json:"writer" validate:"required"`
}

type UpdateBookPayload struct {
	Title  string `json:"title,omitempty" validate:"omitempty"`
	Isbn   string `json:"isbn,omitempty" validate:"omitempty"`
	Writer string `json:"writer,omitempty" validate:"omitempty"`
}
