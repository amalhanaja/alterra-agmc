package services

type ErrUnauthorized struct{}

func (e ErrUnauthorized) Error() string {
	return "unauthorized"
}
