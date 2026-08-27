package domain

/*
Certificazione medica dell'allievo
*/
import "time"

// Certificato medico per l'attivita'
type Certificazione struct {
	Model
	Agonismo     bool       `db:"certificato_agonistico" json:"Certificato_Agonistico"`
	DataRilascio *time.Time `db:"data_rilascio" json:"Data_Rilascio"`
	DataScadenza *time.Time `db:"data_scadenza" json:"Data_Scadenza"`
	Medico       string     `db:"medico" json:"Medico"`
	Protocollo   string     `db:"protocollo" json:"Protocollo"`
	IdAllievo    Uuid       `db:"id_allievo" json:"Id_Allievo"`
}

// Certificato medico di un allievo
type CertificazioneAllievo struct {
	Certificato Certificazione
	NomeAllievo string `db:"nome_cognome_allievo" json:"Nome_Cognome_Allievo"`
}

type CertificazioneDAO interface {
	Find(idAllievo Uuid) (*CertificazioneAllievo, error)
	List() ([]*CertificazioneAllievo, error)
	Create(a *Certificazione) (Uuid, error)
	Delete(id Uuid) error
}
