package auth

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

type FirebaseAuth struct {
	ctx        context.Context
	app        *firebase.App
	authClient *auth.Client
}

// NewFirebaseAuth creates a new instance of FirebaseAuth.
func NewFirebaseAuth(ctx context.Context, app *firebase.App) (*FirebaseAuth, error) {

	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Printf("error initializing Firebase Auth client: %v", err)
		return nil, err
	}

	return &FirebaseAuth{
		ctx:        ctx,
		app:        app,
		authClient: authClient,
	}, nil
}

// RegisterUser registers a new user.
func (fa *FirebaseAuth) RegisterUser(email, password string) (*auth.UserRecord, error) {
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password)
	user, err := fa.authClient.CreateUser(fa.ctx, params)
	if err != nil {
		log.Printf("error registering user: %v", err)
		return nil, err
	}
	return user, nil
}

// VerifyIDToken verifies the ID token.
func (fa *FirebaseAuth) VerifyIDToken(idToken string) (*auth.Token, error) {
	token, err := fa.authClient.VerifyIDToken(fa.ctx, idToken)
	if err != nil {
		log.Printf("error verifying ID token: %v", err)
		return nil, err
	}
	return token, nil
}
