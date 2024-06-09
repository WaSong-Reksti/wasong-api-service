package routes

import (
	"context"
	"encoding/json"
	"example/wasong-api-service/src/auth"
	"log"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func InitializeAuthRoutes(c context.Context, r *gin.Engine, firebaseAuth *auth.FirebaseAuth, firestoreClient *firestore.Client) {
	r.POST("/api/register", func(ctx *gin.Context) {
		var requestData struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			Name     string `json:name`
		}
		if err := ctx.BindJSON(&requestData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		user, userDoc, err := firebaseAuth.RegisterUser(requestData.Email, requestData.Password, requestData.Name, firestoreClient)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"user_id": user.UID, "user_record": userDoc})
	})

	r.POST("/api/login", func(ctx *gin.Context) {
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

	r.GET("/api/session", func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Malformed authorization header"})
			return
		}
		token, err := firebaseAuth.VerifyToken(c, tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "error verifying ID token: " + err.Error()})
			return
		}
		resp := map[string]interface{}{
			"uid":      token.UID,
			"issuer":   token.Issuer,
			"expires":  token.Expires,
			"issuedAt": token.IssuedAt,
			"authTime": token.AuthTime,
			"firebase": token.Firebase,
		}
		ctx.JSON(http.StatusOK, resp)
	})

	r.POST("/api/logout", func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Malformed authorization header"})
			return
		}

		err := firebaseAuth.RevokeToken(c, tokenString)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke refresh tokens"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
	})
	log.Println("Auth route initialized")
}
