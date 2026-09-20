package http

import "context"

const (
	ctxValueUserId = "user_id"
)

func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxValueUserId, id)
}
