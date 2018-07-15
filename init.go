package main

import (
	"github.com/coreos/bbolt"
	"log"
	"github.com/wheresalice/pronounsapi/lib/database"
)

func main() {
	database.DBCon, database.Err = bolt.Open("data/bolt.db", 0600, nil)
	if database.Err != nil {
		log.Fatalf("could not open db, %v", database.Err)
	}

	err := database.CreateUserBucket("alice")
	if err != nil {
		log.Fatalf("could create user bucket, %v", err)
	}

	alice := database.User{Username: "alice", Pronouns: []string{"They/them"}}
	err = database.SetUserData(alice)
	if err != nil {
		log.Fatalf("could set user data, %v", err)
	}

	database.DBCon.View(func(tx *bolt.Tx) error {
		err = tx.CopyFile("data/backup.db", 0600)
		if err != nil {
			log.Fatalf("could not back up database, %v", err)
		}
		return nil
	})

	log.Print("DB Setup Done")
}