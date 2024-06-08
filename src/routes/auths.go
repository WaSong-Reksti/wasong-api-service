package routes

import (
	"context"
	"encoding/json"
	"example/wasong-api-service/src/auth"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func InitializeAuthRoutes(c context.Context, r *gin.Engine, firebaseAuth *auth.FirebaseAuth, firestoreClient *firestore.Client) {
	r.POST("/api/register", func(ctx *gin.Context) {
		var requestData struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := ctx.BindJSON(&requestData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		user, userDoc, err := firebaseAuth.RegisterUser(requestData.Email, requestData.Password, firestoreClient)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"user_id": user.UID, "user_record": userDoc})
	})

	r.GET("/api/login", func(ctx *gin.Context) {
		var requestData struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := ctx.BindJSON(&requestData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		resp, err := firebaseAuth.LoginUser(requestData.Email, requestData.Password)
		if err != nil {
			var errorResp map[string]interface{}
			if jsonErr := json.Unmarshal([]byte(err.Error()), &errorResp); jsonErr == nil {
				ctx.JSON(http.StatusInternalServerError, errorResp)
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		ctx.JSON(http.StatusOK, resp)
	})
	log.Println("Auth route initialized")
}
