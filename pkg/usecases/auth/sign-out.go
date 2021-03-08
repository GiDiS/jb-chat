package auth

import "context"

type SignOut interface {
	SignOut(ctx context.Context, token string) error
}

func (a *authImpl) SignOut(ctx context.Context, token string) error {
	// @todo remove sess
	return nil
}
