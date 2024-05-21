// src/routes/courses.go

package routes

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func InitializeCourseRoutes(ctx context.Context, r *gin.Engine, firestoreClient *firestore.Client) {
	// Add your course routes here
	r.GET("/api/courses", func(c *gin.Context) {
		users, err = models.
			c.JSON(http.StatusOK, gin.H{"message": "List of courses"})
	})
	// Other course routes
}
