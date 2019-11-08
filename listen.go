package main

import (
	"fmt"
	_ "github.com/lib/pq"
)
func main(){
	tryPickupJob := make(chan interface{})
	//equivalent to 'LISTEN ci_jobs_status_channel;'
	listener.Listen("ci_jobs_status_channel")
	go func() {
	for event := range listener.Notify {
		select {
		case tryPickupJob <- true:
		}
	}
	close(tryPickupJob)
	}
	
	for job := range tryPickupJob {
		fmt.Printf("here")
	}
}
