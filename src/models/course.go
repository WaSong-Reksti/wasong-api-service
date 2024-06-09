package models

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

/*
Course Model
*/
type Course struct {
	ID           string                 `firestore:"-"`
	Name         string                 `firestore:"Name,omitempty"`
	Description  string                 `firestore:"Description,omitempty"`
	Instrument   string                 `firestore:"Instrument,omitempty"`
	Instructor   *Instructor            `firestore:"-"`
	InstructorID *firestore.DocumentRef `firestore:"InstructorID,omitempty"`
	Assignments  []Assignment           `firestore:"-"`
}

func GetAllCourses(ctx context.Context, firestore *firestore.Client) ([]Course, error) {
	iter := firestore.Collection("Course").Documents(ctx)
	var courses []Course
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate: %v", err)
		}
		var course Course
		if err := doc.DataTo(&course); err != nil {
			return nil, fmt.Errorf("failed to convert data to Course struct: %v", err)
		}
		course.ID = doc.Ref.ID
		courses = append(courses, course)
	}
	return courses, nil
}

func GetCoursesById(ctx context.Context, firestore *firestore.Client, courseId string) (*Course, error) {
	docSnapshot, err := firestore.Collection("Course").Doc(courseId).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Course document %v", err)
	}
	var course Course
	if err := docSnapshot.DataTo(&course); err != nil {
		return nil, fmt.Errorf("failed to convert data to Course struct: %v", err)
	}
	course.ID = docSnapshot.Ref.ID

	// Fetch instructor data for the course
	instructorDocRef := course.InstructorID
	if instructorDocRef != nil {
		instructor, err := GetInstructor(ctx, instructorDocRef)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch instructor data for course %s: %v", course.ID, err)
		}
		course.Instructor = instructor

		// Fetch user data from instructor
		userDocRef := instructor.UserID
		if userDocRef != nil {
			user, err := getUser(ctx, userDocRef)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch user data for instructor of course %s: %v", course.ID, err)
			}
			instructor.User = user
		}
	}
	return &course, nil
}

func CreateCourse(ctx context.Context, client *firestore.Client, course *Course) (string, *Course, error) {
	docRef, _, err := client.Collection("Course").Add(ctx, course)
	if err != nil {
		return "", nil, fmt.Errorf("error: %v", err)
	}
	createdCourse := course
	createdCourse.ID = docRef.ID
	return docRef.ID, createdCourse, nil
}

func UpdateCourse(ctx context.Context, firestoreClient *firestore.Client, courseID string, keyValue map[string]interface{}) (map[string]interface{}, error) {
	updates := make([]firestore.Update, 0)
	for key, value := range keyValue {
		updates = append(updates, firestore.Update{
			Path:  key,
			Value: value,
		})
	}

	_, err := firestoreClient.Collection("Course").Doc(courseID).Update(ctx, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update course: %v", err)
	}

	docSnap, err := firestoreClient.Collection("Course").Doc(courseID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated course: %v", err)
	}
	return docSnap.Data(), nil
}

func DeleteCourseByID(ctx context.Context, firestoreClient *firestore.Client, courseID string) error {
	_, err := firestoreClient.Collection("Course").Doc(courseID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete: %v", err)
	}

	return nil

}

func AddStudentToCourse(ctx context.Context, firestoreClient *firestore.Client, user *User, courseRef *firestore.DocumentRef) error {
	// Retrieve the student document
	studentDoc, err := user.StudentRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get student document: %v", err)
	}

	// Get the current courses array
	courses, err := studentDoc.DataAt("Courses")
	if err != nil {
		return fmt.Errorf("failed to get Courses field: %v", err)
	}

	// Type assert the courses to an array of DocumentRefs
	// Type assert the courses to an array of interface{}
	coursesArray, ok := courses.([]interface{})
	if !ok {
		return fmt.Errorf("Courses field is not of the expected type: %T", courses)
	}

	// Convert each element to *firestore.DocumentRef
	var updatedCoursesArray []*firestore.DocumentRef
	for _, course := range coursesArray {
		docRef, ok := course.(*firestore.DocumentRef)
		if !ok {
			return fmt.Errorf("Courses field contains non-document reference elements")
		}
		updatedCoursesArray = append(updatedCoursesArray, docRef)
	}

	// Append the new course reference to the courses array
	coursesArray = append(coursesArray, courseRef)

	// Update the Courses field in the student document
	_, err = user.StudentRef.Update(ctx, []firestore.Update{
		{
			Path:  "Courses",
			Value: coursesArray,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update Courses field: %v", err)
	}

	return nil
}

func JoinCourseWithInstructor(ctx context.Context, firestore *firestore.Client) ([]Course, error) {
	// Step 1: Query Courses
	iter := firestore.Collection("Course").Documents(ctx)
	var courses []Course
	// Step 2: Fetch instructor data for each course
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate: %v", err)
		}
		var course Course
		if err := doc.DataTo(&course); err != nil {
			return nil, fmt.Errorf("failed to convert data to Course struct: %v", err)
		}
		course.ID = doc.Ref.ID

		// Fetch instructor data for the course
		instructorDocRef := course.InstructorID
		if instructorDocRef != nil {
			instructor, err := GetInstructor(ctx, instructorDocRef)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch instructor data for course %s: %v", course.ID, err)
			}
			course.Instructor = instructor

			// fetch user data from instructor
			userDocRef := instructor.UserID
			if userDocRef != nil {
				user, err := getUser(ctx, userDocRef)
				if err != nil {
					return nil, fmt.Errorf("failed to fetch user data for instructor of course %s: %v", course.ID, err)
				}
				instructor.User = user
			}
		}

		courses = append(courses, course)
	}

	return courses, nil
}
