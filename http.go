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

	router.HandleFunc("/u/{userName}", handleUser).Methods("GET")
	router.HandleFunc("/.well-known/webfinger", handleWebfinger).Methods("GET")
	router.HandleFunc("/a/{userName}", handleActor).Methods("GET")

	fmt.Println("Running server!")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userName := vars["userName"]
	pronouns := lib.GetUserPronouns(userName)

	pronouns = strings.Replace(pronouns, "{", "[", -1)
	pronouns = strings.Replace(pronouns, "}", "]", -1)
	fmt.Fprintf(w,`
{
"username": "%s",
"pronouns": %s
}`, userName, pronouns)
}

func webfingerResourceToUser(webfingerUser string) (string, string) {
	account := strings.TrimPrefix(webfingerUser, "acct:")
	components := strings.Split(account, "@")
	return components[0], components[1]
}

func handleWebfinger(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	user, _ := webfingerResourceToUser(vars.Get("resource"))
	fmt.Fprintf(w,
		`
{
  "subject": "acct:%s@%s",
  "links": [
    {
      "rel": "self",
      "type": "application/activity+json",
      "href": "http://%s/a/%s"
    }
  ]
}
`, user, r.Host, r.Host, user)
}

func handleActor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userName := vars["userName"]

	fmt.Fprintf(w, `
{
  "@context": [
    "https://www.w3.org/ns/activitystreams",
    "https://w3id.org/security/v1"
  ],
  "id": "http://%s/a/%s",
  "type": "Person",
  "preferredUsername": "%s",
  "inbox": "http://%s/inbox",
  "summary": "Pronouns: %s"
}
`, r.Host, userName, userName, r.Host, strings.Replace(lib.GetUserPronouns(userName), "\"", "'", -1))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Still alive!")
}
