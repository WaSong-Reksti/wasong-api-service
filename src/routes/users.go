// src/routes/users.go

package routes

import (
	"context"
	"example/wasong-api-service/src/models"
	"fmt"
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

	r.GET("/api/users/uid=:uid", func(c *gin.Context) {
		uid := c.Param("uid")

		user, err := models.GetUserByUID(ctx, firestoreClient, uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	})

	r.GET("/api/students/uid=:uid", func(c *gin.Context) {
		uid := c.Param("uid")
		user, err := models.GetUserByUID(ctx, firestoreClient, uid)
		if err != nil {
			fmt.Printf("Error coy: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		student, err := models.GetStudents(ctx, user.StudentRef)
		if err != nil {
			fmt.Printf("Error coy: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, student)
	})

	r.POST("api/enroll", func(c *gin.Context) {

		var requestBody struct {
			UID      string `json:"uid"`
			CourseID string `json:"course_id"`
		}
		if err := c.BindJSON(&requestBody); err != nil {
			fmt.Printf("error: " + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse JSON body"})
			return
		}

		uid := requestBody.UID
		courseId := requestBody.CourseID

		courseRef := firestoreClient.Collection("Course").Doc(courseId)
		user, err := models.GetUserByUID(ctx, firestoreClient, uid)
		if err != nil {
			fmt.Printf("Error coy 1: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = models.AddStudentToCourse(ctx, firestoreClient, user, courseRef)
		if err != nil {
			fmt.Printf("Error coy 2: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		courseDoc, err := courseRef.Get(ctx)
		if err != nil {
			fmt.Printf("Error retrieving course: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		courseName, err := courseDoc.DataAt("Name")
		if err != nil {
			fmt.Printf("Error retrieving course name: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		message := fmt.Sprintf("Student %s Enrolled To Class %s", user.Name, courseName)
		c.JSON(http.StatusOK, gin.H{"message": message})
	})

	fmt.Println("Initialize users route")
}
