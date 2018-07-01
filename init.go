package main

import (
	"github.com/wheresalice/pronounsapi/lib"
	"log"
)

func main() {
	lib.ExecSQL("CREATE TABLE IF NOT EXISTS accounts (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), username STRING, UNIQUE(username))")
	lib.ExecSQL("CREATE TABLE IF NOT EXISTS pronouns (accountid UUID, pronouns STRING[], UNIQUE(accountid))")

	lib.NewUser("alice")
	lib.SetUserPronoun("alice", []string{"They/Them"})

	log.Printf("%s: %v", "alice", lib.GetUserPronouns("alice"))
}