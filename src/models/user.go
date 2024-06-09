package models

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

/*
User Model
*/
type User struct {
	ID         string                 `firestore:"-"`
	Password   string                 `firestore:"Password,omitempty"`
	Email      string                 `firestore:"Email,omitempty"`
	Type       string                 `firestore:"Type,omitempty"`
	Name       string                 `firestore:"Name,omitempty"`
	About      string                 `firestore:"About,omitempty"`
	UID        string                 `firestore:"UID,omitempty"`
	StudentRef *firestore.DocumentRef `firestore:"Student,omitempty"`
	Student    *Students              `firestore:"-"`
}

func CreateUser(ctx context.Context, client *firestore.Client, user User) (string, *User, error) {
	if user.Type == "student" {
		// Create a new document in the "Students" collection
		studentDocRef, _, err := client.Collection("Students").Add(ctx, map[string]interface{}{})
		if err != nil {
			return "", nil, fmt.Errorf("failed to create student document: %v", err)
		}
		// Assign the reference of the created student document to the StudentRef field
		user.StudentRef = studentDocRef
	}

	// Create the user document in the "User" collection
	docRef, _, err := client.Collection("User").Add(ctx, user)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create user document: %v", err)
	}
	createdUser := user
	createdUser.ID = docRef.ID

	return docRef.ID, &createdUser, nil

}

func UpdateUserByID(ctx context.Context, client *firestore.Client, userID string, keyValue map[string]interface{}) (map[string]interface{}, error) {
	updates := make([]firestore.Update, 0)

	// Construct firestore.Update slice based on the keyValue map
	for key, value := range keyValue {
		updates = append(updates, firestore.Update{
			Path:  key,
			Value: value,
		})
	}

	_, err := client.Collection("User").Doc(userID).Update(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	docSnapshot, err := client.Collection("User").Doc(userID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated user data: %v", err)
	}

	return docSnapshot.Data(), nil

}

func DeleteUserByID(ctx context.Context, firestoneClient *firestore.Client, userID string) error {
	_, err := firestoneClient.Collection("User").Doc(userID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete: %v", err)
	}

	return err
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

func GetUserById(ctx context.Context, firestoreClient *firestore.Client, userID string) (*User, error) {
	dsnap, err := firestoreClient.Collection("User").Doc(userID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user document: %v", err)
	}

	var user User
	if err := dsnap.DataTo(&user); err != nil {
		return nil, fmt.Errorf("failed to convert data to User struct: %v", err)
	}

	return &user, nil

}

func getUser(ctx context.Context, userRef *firestore.DocumentRef) (*User, error) {
	// Fetch user document using the document reference
	docSnapshot, err := userRef.Get(ctx)
	if err != nil {
		return &User{}, fmt.Errorf("failed to get user document: %v", err)
	}

	// Convert data to User struct
	var user User
	if err := docSnapshot.DataTo(&user); err != nil {
		return &User{}, fmt.Errorf("failed to convert data to User struct: %v", err)
	}

	if user.StudentRef != nil {
		var student Students
		studentDoc, err := user.StudentRef.Get(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get student document: %v", err)
		}
		if err := studentDoc.DataTo(&student); err != nil {
			return nil, fmt.Errorf("failed to convert student data to Students struct: %v", err)
		}
		user.Student = &student
	}

	// Set the ID field of the user
	user.ID = docSnapshot.Ref.ID

	return &user, nil
}

func GetUserByUID(ctx context.Context, firestoreClient *firestore.Client, UID string) (*User, error) {
	iter := firestoreClient.Collection("User").Where("UID", "==", UID).Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user document: %v", err)
	}

	if err == iterator.Done {

	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user document: %v", err)
	}

	var user User
	if err := doc.DataTo(&user); err != nil {
		return nil, fmt.Errorf("failed to convert data to User struct: %v", err)
	}
	if user.StudentRef != nil {
		var student Students
		studentDoc, err := user.StudentRef.Get(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get student document: %v", err)
		}
		if err := studentDoc.DataTo(&student); err != nil {
			return nil, fmt.Errorf("failed to convert student data to Students struct: %v", err)
		}
		user.Student = &student
	}

	// Set the ID field of the user
	user.ID = doc.Ref.ID

	return &user, nil
}
