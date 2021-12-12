package authentication

import (
	"context"
	"github.com/netscrn/gocookieauth/data/sessions"
	"github.com/netscrn/gocookieauth/data/users"
)

type contextAuthKey int

const (
	userKey contextAuthKey = iota + 1
	tokenKey
)

func ContextWithUser(ctx context.Context, user users.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}
func UserFromContext(ctx context.Context) (users.User, bool) {
	user, ok := ctx.Value(userKey).(users.User)
	return user, ok
}

func ContextWithToken(ctx context.Context, token sessions.Token) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}
func TokenFromContext(ctx context.Context) (sessions.Token, bool) {
	token, ok := ctx.Value(tokenKey).(sessions.Token)
	return token, ok
}
