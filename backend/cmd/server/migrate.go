package main

import (
	"backend/internal/config"
	"backend/internal/repository/sqlite"
	"backend/internal/repository/sqlite/migrations"
	"context"
	"flag"
	"log"
)

func runMigrateCmd(conf config.DatabaseConfig, args []string) {
	fs := flag.NewFlagSet("migrate", flag.ExitOnError)
	status := fs.Bool("status", false, "show migration status")
	err := fs.Parse(args)
	if err != nil {
		log.Fatalf("error parsing migrate flags: %v", err)
	}

	if *status {
		runMigrateStatus(conf)
		return
	}
	runMigrate(conf)
}

func runMigrateStatus(conf config.DatabaseConfig) {
	ctx := context.Background()

	db, err := sqlite.OpenDB(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	statuses, err := migrations.Status(ctx, db)
	if err != nil {
		log.Fatal(err)
	}

	if len(statuses) == 0 {
		log.Println("no migrations found")
		return
	}

	for _, s := range statuses {
		if s.Applied {
			log.Printf("  applied: %03d_%s (at %s)", s.Version, s.Name, s.AppliedAt.Format("2006-01-02 15:04:05"))
		} else {
			log.Printf("  pending: %03d_%s", s.Version, s.Name)
		}
	}
}

func runMigrate(conf config.DatabaseConfig) {
	ctx := context.Background()

	db, err := sqlite.OpenDB(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := migrations.Apply(ctx, db); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
}
