//postgres pub sub
package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "test"
)

type DBRow struct {
	id                 string
	jobname            string
	status             string
	status_change_time string
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, _ := db.Query("SELECT * FROM ci_jobs")
	for rows.Next() {
		var row1 DBRow
		if err := rows.Scan(&row1.id, &row1.jobname, &row1.status, &row1.status_change_time); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s %s %s %s\n", row1.id, row1.jobname, row1.status, row1.status_change_time)
	}

	// tryPickupJob := make(chan interface{})
	//equivalent to 'LISTEN ci_jobs_status_channel;'
	minReconn := 10 * time.Second
	maxReconn := time.Minute
	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	listener := pq.NewListener(psqlInfo, minReconn, maxReconn, reportProblem)
	err = listener.Listen("ci_jobs_status_channel")
	if err != nil {
		panic(err)
	}
	for {
		waitForNotification(listener)
		// _, err := db.Query("UPDATE ci_jobs SET status = 'new' ")
		// if err != nil {
		// 	panic(err)
		// }
	}
} //main
func waitForNotification(l *pq.Listener) {
	select {
	case <-l.Notify:
		fmt.Println("received notification, new work available")
	case <-time.After(5 * time.Second):
		go l.Ping()
		// Check if there's more work available, just in case it takes
		// a while for the Listener to notice connection loss and
		// reconnect.
		fmt.Println("received no work for 90 seconds, checking for new work")
	}
}
