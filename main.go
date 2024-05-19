package main

import (
	"billing-engine/internal/config"
	"billing-engine/internal/delivery/http"
	"billing-engine/internal/repository"
	"billing-engine/internal/usecase"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.LoadConfig()

	dsn := cfg.DatabaseURL()

	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer dbpool.Close()

	loanRepo := repository.NewLoanRepository(dbpool)
	loanUsecase := usecase.NewLoanUsecase(loanRepo)

	e := echo.New()
	http.NewLoanHandler(e, loanUsecase)

	e.Logger.Fatal(e.Start(":8080"))
}
