package services

// Service is the external service interface.
type Service interface {
	Name() string
	CreateAuthURL() string
	Authenticate(code string) error
	GetUsername() (string, error)
}
