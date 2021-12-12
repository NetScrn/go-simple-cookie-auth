package sessions

import "context"

type TokenManager interface {
	Read(ctx context.Context, tokenId string) (Token, error)
	Create(ctx context.Context, token Token) (string, error)
	Revoke(ctx context.Context, tokenId string) error
}
