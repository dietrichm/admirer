package services

// Service is the external service interface.
type Service interface {
	CreateAuthURL() string
	Authenticate(code string)
	GetUsername() string
}
