package domain

/*
Utenze e ruoli del sistema
*/

import "time"

// Ruolo utente per l'accesso alle funzioni dell'applicazione.
//   - Admin    : amministra il sistema creando le altre utenze;
//   - Direzione: direttore artistico che accede a tutte le funzionalità tranne quelle di "admin";
//   - Iscritto : visualizza i propri dati anagrafici, le prenotazioni, i corsi;
//   - Maestro  : visualizza i propri dati anagrafici, le prenotazioni, i corsi;
//   - Staff    : gestisce le anagrafiche, le prenotazioni, i corsi;
type TipoRuolo int

const (
	Admin TipoRuolo = iota + 1
	Direzione
	Iscritto
	Maestro
	Staff
)

func (t TipoRuolo) String() string {
	return [...]string{"Admin", "Direzione", "Iscritto", "Maestro", "Staff"}[t-1]
}

// Utente del sistema
type Utente struct {
	Model
	Nome          string     `db:"nome" json:"Nome"`
	Password      string     `db:"password" json:"Password"`
	Temp_Password string     `db:"temp_password" json:"Temp_Password"`
	Scad_Password *time.Time `db:"scad_password" json:"Scad_Password"`
	Tipo_Ruolo    TipoRuolo  `db:"tipo_ruolo" json:"Tipo_Ruolo"`
}

type UtenteDAO interface {
	Find(id Uuid) (*Utente, error)
	FindByName(name string) (*Utente, error)
	List() ([]*Utente, error)
	Create(a *Utente) (Uuid, error)
	Delete(id Uuid) error
}
