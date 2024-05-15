package main

import (
	// "errors"
	"context"
	"example/wasong-api-service/src/database"
	models "example/wasong-api-service/src/models"
	"fmt"
	"log"
	"net/http"

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

	router.GET("/api/users", func(c *gin.Context) {
		users, err := models.GetUsersFromFirestore(ctx, firestoreClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	})

	router.GET("/api/users/:id", func(c *gin.Context) {
		userID := c.Param("id")

		user, err := models.GetUserById(ctx, firestoreClient, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	})

	router.PATCH("/api/users/:id", func(c *gin.Context) {
		userID := c.Param("id")

		var updateData map[string]interface{}
		if err := c.BindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		data, err := models.UpdateUserByID(ctx, firestoreClient, userID, updateData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, data)
	})

	router.DELETE("/api/users/:id", func(c *gin.Context) {
		userID := c.Param("id")

		if err := models.DeleteUserByID(ctx, firestoreClient, userID); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, "User " + userID + " deleted successfully")
	})

	router.Run(":8080")

}
