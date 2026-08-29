package db

/*
Definisce e implementa i metodi di accesso al db postgresg
Il context e il pool di connessioni vengono gestiti con variabili globali
così da non dover essere passati come parametri a ogni chiamata alle funzioni CRUD
*/
import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"ricantares.com/ballroom/src/internal/domain"
	"ricantares.com/ballroom/src/internal/logger"
)

// Connessione o pool di connessioni alla base dati
type Database struct {
	DB *pgxpool.Pool
}

var dbInstance *Database
var dbContext context.Context

// Inizializzazione connessioni alla base dati
func NewDbConnection(ctx context.Context, databaseUri string) (*Database, error) {
	conn, err := pgxpool.New(ctx, databaseUri)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	dbInstance = &Database{conn}
	dbContext = ctx
	return dbInstance, nil
}

// Chiusura pool di connessioni
func (db *Database) Close() {
	db.DB.Close()
	logger.LogDebug("Disconnesso dal db")
}

// Create
func Create[T any](query string, args pgx.NamedArgs) (resultSet T, e error) {

	rows, _ := dbInstance.DB.Query(dbContext, query, args)
	resultSet, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[T])
	if err != nil {
		return resultSet, err
	}

	return resultSet, nil

}

// Read
func Read[T any](query string, args pgx.NamedArgs) (resultSet []T, e error) {

	rows, _ := dbInstance.DB.Query(dbContext, query, args)
	resultSet, err := pgx.CollectRows(rows, pgx.RowToStructByName[T])
	if err != nil {
		return nil, err
	}

	return resultSet, nil

}

// Update
func Update[T any](query string, args pgx.NamedArgs) (resultSet T, e error) {

	rows, _ := dbInstance.DB.Query(dbContext, query, args)
	resultSet, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[T])
	if err != nil {
		return resultSet, err
	}

	return resultSet, nil

}

// Delete
func Delete[T any](query string, args pgx.NamedArgs) (recordId domain.Uuid, e error) {

	err := dbInstance.DB.QueryRow(dbContext, query, args).Scan(&recordId)
	if err != nil {
		return 0, err
	}
	return recordId, nil

}
