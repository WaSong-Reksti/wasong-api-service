package main

import (
	// "errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"context"
	"example/wasong-api-service/src/database"
	models "example/wasong-api-service/src/models"
	"fmt"
	"log"
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

	router.GET("/api/users", func(c *gin.Context) {
		users, err := models.GetUsersFromFirestore(ctx, firestoreClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	})

	router.Run(":8080")
	
}
