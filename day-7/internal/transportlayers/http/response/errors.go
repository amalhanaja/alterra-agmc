package response

type ErrorResponse struct {
	Status  uint   `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
