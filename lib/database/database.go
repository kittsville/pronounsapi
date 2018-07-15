package database

import (
	"github.com/coreos/bbolt"
	"fmt"
	"encoding/json"
)

var (
	DBCon *bolt.DB
	Err error
)

type User struct {
	Username string
	Pronouns []string
}

func CreateUserBucket(username string) error {
	DBCon.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte("DB"))
		if err != nil {
			return fmt.Errorf("could not open db, %v", err)
		}
		_, err = root.CreateBucketIfNotExists([]byte(username))
		if err != nil {
			return fmt.Errorf("could not create %s's bucket: %v", username, err)
		}
		return nil
	})
	return nil
}

func SetUserData(user User) error {
	CreateUserBucket(user.Username)

	userBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("could not marshal user data: %v", err)
	}
	DBCon.Update(func(tx *bolt.Tx) error {
		err = tx.Bucket([]byte("DB")).Bucket([]byte(user.Username)).Put([]byte("pronouns"), userBytes)
		if err != nil {
			return fmt.Errorf("could not set user data: %v", err)
		}
		return nil
	})
	return nil
}

func GetUserData(username string) User {
	var user User
	DBCon.View(func(tx *bolt.Tx) error {
		userBytes := tx.Bucket([]byte("DB")).Bucket([]byte(username)).Get([]byte("pronouns"))
		err := json.Unmarshal(userBytes, &user)
		if err != nil {
			return fmt.Errorf("could not get user data: %v", err)
		}
		return nil
	})
	return user
}