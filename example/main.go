package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	g "github.com/musaubrian/gno"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!\n")
	g.Log(g.INFO, fmt.Sprintf("endpoint '%s' hit", r.RequestURI))
}

func now(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n", time.Now().Format(time.RFC822))
	g.Log(g.INFO, fmt.Sprintf("endpoint '%s' hit", r.RequestURI))
}

func main() {
	port := ":8080"
	http.HandleFunc("/", home)
	http.HandleFunc("/now", now)
	g.Log(g.INFO, "Server running at: "+port)
	log.Fatal(http.ListenAndServe(port, nil))
}
