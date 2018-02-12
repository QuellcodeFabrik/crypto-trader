package main

import (
    "encoding/json"
    "github.com/gorilla/mux"
    "log"
    "net/http"
)

type Person struct {
    ID        string   `json:"id,omitempty"`
    Firstname string   `json:"firstname,omitempty"`
    Lastname  string   `json:"lastname,omitempty"`
    Address   *Address `json:"address,omitempty"`
}
type Address struct {
    City  string `json:"city,omitempty"`
    State string `json:"state,omitempty"`
}

var people[]Person

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.RequestURI)
        // Call the next handler, which can be another middleware in the chain, or the final handler.
        next.ServeHTTP(w, r)
    })
}

func GetPeople(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(people)
}
func GetPerson(w http.ResponseWriter, r *http.Request) {}
func CreatePerson(w http.ResponseWriter, r *http.Request) {}
func DeletePerson(w http.ResponseWriter, r *http.Request) {}

// our main function
func main() {
    log.Print("Go server starting...")

    people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
    people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
    people = append(people, Person{ID: "3", Firstname: "Francis", Lastname: "Sunday"})

    router := mux.NewRouter()
    router.HandleFunc("/people", GetPeople).Methods("GET")
    router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
    router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
    router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")

    router.Use(loggingMiddleware)

    log.Fatal(http.ListenAndServe(":8205", router))
}
