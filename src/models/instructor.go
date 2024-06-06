// src/models/instructor.go

package models

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

/*
Instructor Model
*/
type Instructor struct {
	ID     string                 `firestore:"-"`
	UserID *firestore.DocumentRef `firestore:"User,omitempty"`
	User   *User                  `firestore:"-"`
}

// func JoinInstructorWithUser(ctx context.Context, firestore *firestore.Client, instructor *Instructor) (*Instructor, error) {
//     // Fetch instructor document
//     // instructorDocRef := firestore.Collection("Instructor").Doc(instructor.ID)
//     instructor, err := GetInstructor(ctx, instructorDocRef)
//     if err != nil {
//         return nil, fmt.Errorf("failed to fetch instructor: %v", err)
//     }

//     // Fetch user data for the instructor
//     user, err := getUser(ctx, firestore, instructor.UserID)
//     if err != nil {
//         return nil, fmt.Errorf("failed to fetch user data for instructor: %v", err)
//     }
//     instructor.User = user

//     return instructor, nil
// }

func GetInstructor(ctx context.Context, instructorRef *firestore.DocumentRef) (*Instructor, error) {
	// Fetch instructor document using the document reference
	docSnapshot, err := instructorRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get instructor document: %v", err)
	}

	// Convert data to Instructor struct
	var instructor Instructor
	if err := docSnapshot.DataTo(&instructor); err != nil {
		return nil, fmt.Errorf("failed to convert data to Instructor struct: %v", err)
	}

	instructor.ID = docSnapshot.Ref.ID

	return &instructor, nil
}
