package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

type Borrower struct {
	ID        int32
	Createdat time.Time
	Updatedat time.Time
	Deletedat *time.Time
	Name      string
	Email     string
	Phone     string
}

func main() {
	// Load environment variables from .env file using Viper
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Database connection string
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_HOST"),
		viper.GetString("DB_PORT"),
		viper.GetString("DB_NAME"))

	// Connect to the database
	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	// Seed data
	err = seedBorrowers(dbpool, 10) // Change the number to the desired number of records
	if err != nil {
		log.Fatalf("Failed to seed borrowers: %v\n", err)
	}

	fmt.Println("Successfully seeded borrowers")
}

func seedBorrowers(dbpool *pgxpool.Pool, count int) error {
	for i := 0; i < count; i++ {
		now := time.Now()
		borrower := Borrower{
			Createdat: now,
			Updatedat: now,
			Deletedat: nil,
			Name:      faker.Name(),
			Email:     faker.Email(),
			Phone:     faker.Phonenumber(),
		}

		_, err := dbpool.Exec(context.Background(), `
			INSERT INTO borrowers (createdat, updatedat, deletedat, name, email, phone)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			borrower.Createdat, borrower.Updatedat, borrower.Deletedat, borrower.Name, borrower.Email, borrower.Phone)
		if err != nil {
			return fmt.Errorf("failed to insert borrower: %w", err)
		}
	}
	return nil
}
