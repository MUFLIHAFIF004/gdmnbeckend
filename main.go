package main

import (
	"fmt"
	"net/http"
	"os"

	"gudangmng/config"
	"gudangmng/routes"
)

func main() {
	config.ConnectDatabase()

	http.HandleFunc("/", routes.URL)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server Gudang Running on Port: " + port)
	http.ListenAndServe(":"+port, nil)
}