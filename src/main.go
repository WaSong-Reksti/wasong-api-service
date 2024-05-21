package main

import (
	// "errors"
	"context"
	"example/wasong-api-service/src/database"
	models "example/wasong-api-service/src/models"
	"example/wasong-api-service/src/routes"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	myUser := models.User{
		Username: "hakuna_matata",
		Password: "123456",
		Email:    "hakuna@example.com",
		Type:     "student",
	}
	fmt.Println(myUser)
	ctx := context.Background()
	firestoreClient, err := database.InitializeFirestoreClient(&ctx)
	if err != nil {
		log.Fatalf("error initializing Firestore client: %v", err)
		return
	}
	fmt.Println("Successfuly connected to firestore client")
	defer firestoreClient.Close()

	routes.InitializeUserRoutes(ctx, router, firestoreClient)
	router.Run(":8080")

}
