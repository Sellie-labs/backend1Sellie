package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var dbPool *sql.DB

// InitPool initializes the database connection pool and runs migrations.
func InitPool() error {
	var err error

	// Construct the database connection string
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Initialize the connection pool
	dbPool, err = sql.Open("pgx", connectionString)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	// Verify the connection
	if err = dbPool.Ping(); err != nil {
		return fmt.Errorf("error pinging database: %v", err)
	}

	// Run migrations
	// if err = runMigrations(dbPool, migrationsPath); err != nil {
	// 	return fmt.Errorf("error running migrations: %v", err)
	// }

	return nil
}

func runMigrations(db *sql.DB, migrationsPath string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

// GetPool returns the database connection pool.
func GetPool() *sql.DB {
	return dbPool
}

// ClosePool closes the database connection pool.
func ClosePool() error {
	if dbPool != nil {
		return dbPool.Close()
	}
	return nil
}
