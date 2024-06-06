// src/models/instructor.go

package models

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

/*
Students Model
*/
type Students struct {
	ID       string                 `firestore:"-"`
	UserID   *firestore.DocumentRef `firestore:"User,omitempty"`
	User     *User                  `firestore:"-"`
	CourseID *firestore.DocumentRef `firestore:"Course,omitempty"`
	Course   *Course                `firestore:"-"`
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

func GetStudents(ctx context.Context, studentsRef *firestore.DocumentRef) (*Students, error) {
	// Fetch student document using the document reference
	docSnapshot, err := studentsRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get students document: %v", err)
	}

	// Convert data to Students struct
	var student Students
	if err := docSnapshot.DataTo(&student); err != nil {
		return nil, fmt.Errorf("failed to convert data to Students struct: %v", err)
	}

	student.ID = docSnapshot.Ref.ID

	return &student, nil
}
