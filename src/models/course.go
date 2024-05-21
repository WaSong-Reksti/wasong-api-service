package models

/*
Course Model
*/
type Course struct {
	Name        string     `firestore: "Name, omitempty"`
	Description string     `firestore: "Description, omitempty"`
	Instrument  string     `firestore: "Instrument, omitempty"`
	Instructor  Instructor `firestore: "Instructor, omitempty"`
}

func GetAllCourses() {

}
