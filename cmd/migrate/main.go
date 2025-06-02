package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/omkar273/codegeeky/ent"
	"github.com/omkar273/codegeeky/internal/config"
	"github.com/omkar273/codegeeky/internal/logger"
)

func main() {
	// Parse command line flags
	dryRun := flag.Bool("dry-run", false, "Print migration SQL without executing it")
	flag.Parse()

	// Load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger, err := logger.NewLogger(cfg)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	// Get DSN from config
	dsn := cfg.Postgres.GetDSN()
	logger.Infow("Connecting to database", "host", cfg.Postgres.Host)

	// Create Ent client
	client, err := ent.Open("postgres", dsn)
	if err != nil {
		logger.Fatalw("Failed to connect to postgres", "error", err)
	}
	//nolint:errcheck
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Run auto migration
	logger.Info("Running database migrations...")

	// Check if we're in dry-run mode
	if *dryRun {
		logger.Info("Dry run mode - printing migration SQL without executing")
		// In dry-run mode, we just print the SQL that would be executed
		err = client.Schema.WriteTo(ctx, os.Stdout)
		if err != nil {
			logger.Fatalw("Failed to generate migration SQL", "error", err)
		}
	} else {
		// Run the actual migration
		err = client.Schema.Create(ctx)
		if err != nil {
			logger.Fatalw("Failed to create schema resources", "error", err)
		}
		logger.Info("Migration completed successfully")
	}

	fmt.Println("Migration process completed")
}
