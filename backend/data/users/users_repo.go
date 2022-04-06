package users

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"github.com/go-sql-driver/mysql"
)

//go:embed sql/get_user_by_id.sql
var getUserByIdQuery string

//go:embed sql/get_user_by_email.sql
var getUserByEmailQuery string

//go:embed sql/save_user.sql
var saveUserQuery string

var ErrNoUserFound = errors.New("no user is found")
var ErrSuchEmailIsAlreadyExists = errors.New("such email is already exists")

type UsersRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UsersRepo {
	return UsersRepo{
		db: db,
	}
}

func (ur UsersRepo) GetUserByEmail(ctx context.Context, email string) (User, error) {
	u := User{}
	row := ur.db.QueryRowContext(ctx, getUserByEmailQuery, email)
	if row.Err() != nil {
		return u, row.Err()
	}

	err := row.Scan(&u.Id, &u.Email, &u.PasswordDigest, &u.IsConfirmed)
	if err == sql.ErrNoRows { // todo: use errors.Is()
		return u, ErrNoUserFound
	} else if err != nil {
		return u, err
	}

	return u, nil
}

func (ur UsersRepo) GetUserByID(ctx context.Context, userId int) (User, error) {
	u := User{}
	row := ur.db.QueryRowContext(ctx, getUserByIdQuery, userId)
	if row.Err() != nil {
		return u, row.Err()
	}

	err := row.Scan(&u.Id, &u.Email, &u.PasswordDigest, &u.IsConfirmed)
	if err == sql.
	{
		return u, ErrNoUserFound
	} else if err != nil {
		return u, err
	}

	return u, nil
}

func (ur UsersRepo) SaveUser(ctx context.Context, user *User) error {
	res, err := ur.db.ExecContext(ctx, saveUserQuery, user.PasswordDigest, user.Email, user.IsConfirmed)
	if err != nil {
		var mysqlError *mysql.MySQLError
		if errors.As(err, &mysqlError) {
			if mysqlError.Number == 1062 {
				return ErrSuchEmailIsAlreadyExists
			}
		}
		return err
	}
	savedUserId, err := res.LastInsertId()
	if err != nil {
		return err
	}
	user.Id = int(savedUserId)
	return nil
}
