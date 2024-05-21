// src/routes/users.go

package routes

import (
	"context"
	"example/wasong-api-service/src/models"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func InitializeUserRoutes(ctx context.Context, r *gin.Engine, firestoreClient *firestore.Client) {
	r.GET("/api/users", func(c *gin.Context) {
		users, err := models.GetUsersFromFirestore(ctx, firestoreClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	})

	r.GET("/api/users/:id", func(c *gin.Context) {
		userID := c.Param("id")

		user, err := models.GetUserById(ctx, firestoreClient, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	})

	r.PATCH("/api/users/:id", func(c *gin.Context) {
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

	r.DELETE("/api/users/:id", func(c *gin.Context) {
		userID := c.Param("id")

		if err := models.DeleteUserByID(ctx, firestoreClient, userID); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, "User "+userID+" deleted successfully")
	})
}
