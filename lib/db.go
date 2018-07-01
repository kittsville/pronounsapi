package lib

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

const dsn = "postgresql://pronounsapi@localhost:26257/pronounsapi?sslmode=disable"

func ExecSQL(sqlStatement string) {
	db, err := sql.Open("postgres", dsn)
	checkErr(err)
	_, err = db.Exec(sqlStatement)
	checkErr(err)
}

func NewUser(userName string) {
	sqlStatement := fmt.Sprintf("INSERT INTO accounts(username) VALUES ('%s')", userName)
	ExecSQL(sqlStatement)
}

func UserNameToID(userName string) string {
	db, err := sql.Open("postgres", dsn)
	checkErr(err)
	sqlStatement := `SELECT id FROM accounts WHERE username=$1`
	var id string
	row := db.QueryRow(sqlStatement, userName)
	switch err := row.Scan(&id); err {
	case sql.ErrNoRows:
		log.Print("No rows were returned!")
	case nil:
		return id
	default:
		panic(err)
	}
	return id
}

func GetUserPronouns(userName string) string {
	db, err := sql.Open("postgres", dsn)
	checkErr(err)
	sqlStatement := `SELECT pronouns FROM pronouns INNER JOIN accounts ON accountid=id WHERE username=$1`
	var pronouns string
	row := db.QueryRow(sqlStatement, userName)
	switch err := row.Scan(&pronouns); err {
	case sql.ErrNoRows:
		return pronouns
	case nil:
		return pronouns
	default:
		panic(err)
	}
	return pronouns
}

func SetUserPronoun(userName string, pronouns []string) {
	id := UserNameToID(userName)
	sqlStatement := fmt.Sprintf("INSERT INTO pronouns(accountid, pronouns)"+
		" VALUES('%s', ARRAY['%s'])" +
		" ON CONFLICT(accountid)" +
		" DO UPDATE SET pronouns = excluded.pronouns",
		id, strings.Join(pronouns, "', '"))
	ExecSQL(sqlStatement)
}

func checkErr(err error) {
	if err != nil {
		log.Print(err)
	}
}
