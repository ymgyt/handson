package context

import (
	"context"

	"github.com/YmgchiYt/golang-handson/web_development_with_go/models"
)

const (
	userKey privateKey = "user"
)

type privateKey string

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	if maybeUser := ctx.Value(userKey); maybeUser != nil {
		if user, ok := maybeUser.(*models.User); ok {
			return user
		}
	}
	return nil
}
