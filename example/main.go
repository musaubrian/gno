package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!\n")
	log.Printf("endpoint '%s' hit", r.RequestURI)
}
func now(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n", time.Now().Format(time.RFC822))
	log.Printf("endpoint '%s' hit", r.RequestURI)
}

func main() {
	port := ":8080"
	http.HandleFunc("/", greet)
	http.HandleFunc("/now", now)
	log.Println("Server running at: ", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
