package domain

/*
	Corsi, programmazione e prenotazioni estemporanee
*/
import "time"

// Tipo corso di insegnamento
//   - Regolare: corso istituito dalla scuola con programmazione regolare destinato agli iscritti alla scuola
//   - Stage   : stage formativo di uno o più giorni aperto anche a non iscritti
//   - Privata : lezione privata di un insegnante con uno o più allievi; insegnante e/o allievi potrebbero essere esterni alla scuola
type TipoCorso int

const (
	Regolare TipoCorso = iota + 1
	Stage
	Privata
)

func (t TipoCorso) String() string {
	return [...]string{"Regolare", "Stage", "Privata"}[t-1]
}

// Giorni della settimana
type Settimana int

const (
	Lunedi Settimana = iota + 1
	Martedi
	Mercoledi
	Giovedi
	Venerdi
	Sabato
	Domenica
)

func (s Settimana) String() string {
	return [...]string{"Lunedì", "Martedì", "Mercoledì", "Giovedì", "Venerdì", "Sabato", "Domenica"}[s-1]
}

type Corso struct {
	Model
	Tipo       TipoCorso  `db:"tipo" json:"Tipo"`
	Titolo     string     `db:"titolo" json:"Titolo"`
	IdSala     Uuid       `db:"id_sala" json:"Id_Sala"`
	DataInizio *time.Time `db:"data_inizio" json:"Data_Inizio"`
	DataFine   *time.Time `db:"data_fine" json:"Data_Fine"`
	Orario     []Orario
}

type Orario struct {
	Giorno    Settimana
	OraInizio *time.Time `db:"ora_inizio" json:"Ora_Inizio"`
	OraFine   *time.Time `db:"ora_fine" json:"Ora_Fine"`
}

type Programmazione struct {
	Model
	IdCorso      Uuid `db:"id_corso" json:"Id_Corso"`
	IdAllievo    Uuid `db:"id_allievo" json:"Id_Allievo"`
	IdInsegnante Uuid `db:"id_insegnante" json:"Id_Insegnante"`
}

type CorsoDAO interface {
	Find(id Uuid) ([]*Corso, error)
	List() ([]*Corso, error)
	Create(c *Corso) (Uuid, error)
	Delete(id Uuid) error
}

type ProgrammazioneDAO interface {
	FindByCorso(idCorso Uuid) (*Corso, []*Persona, []*Persona, error)
	FindByAllievo(idAllievo Uuid) ([]*Corso, []*Persona, *Persona, error)
	FindByInsegnante(idInsegnante Uuid) ([]*Corso, *Persona, []*Persona, error)
	Create(idCorso Uuid, idAllievo []Uuid, idInsegnante []Uuid) (Uuid, error)
	Delete(id Uuid) error
}
