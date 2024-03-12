package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users = []User{
	{ID: 1, Name: "AnuchitO", Age: 18},
}

func usersHandler(write http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		log.Println("GET")
		b, err := json.Marshal(users)
		if err != nil {
			write.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(write, "error: %v", err)
			return
		}
		write.Header().Set("Content-Type", "application/json")
		write.Write(b)
	}
}

func healthHandler(write http.ResponseWriter, request *http.Request) {
	write.WriteHeader(http.StatusOK)
	write.Write([]byte("OK"))
}

type Logger struct {
	Handler http.Handler
}

func (l Logger) ServeHTTP(write http.ResponseWriter, request *http.Request) {
	start := time.Now()
	l.Handler.ServeHTTP(write, request)
	log.Printf("Server http middleware: %s %s %s %s",
		request.RemoteAddr,
		request.Method,
		request.URL,
		time.Since(start))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/users", usersHandler)
	mux.HandleFunc("/health", healthHandler)

	logMux := Logger{mux}

	srv := http.Server{
		Addr:    ":2565",
		Handler: logMux,
	}

	log.Println("Server started at :2565")
	log.Fatal(srv.ListenAndServe())
	log.Println("bye bye!")
}
