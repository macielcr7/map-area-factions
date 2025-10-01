package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
		Use:   "seed",
		Short: "Database seeding tool for Map Area Factions",
		Long:  "A CLI tool to populate the database with seed data for the Map Area Factions backend.",
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

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(listCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run all seed files",
	Run: func(cmd *cobra.Command, args []string) {
		runSeeds()
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available seed files",
	Run: func(cmd *cobra.Command, args []string) {
		listSeeds()
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

func runSeeds() {
	db := getDB()
	defer db.Close()

	// Get list of seed files
	seedFiles, err := getSeedFiles()
	if err != nil {
		log.Fatal("Failed to get seed files:", err)
	}

	if len(seedFiles) == 0 {
		fmt.Println("No seed files found")
		return
	}

	fmt.Printf("Found %d seed files\n", len(seedFiles))

	// Run seeds
	for _, seedFile := range seedFiles {
		fmt.Printf("Running seed: %s\n", seedFile)

		content, err := ioutil.ReadFile(filepath.Join("seeds", seedFile))
		if err != nil {
			log.Fatal("Failed to read seed file:", err)
		}

		// Split content by semicolon and execute each statement
		statements := strings.Split(string(content), ";")
		for _, statement := range statements {
			statement = strings.TrimSpace(statement)
			if statement == "" {
				continue
			}

			_, err = db.Exec(statement)
			if err != nil {
				log.Printf("Warning: Failed to execute statement: %v", err)
				// Continue with other statements
			}
		}

		fmt.Printf("✅ Seed %s completed\n", seedFile)
	}

	fmt.Println("All seeds completed successfully!")
}

func listSeeds() {
	seedFiles, err := getSeedFiles()
	if err != nil {
		log.Fatal("Failed to get seed files:", err)
	}

	if len(seedFiles) == 0 {
		fmt.Println("No seed files found")
		return
	}

	fmt.Println("Available seed files:")
	fmt.Println("====================")
	for _, file := range seedFiles {
		fmt.Printf("- %s\n", file)
	}
}

func getSeedFiles() ([]string, error) {
	files, err := filepath.Glob("seeds/*.sql")
	if err != nil {
		return nil, err
	}

	var seedFiles []string
	for _, file := range files {
		filename := filepath.Base(file)
		seedFiles = append(seedFiles, filename)
	}

	return seedFiles, nil
}
