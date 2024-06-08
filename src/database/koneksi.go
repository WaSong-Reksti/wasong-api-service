package database

import (
	"context"
	"log"
	"path/filepath"
	"runtime"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func InitializeFirestoreClient(ctx *context.Context, firebaseApp *firebase.App) (*firestore.Client, error) {
	client, err := firebaseApp.Firestore(*ctx)
	if err != nil {
		log.Fatalf("error initializing Firestore client: %v", err)
		return nil, err
	}
	return client, nil
}

func InitializeFirebaseApp(ctx context.Context) (*firebase.App, error) {
	_, currentFilePath, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(currentFilePath)))

	// Construct the path to serviceAccount.json relative to the project root
	serviceAccountPath := filepath.Join(projectRoot, "keys", "wasong-reksti-firebase-adminsdk-8zw1j-15e75775b3.json")

	sa := option.WithCredentialsFile(serviceAccountPath)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v", err)
		return nil, err
	}
	return app, nil
}
