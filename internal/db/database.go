package db

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type PostgresStorage struct {
	Database *sql.DB
}

func NewPostgresStorage() *PostgresStorage {
	connectionStr := fmt.Sprintf(
		`host=%s port=%s user=%s dbname=%s password=%s sslmode=require`,
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USERNAME"),
		os.Getenv("POSTGRES_DATABASE"),
		os.Getenv("POSTGRES_PASSWORD"),
	)

	db, err := sql.Open("postgres", connectionStr)

	if err != nil {
		log.Fatalf("Database connection error: %s", err.Error())
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Database connection error: %s", err.Error())
	}

	return &PostgresStorage{Database: db}
}

func (self PostgresStorage) Migrate() {
	driver, err := postgres.WithInstance(self.Database, &postgres.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to create migrator 1: %v", err.Error()))
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		"file://internal/db/migrations",
		"postgres", driver)

	if err != nil {
		panic(fmt.Sprintf("Failed to create migrator 2: %v", err.Error()))
	}

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		panic(fmt.Sprintf("Failed to apply migrations: %v", err.Error()))
	}

	fmt.Println("Migrations applied successfully")
}

/**
func (self PostgresStorage) Seed() {
	// tender
	tenderRep := NewTenderRepository(&self)
	roadTender, err := tenderRep.CreateTender("Big road", "Avito need big road!", enums.ServiceTypeConstruction, alexAvito)
	_, err = tenderRep.CreateTender("Food", "Avito need some food", enums.ServiceTypeDelivery, ivanAvito)
	_, err = tenderRep.CreateTender("Brick production", "Lenta wants some bricks", enums.ServiceTypeConstruction, danilLenta)
	handleSeedError(err)

	// bid
	bidRep := NewBidRepository(&self)
	_, err = bidRep.CreateBid("We can build road", enums.AuthorTypeOrganization, roadTender, alexEmplId)
	handleSeedError(err)

	println("Seeded successfully")
}

func handleSeedError(err error) {
	if err != nil {
		log.Fatalf("Seed error: %s", err.Error())
	}
}
*/
