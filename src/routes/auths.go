package routes

import (
	"context"
	"example/wasong-api-service/src/auth"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitializeAuthRoutes(c context.Context, r *gin.Engine, firebaseAuth *auth.FirebaseAuth) {
	r.POST("/api/register", func(ctx *gin.Context) {
		var requestData struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := ctx.BindJSON(&requestData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		user, err := firebaseAuth.RegisterUser(requestData.Email, requestData.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"user_id": user.UID})
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
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, resp)
	})
	log.Println("Auth route initialized")
}
