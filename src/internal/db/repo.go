package db

// Interfacce di servizio per l'accesso ai dati

import "ricantares.com/ballroom/src/internal/domain"

// Interfacce di servizio che uno specifico Repository deve implementare
type DB interface {
	GetScuola() (domain.Scuola, error)
	GetSala(domain.Uuid) (domain.Sala, error)
	ListSala() ([]domain.Sala, error)
	CreateSala(domain.Sala) (result domain.Sala, err error)
	UpdateSala(domain.Sala) (domain.Sala, error)
	DeleteSala(domain.Uuid) (domain.Uuid, error)
	GetUtente(domain.Uuid) (domain.Utente, error)
	GetUtenteByName(string) (domain.Utente, error)
	ListUtente() ([]domain.Utente, error)
	ListUtenteByRuolo(domain.TipoRuolo) ([]domain.Utente, error)
	CreateUtente(domain.Utente) (result domain.Utente, err error)
	UpdateUtente(domain.Utente) (domain.Utente, error)
	DeleteUtente(domain.Uuid) (domain.Uuid, error)
}

type Repository struct {
	db DB
}

func NewRepository(db DB) *Repository {
	return &Repository{db: db}
}

// Per una precisa scelta progettuale, l'applicazione gestisce una singola scuola
// Pertanto, la scuola e' un oggetto unico e immutabile creato in fase di inizializzazione del sistema
// e ogni altra entita' e' implicitamente dipendente (i.e. nessun attributo di relazione)
func (r *Repository) GetScuola() (domain.Scuola, error) {
	return r.db.GetScuola()
}

/*
Gestione sala
*/
func (r *Repository) GetSala(id domain.Uuid) (domain.Sala, error) {
	return r.db.GetSala(id)
}

func (r *Repository) ListSala() ([]domain.Sala, error) {
	return r.db.ListSala()
}

func (r *Repository) CreateSala(newSala domain.Sala) (domain.Sala, error) {
	return r.db.CreateSala(newSala)
}

func (r *Repository) UpdateSala(updSala domain.Sala) (domain.Sala, error) {
	return r.db.UpdateSala(updSala)
}

func (r *Repository) DeleteSala(id domain.Uuid) (domain.Uuid, error) {
	return r.db.DeleteSala(id)
}

/*
Gestione utente
*/
func (r *Repository) GetUtente(id domain.Uuid) (domain.Utente, error) {
	return r.db.GetUtente(id)
}

func (r *Repository) GetUtenteByName(name string) (domain.Utente, error) {
	return r.db.GetUtenteByName(name)
}

func (r *Repository) ListUtente() ([]domain.Utente, error) {
	return r.db.ListUtente()
}

func (r *Repository) ListUtenteByRuolo(domain.TipoRuolo) ([]domain.Utente, error) {
	return r.db.ListUtente()
}

func (r *Repository) CreateUtente(newUtente domain.Utente) (domain.Utente, error) {
	return r.db.CreateUtente(newUtente)
}

func (r *Repository) UpdateUtente(updUtente domain.Utente) (domain.Utente, error) {
	return r.db.UpdateUtente(updUtente)
}

func (r *Repository) DeleteUtente(id domain.Uuid) (domain.Uuid, error) {
	return r.db.DeleteUtente(id)
}
