package main

import (
	"backend/internal/config"
	"backend/internal/repository/sqlite"
	"backend/internal/repository/sqlite/migrations"
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func runMigrateCmd(conf config.DatabaseConfig, args []string) {
	cmd := ""
	if len(args) > 0 {
		cmd = args[0]
	}

	switch cmd {
	case "up":
		runMigrateUp(conf)
	case "down":
		runMigrateDown(conf, args[1:])
	case "status":
		runMigrateStatus(conf)
	default:
		log.Fatalf("unknown migrate command: %s (expected 'up', 'down', or 'status')", cmd)
	}
}

func runMigrateUp(conf config.DatabaseConfig) {
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

func runMigrateDown(conf config.DatabaseConfig, args []string) {
	fs := flag.NewFlagSet("migrate down", flag.ExitOnError)
	to := fs.Int("to", -1, "target version to roll back to (required)")
	err := fs.Parse(args)
	if err != nil {
		log.Fatalf("error parsing migrate down flags: %v", err)
	}

	if *to < 0 {
		log.Fatal("migrate down: -to flag is required (e.g. migrate down -to 2)")
	}

	ctx := context.Background()

	db, err := sqlite.OpenDB(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	reader := bufio.NewReader(os.Stdin)
	confirm := func(m migrations.Migration) bool {
		fmt.Printf("Roll back %03d_%s? [y/N] ", m.Version, m.Name)
		answer, _ := reader.ReadString('\n')
		return strings.TrimSpace(strings.ToLower(answer)) == "y"
	}

	if err := migrations.Rollback(ctx, db, *to, confirm); err != nil {
		log.Fatalf("rollback failed: %v", err)
	}
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
