package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

func handleJobs(w http.ResponseWriter, r *http.Request) {

	fmt.Println("here")
	//msg := fmt.Sprintf("jobs")
	database, _ := sql.Open("sqlite3", "./localsqllite.db")
	//message("test one 1 2 3")
	rows, _ := database.Query("SELECT id, targetservers, ixconfs, results FROM jobs")
	var id int
	var ts string
	var ixc string
	var res string
	for rows.Next() {
		rows.Scan(&id, &ts, &ixc, res)
		message(fmt.Sprintf(strconv.Itoa(id) + ": " + ts + " " + ixc + " " + res))
		//	spew.Dump(err)
	}

}
