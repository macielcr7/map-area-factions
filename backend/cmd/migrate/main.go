package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var (
	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string
	dbSSLMode  string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Database migration tool for Map Area Factions",
		Long:  "A CLI tool to manage database migrations for the Map Area Factions backend.",
	}

	rootCmd.PersistentFlags().StringVar(&dbHost, "host", "localhost", "Database host")
	rootCmd.PersistentFlags().StringVar(&dbPort, "port", "5432", "Database port")
	rootCmd.PersistentFlags().StringVar(&dbUser, "user", "postgres", "Database user")
	rootCmd.PersistentFlags().StringVar(&dbPassword, "password", "postgres", "Database password")
	rootCmd.PersistentFlags().StringVar(&dbName, "dbname", "map_factions", "Database name")
	rootCmd.PersistentFlags().StringVar(&dbSSLMode, "sslmode", "disable", "SSL mode")

	// Override with environment variables
	if host := os.Getenv("DB_HOST"); host != "" {
		dbHost = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		dbPort = port
	}
	if user := os.Getenv("DB_USER"); user != "" {
		dbUser = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		dbPassword = password
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		dbName = name
	}
	if sslMode := os.Getenv("DB_SSL_MODE"); sslMode != "" {
		dbSSLMode = sslMode
	}

	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
	rootCmd.AddCommand(statusCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all pending migrations",
	Run: func(cmd *cobra.Command, args []string) {
		runMigrations("up")
	},
}

var downCmd = &cobra.Command{
	Use:   "down [count]",
	Short: "Rollback migrations",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		count := 1
		if len(args) > 0 {
			fmt.Sscanf(args[0], "%d", &count)
		}
		rollbackMigrations(count)
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	Run: func(cmd *cobra.Command, args []string) {
		showStatus()
	},
}

func getDB() *sql.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	return db
}

func runMigrations(direction string) {
	db := getDB()
	defer db.Close()

	// Create migrations table if it doesn't exist
	createMigrationsTable(db)

	// Get list of migration files
	migrationFiles, err := getMigrationFiles(direction)
	if err != nil {
		log.Fatal("Failed to get migration files:", err)
	}

	// Get applied migrations
	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		log.Fatal("Failed to get applied migrations:", err)
	}

	// Filter migrations
	var migrationsToRun []string
	for _, file := range migrationFiles {
		if direction == "up" {
			if !contains(appliedMigrations, file) {
				migrationsToRun = append(migrationsToRun, file)
			}
		} else {
			if contains(appliedMigrations, file) {
				migrationsToRun = append(migrationsToRun, file)
			}
		}
	}

	if len(migrationsToRun) == 0 {
		fmt.Println("No migrations to run")
		return
	}

	// Run migrations
	for _, migration := range migrationsToRun {
		fmt.Printf("Running migration: %s\n", migration)

		content, err := ioutil.ReadFile(filepath.Join("migrations", migration))
		if err != nil {
			log.Fatal("Failed to read migration file:", err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			log.Fatal("Failed to execute migration:", err)
		}

		// Record migration
		if direction == "up" {
			_, err = db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", migration)
		} else {
			_, err = db.Exec("DELETE FROM schema_migrations WHERE version = $1", migration)
		}

		if err != nil {
			log.Fatal("Failed to record migration:", err)
		}

		fmt.Printf("✅ Migration %s completed\n", migration)
	}

	fmt.Println("All migrations completed successfully!")
}

func rollbackMigrations(count int) {
	fmt.Printf("Rolling back %d migrations...\n", count)
	// Implementation for rollback
}

func showStatus() {
	db := getDB()
	defer db.Close()

	// Create migrations table if it doesn't exist
	createMigrationsTable(db)

	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		log.Fatal("Failed to get applied migrations:", err)
	}

	migrationFiles, err := getMigrationFiles("up")
	if err != nil {
		log.Fatal("Failed to get migration files:", err)
	}

	fmt.Println("Migration Status:")
	fmt.Println("================")

	for _, file := range migrationFiles {
		status := "❌ Pending"
		if contains(appliedMigrations, file) {
			status = "✅ Applied"
		}
		fmt.Printf("%-30s %s\n", file, status)
	}
}

func createMigrationsTable(db *sql.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT NOW()
		)
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create migrations table:", err)
	}
}

func getMigrationFiles(direction string) ([]string, error) {
	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		return nil, err
	}

	var migrationFiles []string
	for _, file := range files {
		filename := filepath.Base(file)
		if strings.Contains(filename, direction) {
			migrationFiles = append(migrationFiles, filename)
		}
	}

	sort.Strings(migrationFiles)
	return migrationFiles, nil
}

func getAppliedMigrations(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []string
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		migrations = append(migrations, version)
	}

	return migrations, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
