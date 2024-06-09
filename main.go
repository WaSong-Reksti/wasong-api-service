package main

import (
	// "errors"
	"context"
	"example/wasong-api-service/src/auth"
	"example/wasong-api-service/src/database"
	"example/wasong-api-service/src/routes"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	ctx := context.Background()
	firebaseApp, err := database.InitializeFirebaseApp(ctx)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v", err)
		return
	}
	fmt.Println("Successfuly initialize firebase app")
	firestoreClient, err := database.InitializeFirestoreClient(&ctx, firebaseApp)
	if err != nil {
		log.Fatalf("error initializing Firestore client: %v", err)
		return
	}
	fmt.Println("Successfuly connected to firestore client")
	defer firestoreClient.Close()
	firebaseAuth, err := auth.NewFirebaseAuth(ctx, firebaseApp)
	if err != nil {
		log.Fatalf("error initializing Firebase auth: %v", err)
		return
	}
	log.Println("Successfully initialize firebase auth")

	// Run assignment routes in a separate goroutine
	routes.InitializeAssignmentsRoutes(ctx, router, firestoreClient)
	routes.InitializeUserRoutes(ctx, router, firestoreClient)
	routes.InitializeCourseRoutes(ctx, router, firestoreClient)
	routes.InitializeAuthRoutes(ctx, router, firebaseAuth, firestoreClient)

	router.Run("localhost:8080")
	fmt.Println("Sucessfully run on localhost:8080")

}
