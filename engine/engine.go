package engine

import (
	"database/sql"
	"log"
	"strconv"
)

// Status  - TOP20 in chat
func Status(db *sql.DB, ID int) string {

	rows, err := db.Query("select user.username, user.firstname, num_messages.userid, num_messages.count from user inner join num_messages on user.userid=num_messages.userid order by count desc")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	output := " -= TOP LIST =- \r\n"
	index := 1
	limit := 20
	for rows.Next() { // Generating TOP 20
		var username string
		var firstname string
		var count string
		var userid int
		err = rows.Scan(&username, &firstname, &userid, &count)
		if username == "" {
			username = firstname
		}

		if err != nil {
			log.Fatal(err)
		}
		log.Println(strconv.Itoa(index), username, count)

		if index <= limit {
			output = output + strconv.Itoa(index) + ". " + username + " = " + count + "\r\n"
		}
		if (ID == userid) && (index > limit) {
			output = output + "...\r\n" + strconv.Itoa(index) + ". " + username + " = " + count + "\r\n"
		}
		index++

	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return output
}
