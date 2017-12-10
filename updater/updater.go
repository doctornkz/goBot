package updater

import (
	"database/sql"
	"fmt"
	"log"
	// INI parser
	"github.com/asjustas/goini"
	// SQLite3 package
	_ "github.com/mattn/go-sqlite3"
)

// Update message's counter
func Update(ID int, UserName string, FirstName string, LastName string) {
	// Open DB
	//sqlite3.ErrCantOpen
	conf, err := goini.Load("./settings.ini")
	if err != nil {
		panic(err)
	}

	sqliteDB := conf.Str("main", "SQLITE_DB")

	db, err := sql.Open("sqlite3", sqliteDB) // TODO Dublicated INI parser
	if err != nil {
		log.Fatal(err)
		// TODO: Input "Create DB?"
	}
	defer db.Close()

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
		UserUpdate(ID, UserName, FirstName, LastName)
	}

	log.Printf("Updater: UserID %d found with %d messages", ID, count)

	// Insert section

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
		//sqlQueryName = "insert into user (userid, username, firstname, lastname) values (?, ?, ?, ?)"
		//log.Printf("Updater: SQL Insert UserName %s", sqlQueryName)

	} else {

		sqlQuery = "update num_messages set count =? where userid = ?"
		log.Printf("Updater: SQL Update %s", sqlQuery)
		//update userinfo set username=? where uid=?
	}

	smth, err := tx.Prepare(sqlQuery)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer smth.Close()
	// Exec section
	log.Printf("Updater: SQL Insert %s", sqlQuery)
	count++
	log.Printf("Updater: UserID: %d count: %d", ID, count)
	if operation == "insert" {
		_, err = smth.Exec(ID, count)
		if err != nil {
			log.Println(err)
			log.Fatal(err)
		}

	} else {
		_, err = smth.Exec(count, ID) // WTF?
		if err != nil {
			log.Println(err)
			log.Fatal(err)
		}
	}

	// Commit section
	log.Println("Updater: Pre commit step")
	tx.Commit()
	log.Println("Updater: Committed")

	// Username database updating

}

// UserUpdate  Username updating
func UserUpdate(ID int, UserName string, FirstName string, LastName string) {
	// Open DB
	//sqlite3.ErrCantOpen
	db, err := sql.Open("sqlite3", "./exyandex.db")
	if err != nil {
		log.Fatal(err)
		// TODO: Input "Create DB?"
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	sqlQueryUser := "insert into user (userid, username, firstname, lastname) values (?, ?, ?, ?)"
	log.Printf("Updater: SQL Insert UserName %s", sqlQueryUser)
	//sqlQueryName = "insert into user (userid, username, firstname, lastname) values (?, ?, ?, ?)"
	//log.Printf("Updater: SQL Insert UserName %s", sqlQueryName)
	smth, err := tx.Prepare(sqlQueryUser)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer smth.Close()
	_, err = smth.Exec(ID, UserName, FirstName, LastName)
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}
	tx.Commit()

}
