package main

import (
	"devconnect/config"
	"devconnect/routes"
	"fmt"
	"log"
)

func main() {
	config.ConnectDB()
	r := routes.SetupRouter()
	fmt.Println("Server running on port 8080")
	log.Fatal(r.Run(":8080"))
}
