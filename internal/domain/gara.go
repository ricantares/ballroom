package domain

import "time"

// Gara di ballo
type TipoGara int

const (
	Locale TipoGara = iota + 1
	Nazionale
	Internazionale
)

func (t TipoGara) String() string {
	return [...]string{"Locale", "Nazionale", "Internazionale"}[t-1]
}

type Gara struct {
	Model
	Tipo       TipoGara   `db:"tipo" json:"Tipo"`
	Titolo     string     `db:"titolo" json:"Titolo"`
	Luogo      string     `db:"id_sala" json:"Id_Sala"`
	DataInizio *time.Time `db:"data_inizio" json:"Data_Inizio"`
	DataFine   *time.Time `db:"data_fine" json:"Data_Fine"`
}

type Iscrizione struct {
	Model
	IdGara    Uuid   `db:"id_gara" json:"Id_Gara"`
	IdAllievo []Uuid `db:"id_allievo" json:"Id_Allievo"`
}

type GaraDAO interface {
	Find(id Uuid) ([]*Gara, error)
	List() ([]*Gara, error)
	Create(c *Gara) (Uuid, error)
	Delete(id Uuid) error
}

type IscrizioneDAO interface {
	FindByGara(idGara Uuid) (*Gara, []*Iscrizione, []*Persona, error)
	FindByAllievo(idAllievo Uuid) ([]*Gara, []*Iscrizione, *Persona, error)
	Create(idGara Uuid, idAllievo []Uuid) (Uuid, error)
	Delete(id Uuid) error
}
