// src/models/instructor.go

package models

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/firestore"
)

/*
Students Model
*/
type Students struct {
	ID      string   `firestore:"-"`
	Courses []Course `firestore:"-"`
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

	//iterate through Courses
	coursesData, err := docSnapshot.DataAt("Courses")
	if err != nil {
		return nil, fmt.Errorf("failed to get Courses data: %v", err)
	}
	coursesSlice, ok := coursesData.([]interface{})
	if !ok {
		return nil, fmt.Errorf("courses data is not an array")
	}

	for _, courseData := range coursesSlice {
		var course Course
		courseMap, ok := courseData.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("course data is not a map")
		}

		// Convert map to Course struct
		if err := mapToStruct(courseMap, &course); err != nil {
			return nil, fmt.Errorf("failed to convert course data to Course struct: %v", err)
		}

		student.Courses = append(student.Courses, course)
	}

	return &student, nil
}

func mapToStruct(m map[string]interface{}, s interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal map: %v", err)
	}
	if err := json.Unmarshal(data, s); err != nil {
		return fmt.Errorf("failed to unmarshal data to struct: %v", err)
	}
	return nil
}
