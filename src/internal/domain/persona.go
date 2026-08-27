package domain

import "time"

// Tipi di anagrafiche gestite dall'applicazione.
//   - Insegnante : maestro di ballo di uno o più corsi;
//   - Allievo    : allievo partecipante a uno o più corsi di ballo;
//   - Ospite     : insegnante o allievo occasionale (stage, gare interne, altro);
type TipoPersona int

const (
	Insegnante TipoPersona = iota + 1
	Allievo
	Ospite
)

func (t TipoPersona) String() string {
	return [...]string{"Insegnante", "Allievo", "Ospite"}[t-1]
}

// Anagrafica persona
//   - Nome e cognome dati obbligatori
//   - DataNascita e CodiceFiscale opzionali
//   - IdCrew indicato in caso l'allievo partecipi a corsi che prevedono uno o più partner
//   - IdUtenteSistema nel caso in cui la persona sia registrata come utente del sistema
type Persona struct {
	Model
	Tipo            TipoPersona `db:"tipo_persona" json:"Tipo_Persona"`
	Nome            string      `db:"nome" json:"Nome"`
	Cognome         string      `db:"cognome" json:"Cognome"`
	DataNascita     *time.Time  `db:"data_nascita" json:"Data_Nascita"`
	CodiceFiscale   string      `db:"codice_fiscale" json:"Codice_Fiscale"`
	IdCrew          Uuid        `db:"id_crew" json:"Id_Crew"`
	IdUtenteSistema Uuid        `db:"id_utente" json:"Id_Utente"`
}

type Crew struct {
	Model
	IdAllievi []Uuid `db:"id_allievi" json:"Id_Allievi"`
}

type PersonaDAO interface {
	Find(id Uuid) (*Persona, error)
	List() ([]*Persona, error)
	Create(a *Persona) (Uuid, error)
	Delete(id int32) error
}
