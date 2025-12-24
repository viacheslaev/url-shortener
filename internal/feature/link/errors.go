package link

var (
	ErrNotFound    = errorString("link not found")
	ErrLinkExpired = errorString("link expired")
)

type errorString string

func (e errorString) Error() string { return string(e) }
