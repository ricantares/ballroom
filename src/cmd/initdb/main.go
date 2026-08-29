package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"ricantares.com/ballroom/src/internal/security"
)

func main() {
	resourcesDir := flag.String("resources-dir", "src/resources", "directory containing the SQL files")
	envFile := flag.String("env-file", "src/.env", "environment file containing DBCONNECT")
	flag.Parse()

	if err := initializeDatabase(*envFile, *resourcesDir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("Database inizializzato correttamente.")
}

func initializeDatabase(envFile, resourcesDir string) error {
	env, err := godotenv.Read(envFile)
	if err != nil {
		return fmt.Errorf("errore nel caricamento del file .env %q: %w", envFile, err)
	}

	databaseURL := env["DBCONNECT"]
	if databaseURL == "" {
		return fmt.Errorf("errore: DBCONNECT non definita in %q", envFile)
	}

	ctx := context.Background()
	database, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("errore nella connessione al database: %w", err)
	}
	defer database.Close()

	if err := database.Ping(ctx); err != nil {
		return fmt.Errorf("errore nel ping del database: %w", err)
	}

	schema, err := os.ReadFile(filepath.Join(resourcesDir, "schema.sql"))
	if err != nil {
		return fmt.Errorf("errore nella lettura di schema.sql: %w", err)
	}
	seed, err := os.ReadFile(filepath.Join(resourcesDir, "init_scuola.sql"))
	if err != nil {
		return fmt.Errorf("errore nella lettura di init_scuola.sql: %w", err)
	}

	transaction, err := database.Begin(ctx)
	if err != nil {
		return fmt.Errorf("errore nell'avvio della transazione: %w", err)
	}
	defer transaction.Rollback(ctx)

	if _, err := transaction.Exec(ctx, string(schema)); err != nil {
		return fmt.Errorf("errore nell'esecuzione di schema.sql: %w", err)
	}
	if _, err := transaction.Exec(ctx, string(seed)); err != nil {
		return fmt.Errorf("errore nell'esecuzione di init_scuola.sql: %w", err)
	}

	hashPassword, err := security.HashPassword("Password123")
	if err != nil {
		return fmt.Errorf("errore nella generazione dell'hash della password admin: %w", err)
	}

	if _, err := transaction.Exec(ctx, `
		INSERT INTO utente (nome, password, temp_password, tipo_ruolo)
		SELECT $1, $2, $2, 1
		WHERE NOT EXISTS (
			SELECT 1 FROM utente WHERE nome = $1
		)
	`, "Admin", hashPassword); err != nil {
		return fmt.Errorf("errore nell'inserimento dell'utente admin: %w", err)
	}

	if err := transaction.Commit(ctx); err != nil {
		return fmt.Errorf("errore nel commit dell'inizializzazione: %w", err)
	}

	return nil
}
