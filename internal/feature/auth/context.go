package auth

import "context"

type ctxKey string

const accountPublicIDKey ctxKey = "account_public_id"

// WithAccountPublicID injects the authenticated account_public_id into request context.
func WithAccountPublicID(ctx context.Context, publicID string) context.Context {
	return context.WithValue(ctx, accountPublicIDKey, publicID)
}

// AccountPublicIDFromContext retrieves the authenticated account public_id from context.
// Returns (id, true) when present and valid; ("", false) otherwise.
func AccountPublicIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(accountPublicIDKey)
	id, ok := v.(string)
	if !ok || id == "" {
		return "", false
	}
	return id, true
}
