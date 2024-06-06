package models

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

/*
Assignment Model
*/
type Assignment struct {
	ID          string        `firestore:"-"`
	Name        string        `firestore:"Name,omitempty"`
	StartTime   time.Time     `firestore:"StartTime,omitempty"`
	EndTime     time.Time     `firestore:"EndTime,omitempty"`
	Submissions []Submissions `firestore:"-"`
	Description string        `firestore:"Description,omitempty"`
}

type Submissions struct {
	attachmentPath string                 `firestore:"AttachmentPath,omitempty"`
	StudentID      *firestore.DocumentRef `firestore:"StudentID,omitempty"`
	Student        *Students              `firestore:"-"`
}

func GetCourseAssignmentById(ctx context.Context, firestoreClient *firestore.Client, courseId string, assignmentId string) (*Assignment, error) {
	assignmentDocRef := firestoreClient.Collection("Course").Doc(courseId).Collection("Assignments").Doc(assignmentId)
	assignmentDoc, err := assignmentDocRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignment document: %v", err)
	}

	var assignment Assignment
	if err := assignmentDoc.DataTo(&assignment); err != nil {
		return nil, fmt.Errorf("failed to convert data to Assignment struct: %v", err)
	}
	assignment.ID = assignmentDoc.Ref.ID

	// Fetch Submissions (if any)
	if assignment.Submissions != nil {
		for i, submission := range assignment.Submissions {
			if submission.StudentID != nil {
				studentDoc, err := submission.StudentID.Get(ctx)
				if err != nil {
					fmt.Printf("failed to fetch student data for submission %d in assignment %s: %v", i, assignment.ID, err)
					continue
				}
				var student Students
				if err := studentDoc.DataTo(&student); err != nil {
					fmt.Printf("failed to convert student data to Student struct for submission %d in assignment %s: %v", i, assignment.ID, err)
					continue
				}
				submission.Student = &student
				assignment.Submissions[i] = submission
			}
		}
	}

	return &assignment, nil
}

func GetCourseAssignments(ctx context.Context, firestoreClient *firestore.Client, courseId string) ([]Assignment, error) {
	// Fetch the course document
	courseDocRef := firestoreClient.Collection("Course").Doc(courseId)
	courseDoc, err := courseDocRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Course document: %v", err)
	}

	// Convert the course document data to a Course struct
	var course Course
	if err := courseDoc.DataTo(&course); err != nil {
		return nil, fmt.Errorf("failed to convert data to Course struct: %v", err)
	}
	course.ID = courseDoc.Ref.ID

	// Fetch assignments subcollection for the course
	assignmentsIter := courseDocRef.Collection("Assignments").Documents(ctx)
	var assignments []Assignment
	for {
		doc, err := assignmentsIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate over assignments: %v", err)
		}

		var assignment Assignment
		if err := doc.DataTo(&assignment); err != nil {
			return nil, fmt.Errorf("failed to convert data to Assignment struct: %v", err)
		}
		assignment.ID = doc.Ref.ID

		// Fetch Submissions (if any)
		if assignment.Submissions != nil {
			for i, submission := range assignment.Submissions {
				if submission.StudentID != nil {
					studentDoc, err := submission.StudentID.Get(ctx)
					if err != nil {
						log.Printf("failed to fetch student data for submission %d in assignment %s: %v", i, assignment.ID, err)
						continue
					}
					var student Students
					if err := studentDoc.DataTo(&student); err != nil {
						log.Printf("failed to convert student data to Student struct for submission %d in assignment %s: %v", i, assignment.ID, err)
						continue
					}
					submission.Student = &student
					assignment.Submissions[i] = submission
				}
			}
		}

		assignments = append(assignments, assignment)
	}

	course.Assignments = assignments
	return assignments, nil
}

func CreateAssignment(ctx context.Context, client *firestore.Client, courseId string, assignment *Assignment) (string, *Assignment, error) {
	docRef, _, err := client.Collection("Course").Doc(courseId).Collection("Assignments").Add(ctx, assignment)
	if err != nil {
		return "", nil, fmt.Errorf("error: %v", err)
	}
	createdAssignment := assignment
	createdAssignment.ID = docRef.ID
	return docRef.ID, createdAssignment, nil
}

func UpdateAssignment(ctx context.Context, firestoreClient *firestore.Client, assignmentID string, keyValue map[string]interface{}) (map[string]interface{}, error) {
	updates := make([]firestore.Update, 0)
	for key, value := range keyValue {
		updates = append(updates, firestore.Update{
			Path:  key,
			Value: value,
		})
	}

	_, err := firestoreClient.Collection("Assignment").Doc(assignmentID).Update(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update assignment: %v", err)
	}

	docSnap, err := firestoreClient.Collection("Assignment").Doc(assignmentID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated assignment: %v", err)
	}
	return docSnap.Data(), nil
}

func DeleteAssignmentByID(ctx context.Context, firestoreClient *firestore.Client, courseID string, assignmentID string) error {
	_, err := firestoreClient.Collection("Course").Doc(courseID).Collection("Assignments").Doc(assignmentID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete: %v", err)
	}

	return err

}

// func JoinCourseWithInstructor(ctx context.Context, firestore *firestore.Client) ([]Course, error) {
// 	// Step 1: Query Courses
// 	iter := firestore.Collection("Course").Documents(ctx)
// 	var courses []Course
// 	// Step 2: Fetch instructor data for each course
// 	for {
// 		doc, err := iter.Next()
// 		if err == iterator.Done {
// 			break
// 		}
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to iterate: %v", err)
// 		}
// 		var course Course
// 		if err := doc.DataTo(&course); err != nil {
// 			return nil, fmt.Errorf("failed to convert data to Course struct: %v", err)
// 		}
// 		course.ID = doc.Ref.ID

// 		// Fetch instructor data for the course
// 		instructorDocRef := course.InstructorID
// 		if instructorDocRef != nil {
// 			instructor, err := GetInstructor(ctx, instructorDocRef)
// 			if err != nil {
// 				return nil, fmt.Errorf("failed to fetch instructor data for course %s: %v", course.ID, err)
// 			}
// 			course.Instructor = instructor

// 			// fetch user data from instructor
// 			userDocRef := instructor.UserID
// 			if userDocRef != nil {
// 				user, err := getUser(ctx, userDocRef)
// 				if err != nil {
// 					return nil, fmt.Errorf("failed to fetch user data for instructor of course %s: %v", course.ID, err)
// 				}
// 				instructor.User = user
// 			}
// 		}

// 		courses = append(courses, course)
// 	}

// 	return courses, nil
// }
