package main

import (
	"api/db"
	"api/service"
	"log"
)

var Version = "0.0.3"
var BuildTime string

func main() {
	var err error

	err = db.Connect()
	if err != nil {
		// TODO change to fatal
		log.Printf("Could not connect to Postgres: %s", err)
	}

	//err = db.RunMigrations()
	//if err != nil {
	//	log.Fatalf("Migration error: %s", err)
	//}

	service.RunHttpServer(Version, BuildTime)
}
