package data

import "context"

type UsersManger interface {
	GetUserByID(ctx context.Context, userId int) (*User, error)
	SaveUser(ctx context.Context, user *User) error
}
