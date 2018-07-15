package lib

import (
	"strings"
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"encoding/json"
	"github.com/wheresalice/pronounsapi/lib/database"
)

func HandleUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userName := vars["userName"]
	pronouns := database.GetUserData(userName)

	userString, _ := json.Marshal(pronouns)
	fmt.Fprintf(w, string(userString))
}

func webfingerResourceToUser(webfingerUser string) (string, string) {
	account := strings.TrimPrefix(webfingerUser, "acct:")
	components := strings.Split(account, "@")
	return components[0], components[1]
}

func HandleWebfinger(w http.ResponseWriter, r *http.Request) {
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

func HandleActor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userName := vars["userName"]
	userData := database.GetUserData(userName)

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
`, r.Host, userName, userName, r.Host, userData.Pronouns)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Still alive!")
}
