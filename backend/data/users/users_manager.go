package users

import "context"

type Manger interface {
	GetUserByID(ctx context.Context, userId int) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	SaveUser(ctx context.Context, user *User) error
}
