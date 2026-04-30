package main

import (
	"backend/internal/config"
	"backend/internal/logging"
	"log"
	"os"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	logging.Init(conf.Logging)

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		if len(os.Args) > 2 && os.Args[2] == "status" {
			runMigrateStatus(conf.Database)
			return
		}
		runMigrate(conf.Database)
		return
	}

	runServer(conf)
}
