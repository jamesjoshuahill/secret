package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	log.Fatalln(err)
}
