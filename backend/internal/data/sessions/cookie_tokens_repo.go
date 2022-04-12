package sessions

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log"
	"time"
)

type Token struct {
	UserID     int
	Expiry     time.Time
	Attributes map[string]string
}

//go:embed sql/create_token.sql
var createTokenQuery string

//go:embed sql/read_token.sql
var readTokenQuery string

//go:embed sql/revoke_token.sql
var revokeTokenQuery string

var ErrTokenNotFound = errors.New("token not found")
var ErrTokenIsNotActive = errors.New("token is not active")
var ErrTokenIsExpired = errors.New("token is expired")
var ErrNoTokenWasDeleted = errors.New("no token was deleted")

type TokensRepo struct {
	db *sql.DB
}

func NewTokensRepo(db *sql.DB) TokensRepo {
	return TokensRepo{
		db: db,
	}
}

type tokenDBRow struct {
	uuid   string
	userId int
	expiry time.Time
	active bool
	attrs  string
}

func (tr TokensRepo) Create(ctx context.Context, token Token) (string, error) {
	tUuid := uuid.New().String()
	tokenAttrsString, err := json.Marshal(token.Attributes)
	if err != nil {
		log.Println("can't marshal token attributes: ", token.Attributes)
		tokenAttrsString = nil
	}

	tdbr := tokenDBRow{
		uuid:   tUuid,
		userId: token.UserID,
		expiry: token.Expiry,
		active: true,
		attrs:  string(tokenAttrsString),
	}

	_, err = tr.db.ExecContext(
		ctx,
		createTokenQuery,
		tdbr.uuid,
		tdbr.userId,
		tdbr.expiry,
		tdbr.active,
		tdbr.attrs,
	)
	if err != nil {
		return "", err
	}

	return tUuid, nil
}

func (tr TokensRepo) Read(ctx context.Context, tokenId string) (Token, error) {
	row := tr.db.QueryRowContext(ctx, readTokenQuery, tokenId)
	if row.Err() == sql.ErrNoRows {
		return Token{}, ErrTokenNotFound
	} else if row.Err() != nil {
		return Token{}, row.Err()
	}

	tdbr := tokenDBRow{}
	err := row.Scan(&tdbr.uuid, &tdbr.userId, &tdbr.active, &tdbr.expiry, &tdbr.attrs)
	if err != nil {
		return Token{}, err
	}
	if !tdbr.active {
		return Token{}, ErrTokenIsNotActive
	}
	if time.Now().After(tdbr.expiry) {
		return Token{}, ErrTokenIsExpired
	}

	token := Token{
		UserID:     tdbr.userId,
		Expiry:     tdbr.expiry,
		Attributes: nil,
	}

	if len(tdbr.attrs) > 0 {
		err = json.Unmarshal([]byte(tdbr.attrs), &token.Attributes)
		if err != nil {
			log.Printf("can't parse token attrs, uuid(%s): %v\n", tdbr.uuid, err)
		}
	}

	return token, nil
}

func (tr TokensRepo) Revoke(ctx context.Context, tokenId string) error {
	r, err := tr.db.ExecContext(ctx, revokeTokenQuery, tokenId)
	if err != nil {
		return err
	}
	ra, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if ra == 0 {
		return ErrNoTokenWasDeleted
	}
	return err
}
