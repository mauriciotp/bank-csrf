package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/mauriciotp/bank-csrf/internal/api"
	"github.com/mauriciotp/bank-csrf/internal/store/pgstore"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	ctx := context.Background()

	port := os.Getenv("PORT")
	databaseUser := os.Getenv("DATABASE_USER")
	databasePassword := os.Getenv("DATABASE_PASSWORD")
	databaseHost := os.Getenv("DATABASE_HOST")
	databasePort := os.Getenv("DATABASE_PORT")
	databaseName := os.Getenv("DATABASE_NAME")

	pool, err := pgxpool.New(
		ctx,
		fmt.Sprintf(
			"user=%s password=%s host=%s port=%s dbname=%s",
			databaseUser,
			databasePassword,
			databaseHost,
			databasePort,
			databaseName,
		))
	if err != nil {
		panic(err)
	}

	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	a := api.API{
		Router:  chi.NewMux(),
		Pool:    pool,
		Queries: pgstore.New(pool),
	}

	slog.Info(fmt.Sprintf("Starting Server on port %s 🚀 ", port))
	if err := http.ListenAndServe(fmt.Sprintf("localhost:%s", port), a.Router); err != nil {
		panic(err)
	}
}
