package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"example/wasong-api-service/src/models"
	"io"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"golang.org/x/crypto/bcrypt"

	"github.com/joho/godotenv"
)

func init() {
	// Load .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
}

type FirebaseAuth struct {
	ctx        context.Context
	app        *firebase.App
	authClient *auth.Client
}

type Credentials struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

// LoginResponse represents the response received after login.
type LoginResponse struct {
	IDToken      string `json:"idToken"`
	Email        string `json:"email"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalID      string `json:"localId"`
	Registered   bool   `json:"registered"`
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
func (fa *FirebaseAuth) RegisterUser(email, password, name string, firestoreClient *firestore.Client) (*auth.UserRecord, *models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error hashing password: %v", err)
		return nil, nil, err
	}
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password)
	user, err := fa.authClient.CreateUser(fa.ctx, params)
	if err != nil {
		log.Printf("error registering user: %v", err)
		return nil, nil, err
	}
	_, userRecord, err := models.CreateUser(fa.ctx, firestoreClient, models.User{
		Email:    user.Email,
		Name:     name,
		UID:      user.UID,
		Password: string(hashedPassword),
	})
	if err != nil {
		return nil, nil, err
	}
	return user, userRecord, nil
}

// VerifyIDToken verifies the ID token.
func (fa *FirebaseAuth) LoginUser(email, password string) (*LoginResponse, error) {
	// Prepare credentials
	creds := Credentials{
		Email:             email,
		Password:          password,
		ReturnSecureToken: true,
	}

	// Convert credentials to JSON
	jsonData, err := json.Marshal(creds)
	if err != nil {
		log.Printf("error encoding credentials: %v", err)
		return nil, err
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY not found in environment variables")
	}

	// Send POST request to Firebase Auth REST API
	resp, err := http.Post("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key="+apiKey, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("error sending login request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %v", err)
		return nil, err
	}

	// Check if request was successful
	if resp.StatusCode != http.StatusOK {
		log.Printf("login request failed: %s", string(body))

		return nil, errors.New(string(body))
	}

	// Decode response
	var loginResp LoginResponse
	err = json.Unmarshal(body, &loginResp)
	if err != nil {
		log.Printf("error decoding login response: %v", err)
		return nil, err
	}

	return &loginResp, nil
}

func (fa *FirebaseAuth) VerifyToken(ctx context.Context, idToken string) (*auth.Token, error) {
	token, err := fa.authClient.VerifyIDTokenAndCheckRevoked(ctx, idToken)
	if err != nil {
		if err.Error() == "ID token has been revoked" {
			// Token is revoked. Inform the user to reauthenticate or signOut() the user.
			log.Printf("the ID Token has been revoked")
			return nil, err
		} else {
			log.Printf("error: %v", err)
		}
		return nil, err
	}
	return token, err
}

func (fa *FirebaseAuth) RevokeToken(ctx context.Context, idToken string) error {
	token, err := fa.authClient.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return err
	}
	err = fa.authClient.RevokeRefreshTokens(ctx, token.UID)
	if err != nil {
		return err
	}

	return nil

}
