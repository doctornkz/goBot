package updater

import (
	"database/sql"
	"fmt"
	"log"

	// SQLite3 package
	_ "github.com/mattn/go-sqlite3"
)

// Update message's counter
func Update(db *sql.DB, ID int, UserName string, FirstName string, LastName string, Date int, Text string) {
	// Select rows with ID
	sqlSelectQuery := "select count from num_messages where userid= ?"
	log.Printf("Updater: SQL Select %s", sqlSelectQuery)
	stmt, err := db.Prepare(sqlSelectQuery)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Query section
	//var userid int
	var count int
	operation := "update" // TODO Make simplier, or goto?
	err = stmt.QueryRow(ID).Scan(&count)

	if err != nil {
		// New user detected
		fmt.Println(err)
		//userid = ID
		count = 0

		operation = "insert"
		userUpdate(db, ID, UserName, FirstName, LastName)
	}

	// Update messages
	log.Printf("Updater: Place new message in DB")

	messagesUpdate(db, ID, Date, Text)

	log.Printf("Updater: UserID %d found with %d messages", ID, count)

	// Insert section
	log.Printf("Updater: Operation flag is %s", operation)
	// Begin section
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	var sqlQuery string
	// Prepare section

	if operation == "insert" {

		sqlQuery = "insert into num_messages (userid, count) values (?, ?)"
		log.Printf("Updater: SQL Insert %s", sqlQuery)

	} else {

		sqlQuery = "update num_messages set count =? where userid = ?"
		log.Printf("Updater: SQL Update %s", sqlQuery)
	}

	smth, err := tx.Prepare(sqlQuery)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer smth.Close()
	// Exec section

	count++
	log.Printf("Updater: UserID: %d count: %d", ID, count)
	if operation == "insert" {
		_, err = smth.Exec(ID, count)
		if err != nil {
			log.Println(err)
			log.Fatal(err)
		}

	} else {
		_, err = smth.Exec(count, ID)
		if err != nil {
			log.Println(err)
			log.Fatal(err)
		}
	}

	// Commit section
	log.Println("Updater: Pre commit step")
	tx.Commit()
	log.Println("Updater: Committed")
}

// MessagesUpdate - updater messages in chat
func messagesUpdate(db *sql.DB, ID int, Date int, Text string) {
	log.Println("Updater: Messages insert section")
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	sqlMessQuery := "insert into messages (userid, date, text) values (?, ?, ?)"
	log.Printf("Updater: SQL Insert %s", sqlMessQuery)

	updateMessageState, err := tx.Prepare(sqlMessQuery)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer updateMessageState.Close()

	_, err = updateMessageState.Exec(ID, Date, Text)
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}

	tx.Commit()
	log.Println("Updater: Message committed")

}

// UserUpdate  Username updating
func userUpdate(db *sql.DB, ID int, UserName string, FirstName string, LastName string) {

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	sqlQueryUser := "insert into user (userid, username, firstname, lastname) values (?, ?, ?, ?)"
	log.Printf("Updater: SQL Insert UserName %s", sqlQueryUser)

	updateUserState, err := tx.Prepare(sqlQueryUser)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer updateUserState.Close()

	_, err = updateUserState.Exec(ID, UserName, FirstName, LastName)
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}
	tx.Commit()
}
