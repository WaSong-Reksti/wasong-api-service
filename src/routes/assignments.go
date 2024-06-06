// src/routes/assignments.go

package routes

import (
	"context"
	"example/wasong-api-service/src/models"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func InitializeAssignmentsRoutes(ctx context.Context, r *gin.Engine, firestoreClient *firestore.Client) {
	// Add your assignments routes here
	r.GET("/api/assignments/:courseId", func(c *gin.Context) {
		courseID := c.Param("courseId")
		assignments, err := models.GetCourseAssignments(ctx, firestoreClient, courseID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, assignments)
	})

	r.GET("/api/assignments/:courseId/:assignmentId", func(c *gin.Context) {
		courseID := c.Param("courseId")
		assignmentID := c.Param("assignmentId")
		assignments, err := models.GetCourseAssignmentById(ctx, firestoreClient, courseID, assignmentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, assignments)
	})

	r.POST("/api/assignments/:courseId", func(c *gin.Context) {
		courseID := c.Param("courseId")

		var requestBody struct {
			Name        string `json:"name"`
			StartTime   string `json:"startTime"`
			EndTime     string `json:"endTime"`
			Description string `json:"description"`
		}
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse JSON body"})
			return
		}

		// Validate required fields
		if requestBody.Name == "" || requestBody.StartTime == "" || requestBody.EndTime == "" || requestBody.Description == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing required fields"})
			return
		}

		// Create Assignment struct
		startTime, err1 := time.Parse("2006-01-02T15:04:05Z07:00", requestBody.StartTime)
		endTime, err2 := time.Parse("2006-01-02T15:04:05Z07:00", requestBody.EndTime)
		if err1 != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot parse into time.Time format"})
			return
		}
		newAssignment := models.Assignment{
			Name:        requestBody.Name,
			StartTime:   startTime,
			EndTime:     endTime,
			Description: requestBody.Description,
		}

		// Create Assignment in Firestore
		createdAssignmentID, createdAssignment, err := models.CreateAssignment(ctx, firestoreClient, courseID, &newAssignment)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return the created Assignment
		c.JSON(http.StatusCreated, gin.H{"id": createdAssignmentID, "assignment": createdAssignment})
	})

	r.DELETE("/api/assignments/:courseId/:assignmentId", func(c *gin.Context) {
		courseID := c.Param("courseId")
		assignmentID := c.Param("assignmentId")

		if err := models.DeleteAssignmentByID(ctx, firestoreClient, courseID, assignmentID); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, "Assignment "+assignmentID+" deleted successfully")
	})

	fmt.Println("Initialize assignments route")
}
