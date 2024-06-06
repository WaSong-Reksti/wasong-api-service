package main

import (
	// "errors"
	"context"
	"example/wasong-api-service/src/database"
	"example/wasong-api-service/src/routes"
	"fmt"
	"log"
	"net"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	ctx := context.Background()
	firestoreClient, err := database.InitializeFirestoreClient(&ctx)
	if err != nil {
		log.Fatalf("error initializing Firestore client: %v", err)
		return
	}
	fmt.Println("Successfuly connected to firestore client")
	defer firestoreClient.Close()

	routes.InitializeUserRoutes(ctx, router, firestoreClient)
	routes.InitializeCourseRoutes(ctx, router, firestoreClient)
	fmt.Println("Initialize courses route")

	router.Run("localhost:8080")
	fmt.Println("Sucessfully run on localhost:8080")

}

func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no IP address found")
}
