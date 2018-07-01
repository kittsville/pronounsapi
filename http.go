package main

import (
"encoding/json"
"fmt"
"github.com/gorilla/mux"
"log"
"net/http"
	"github.com/wheresalice/pronounsapi/lib"
	"strings"
)

func main() {
	var router = mux.NewRouter()
	router.HandleFunc("/healthcheck", healthCheck).Methods("GET")
	router.HandleFunc("/message", handleQryMessage).Methods("GET")
	router.HandleFunc("/m/{msg}", handleUrlMessage).Methods("GET")
	router.HandleFunc("/u/{userName}", handleUser).Methods("GET")

	fmt.Println("Running server!")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userName := vars["userName"]
	pronouns := lib.GetUserPronouns(userName)

	pronouns = strings.Replace(pronouns, "{", "[", -1)
	pronouns = strings.Replace(pronouns, "}", "]", -1)
	fmt.Fprintf(w,`{"pronouns": %s}`, pronouns)
}

func handleQryMessage(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	message := vars.Get("msg")

	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func handleUrlMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	message := vars["msg"]

	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Still alive!")
}
