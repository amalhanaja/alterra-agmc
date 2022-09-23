package response

type SuccessResponse[T any] struct {
	Status uint `json:"status"`
	Data   T    `json:"data"`
}
