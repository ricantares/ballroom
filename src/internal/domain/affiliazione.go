package domain

import "time"

// Affiliazione (Federazione, Associazione o altro cui la scuola e' affiliata)
type Affiliazione struct {
	Model
	Codice                         string     `db:"codice" json:"Codice"`
	Nome                           string     `db:"nome" json:"Nome"`
	DataPrimaIscrizione            *time.Time `db:"data_prima_iscrizione" json:"Data_Prima_Iscrizione"`
	DataScadenzaIscrizioneCorrente *time.Time `db:"data_scadenza_iscrizione_ac" json:"Data_Scadenza_Iscrizione_AC"`
	DataInizioAnnoSportivo         *time.Time `db:"data_inizio_as" json:"Data_Inizio_AS"`
	DataFineAnnoSportivo           *time.Time `db:"data_fine_as" json:"Data_Fine_AS"`
}

// Gli allievi sono tesserati a una o piu' delle organizzazioni cui la scuola e' affiliata
type Tesseramento struct {
	Model
	IdTesserato                 int32      `db:"id_tesserato" json:"Id_Tesserato"`
	IdAffiliazione              int32      `db:"id_affiliazione" json:"Id_Affiliazione"`
	TesseratoAgonistico         bool       `db:"tesserato_agonistico" json:"Tesserato_Agonistico"`
	NumeroTessera               string     `db:"nro_tessera" json:"Nro_Tessera"`
	DataPrimoTesseramento       *time.Time `db:"data_primo_tesseramento" json:"Data_Primo_Tesseramento"`
	DataScadenzaTesseraCorrente *time.Time `db:"data_scadenza_tessera_ac" json:"Data_Scadenza_Tessera_AC"`
}

// Tesseramento alunno/affiliazione
type TesseraAffiliato struct {
	Tesseramento     Tesseramento
	NomeAffiliazione string `db:"nome_affiliazione" json:"Nome_Affiliazione"`
	NomeTesserato    string `db:"nome_tesserato" json:"Nome_Tesserato"`
}

type AffiliazioneDAO interface {
	Find(id Uuid) (*Affiliazione, error)
	List() ([]*Affiliazione, error)
	Create(a *Affiliazione) (Uuid, error)
	Delete(id Uuid) error
}

type TesseramentoDAO interface {
	Find(idTesserato Uuid) ([]*TesseraAffiliato, error)
	List(idAffiliazione Uuid) ([]*TesseraAffiliato, error)
	Create(idTesserato Uuid, a *Tesseramento) (Uuid, error)
	Delete(id Uuid) error
}
