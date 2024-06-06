// src/routes/courses.go

package routes

import (
	"context"
	"example/wasong-api-service/src/models"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func InitializeCourseRoutes(ctx context.Context, r *gin.Engine, firestoreClient *firestore.Client) {
	// Add your course routes here
	r.GET("/api/courses", func(c *gin.Context) {
		courses, err := models.JoinCourseWithInstructor(ctx, firestoreClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, courses)
	})

	r.GET("api/courses/:id", func(c *gin.Context) {
		courseID := c.Param("id")
		course, err := models.GetCoursesById(ctx, firestoreClient, courseID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, course)
	})

	r.POST("/api/courses", func(c *gin.Context) {
		var requestBody struct {
			Name         string `json:"name"`
			Description  string `json:"description"`
			Instrument   string `json:"instrument"`
			InstructorID string `json:"instructor_id"`
		}
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse JSON body"})
			return
		}

		// Validate required fields
		if requestBody.Name == "" || requestBody.Description == "" || requestBody.Instrument == "" || requestBody.InstructorID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing required fields"})
			return
		}

		// Create Course struct
		newCourse := models.Course{
			Name:         requestBody.Name,
			Description:  requestBody.Description,
			Instrument:   requestBody.Instrument,
			InstructorID: firestoreClient.Collection("Instructor").Doc(requestBody.InstructorID),
		}

		// Create Course in Firestore
		createdCourseID, createdCourse, err := models.CreateCourse(ctx, firestoreClient, &newCourse)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return the created Course
		c.JSON(http.StatusCreated, gin.H{"id": createdCourseID, "course": createdCourse})
	})

	r.PATCH("/api/courses/:id", func(c *gin.Context) {
		courseID := c.Param("id")

		var updateData map[string]interface{}

		if err := c.BindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		data, err := models.UpdateCourse(ctx, firestoreClient, courseID, updateData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, data)

	})

	r.DELETE("/api/courses/:id", func(c *gin.Context) {
		courseID := c.Param("id")

		if err := models.DeleteCourseByID(ctx, firestoreClient, courseID); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, "Course "+courseID+" deleted successfully")
	})

}
