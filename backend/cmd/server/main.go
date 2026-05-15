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

	cmd := ""
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}
	switch cmd {
	case "generate-vapid":
		runGenerateVAPIDCmd()
	case "migrate":
		runMigrateCmd(conf.Database, os.Args[2:])
	case "", "serve":
		runServer(conf)
	default:
		log.Fatalf("unknown command: %s", cmd)
	}
}
