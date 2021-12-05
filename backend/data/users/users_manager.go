package users

import (
	"context"
)

type Manger interface {
	GetUserByID(ctx context.Context, userId int) (*User, error)
	SaveUser(ctx context.Context, user *User) error
}
