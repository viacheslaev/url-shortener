package account

type errorString string

func (e errorString) Error() string { return string(e) }

var (
	ErrEmailAlreadyExists = errorString("email already exists")
	ErrAccountNotFound    = errorString("account not found")
)
