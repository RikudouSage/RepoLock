package http

import (
	"context"

	"github.com/google/uuid"
	"go.chrastecky.dev/repolock/entity"
)

const (
	ctxValueUserId = "user_id"
	ctxValueUser   = "user"
)

func WithUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, ctxValueUserId, id)
}

func GetUserID(ctx context.Context) (userID uuid.UUID, ok bool) {
	userID, ok = ctx.Value(ctxValueUserId).(uuid.UUID)
	return
}

func WithUser(ctx context.Context, user *entity.User) context.Context {
	return context.WithValue(ctx, ctxValueUser, user)
}

func GetUser(ctx context.Context) (user *entity.User, ok bool) {
	user, ok = ctx.Value(ctxValueUser).(*entity.User)
	return
}

func MustGetUser(ctx context.Context) *entity.User {
	user, ok := GetUser(ctx)
	if !ok {
		panic("user not found in context")
	}

	return user
}
