package main

import (
	// "errors"
	// "github.com/gin-gonic/gin"
	// "net/http"
	"fmt"
)

type student struct {
	ID		string
	Name	string
	About	string

}

/*
user type 
*/
type user struct {
	Username 	string
	Password 	string
	Email 		string
	Type 		string
}

func main() {
	var myUser user = user{"hakuna_matata", "123456", "hakuna@example.com", "student"}
	fmt.Println(myUser);
}	
