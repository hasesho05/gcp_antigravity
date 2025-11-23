package auth

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/cockroachdb/errors"
)

func NewClient(ctx context.Context) (*auth.Client, error) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error initializing firebase app")
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error initializing auth client")
	}

	return client, nil
}
