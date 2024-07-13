package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func ping(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ping", ping).Methods("GET")
	http.ListenAndServe(":8084", r)
}
