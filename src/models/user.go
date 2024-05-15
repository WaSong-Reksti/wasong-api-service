package models

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

/*
User Model
*/
type User struct {
	Username string `json: username`
	Password string `json: password`
	Email    string `json: email`
	Type     string `json: type`
	Name     string `json: name`
	About    string `json: about`
}

func (u *User) createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (u *User) updateUser(c *gin.Context) {

}

func GetUsersFromFirestore(ctx context.Context, firestoreClient *firestore.Client) ([]User, error) {
	iter := firestoreClient.Collection("User").Documents(ctx)
	var users []User
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate: %v", err)
		}
		var user User
		if err := doc.DataTo(&user); err != nil {
			return nil, fmt.Errorf("failed to convert data to User struct: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}
