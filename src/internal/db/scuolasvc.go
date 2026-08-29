package db

// Implementazione Postgresql delle interfacce di servizio per l'accesso ai dati della scuola e delle relative sale

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"ricantares.com/ballroom/src/internal/domain"
)

// Scuola services
func (r *Database) GetScuola() (domain.Scuola, error) {
	var query = "SELECT * FROM scuola"
	var args = pgx.NamedArgs{}
	result, err := Read[domain.Scuola](query, args)
	if len(result) > 0 {
		return result[0], err
	}
	return domain.Scuola{}, errors.New("not found")
}

// Sala services
func (r *Database) GetSala(id domain.Uuid) (resultset domain.Sala, err error) {
	var query = "SELECT * FROM sala WHERE Id = @id AND Deleted = @del"
	var args = pgx.NamedArgs{"id": id, "del": false}
	result, err := Read[domain.Sala](query, args)
	if len(result) > 0 {
		return result[0], err
	}
	return domain.Sala{}, errors.New("not found")
}

func (r *Database) ListSala() (resultset []domain.Sala, err error) {
	var query = "SELECT * FROM sala WHERE Deleted = @del"
	var args = pgx.NamedArgs{"del": false}
	result, err := Read[domain.Sala](query, args)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Database) CreateSala(newSala domain.Sala) (result domain.Sala, err error) {
	var query = "INSERT INTO sala (nome) VALUES (@nome) RETURNING id, created_at, updated_at, deleted, nome"
	var args = pgx.NamedArgs{"nome": newSala.Nome}
	new, err := Create[domain.Sala](query, args)
	if err != nil {
		return new, err
	}
	return new, nil
}

func (r *Database) UpdateSala(updSala domain.Sala) (result domain.Sala, err error) {
	var query = "UPDATE sala SET nome = @nome, updated_at = now() WHERE id = @id RETURNING id, created_at, updated_at, deleted, nome"
	var args = pgx.NamedArgs{"id": updSala.Model.Id, "nome": updSala.Nome}
	new, err := Update[domain.Sala](query, args)
	if err != nil {
		return new, err
	}
	return new, nil
}

func (r *Database) DeleteSala(id domain.Uuid) (result domain.Uuid, err error) {
	var query = "DELETE FROM sala WHERE id = @id RETURNING id"
	var args = pgx.NamedArgs{"id": id}
	deleted, err := Delete[domain.Sala](query, args)
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
