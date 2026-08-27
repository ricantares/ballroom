package db

// Implementazione Postgresql delle interfacce di servizio per l'accesso ai dati degli utenti

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"ricantares.com/ballroom/internal/domain"
)

// Utente services
func (r *Database) GetUtente(id domain.Uuid) (resultset domain.Utente, err error) {
	var query = "SELECT * FROM utente WHERE Id = @id AND Deleted = @del"
	var args = pgx.NamedArgs{"id": id, "del": false}
	//logger.LogDebug("boh: " + query + " con id=" + strconv.Itoa(int(id)))
	result, err := Read[domain.Utente](query, args)
	if len(result) > 0 {
		return result[0], err
	}
	return domain.Utente{}, errors.New("not found")
}

func (r *Database) GetUtenteByName(name string) (resultset domain.Utente, err error) {
	var query = "SELECT * FROM utente WHERE Nome = @nome AND Deleted = @del"
	var args = pgx.NamedArgs{"nome": name, "del": false}
	result, err := Read[domain.Utente](query, args)
	if len(result) > 0 {
		return result[0], err
	}
	return domain.Utente{}, errors.New("not found")
}

func (r *Database) ListUtente() (resultset []domain.Utente, err error) {
	var query = "SELECT * FROM utente WHERE Deleted = @del"
	var args = pgx.NamedArgs{"del": false}
	result, err := Read[domain.Utente](query, args)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Database) ListUtenteByRuolo(ruolo domain.TipoRuolo) (resultset []domain.Utente, err error) {
	var query = "SELECT * FROM utente WHERE Tipo_Ruolo = @ruolo AND Deleted = @del"
	var args = pgx.NamedArgs{"ruolo": ruolo, "del": false}
	result, err := Read[domain.Utente](query, args)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Database) CreateUtente(newUtente domain.Utente) (result domain.Utente, err error) {
	var query = "INSERT INTO utente (nome, password, temp_password, scad_password, tipo_ruolo) VALUES (@nome, @pwd, @tmppwd, @scadpwd, @ruolo) RETURNING id, created_at, updated_at, deleted, nome, password, temp_password, scad_password, tipo_ruolo"
	var args = pgx.NamedArgs{"nome": newUtente.Nome, "pwd": newUtente.Password, "tmppwd": newUtente.Temp_Password, "scadpwd": newUtente.Scad_Password, "ruolo": newUtente.Tipo_Ruolo}
	new, err := Create[domain.Utente](query, args)
	if err != nil {
		return new, err
	}
	return new, nil
}

func (r *Database) UpdateUtente(updUtente domain.Utente) (result domain.Utente, err error) {
	var query = "UPDATE utente SET updated_at = now() WHERE id = @id RETURNING id, created_at, updated_at, deleted, nome, password, scad_password, tipo_ruolo"
	var args = pgx.NamedArgs{"id": updUtente.Model.Id, "pwd": updUtente.Password, "scadpwd": updUtente.Scad_Password, "ruolo": updUtente.Tipo_Ruolo}
	new, err := Update[domain.Utente](query, args)
	if err != nil {
		return new, err
	}
	return new, nil
}

func (r *Database) DeleteUtente(id domain.Uuid) (result domain.Uuid, err error) {
	var query = "DELETE FROM utente WHERE id = @id RETURNING id"
	var args = pgx.NamedArgs{"id": id}
	deleted, err := Delete[domain.Utente](query, args)
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
