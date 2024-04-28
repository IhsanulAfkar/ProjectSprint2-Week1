package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func homePage(w http.ResponseWriter, r *http.Request){
	envValue := os.Getenv("TEST")
    fmt.Fprintf(w, envValue)
    fmt.Println("Endpoint Hit: homePage")
}
func handleRequests() {
    http.HandleFunc("/", homePage)
    log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
    handleRequests()
}