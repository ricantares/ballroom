package domain

// Scuola
type Scuola struct {
	Model
	Nome          string `db:"nome" json:"Nome"`
	CodiceFiscale string `db:"codice_fiscale" json:"Codice_Fiscale"`
}

// Sale di allenamento
type Sala struct {
	Model
	Nome string `db:"nome" json:"Nome"`
}
