//postgres pub sub
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
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
	for i := 0; i < 3; i++ {
		waitForNotification(listener, db)
		_, err := db.Query("UPDATE ci_jobs SET status='initializing' WHERE id = (SELECT id FROM ci_jobs WHERE status='new' ORDER BY id FOR UPDATE SKIP LOCKED LIMIT 1) RETURNING *; ")
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("close db connection")
	db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
} //main
func waitForNotification(l *pq.Listener, db *sql.DB) {
	select {
	case n := <-l.Notify:
		fmt.Println("received notification")
		res := strings.Split(n.Extra, ",")
		fmt.Printf("%s %s \n", res[1], res[3])

		// queryDB(db)
	case <-time.After(2 * time.Second):
		go l.Ping()
		// Check if there's more work available, just in case it takes
		// a while for the Listener to notice connection loss and
		// reconnect.
		fmt.Println("received no work for x seconds, checking for new work")
	}
}
func queryDB(db *sql.DB) {
	rows, _ := db.Query("SELECT * FROM ci_jobs")
	for rows.Next() {
		var row1 DBRow
		if err := rows.Scan(&row1.id, &row1.jobname, &row1.status, &row1.status_change_time); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s %s %s %s\n", row1.id, row1.jobname, row1.status, row1.status_change_time)
	}
}
