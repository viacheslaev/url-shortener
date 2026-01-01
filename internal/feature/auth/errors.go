package auth

type errorString string

func (e errorString) Error() string { return string(e) }

var (
	ErrInvalidCredentials = errorString("invalid credentials")
	ErrUnauthorized       = errorString("unauthorized")
)
