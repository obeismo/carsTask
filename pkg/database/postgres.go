package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/sirupsen/logrus"

	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func InitDB() (*sql.DB, error) {
	cfg := Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("SSL_MODE"),
	}
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logrus.Fatalf("failed to open db connection: %s", err.Error())
		return db, err
	}

	err = db.Ping()
	if err != nil {
		logrus.Fatalf("failed to ping db: %s", err.Error())
		return db, err
	}

	logrus.Info("Database connection established")

	driverURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	err = RunMigrations(driverURL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigrations(driverURL string) error {
	driver, err := migrate.New(fmt.Sprintf("file://%s/em_test/schema", os.Getenv("PROGRAM_DIRECTORY_PATH")), driverURL)
	// driver, err := migrate.New("file:///home/maxim/code/em_test/schema", driverURL)
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	err = driver.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	logrus.Info("Migrations applied successfully")
	return nil
}
