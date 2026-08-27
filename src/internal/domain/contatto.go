package domain

/*
Contatto di una persona (una persona può avere più contatti, solo uno come preferito)
*/

// Tipo Contatto
//   - Mail
//   - Telefono
type TipoContatto int

const (
	Telefono TipoContatto = iota + 1
	Mail
)

func (t TipoContatto) String() string {
	return [...]string{"Telefono", "Mail"}[t-1]
}

type Contatto struct {
	Model
	Tipo           TipoContatto `db:"tipo" json:"Tipo"`
	Valore         string       `db:"valore" json:"Valore"`
	IdProprietario Uuid         `db:"id_owner" json:"Id_Proprietario"`
	Preferito      bool         `db:"id_preferito" json:"Id_Preferito"`
}

// Uno o piu' contatti sono sempre riferiti a un "Proprietario" (scuola, allievi, insegnanti)
type ContattoDAO interface {
	List(IdProprietario Uuid) ([]*Contatto, error)
	Create(IdProprietario Uuid, a *Contatto) (Uuid, error)
	Delete(id Uuid) error
}
