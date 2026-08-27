package db

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"ricantares.com/ballroom/internal/domain"
	"ricantares.com/ballroom/internal/logger"
	"ricantares.com/ballroom/internal/security"
)

type RepoTestSuite struct {
	suite.Suite
}

var repository *Repository
var database *Database
var nomeutente string
var ruoloutente domain.TipoRuolo
var idutente domain.Uuid

func (suite *RepoTestSuite) SetupSuite() {

	// Caricamento environment
	err := godotenv.Load("/home/ubuntu/go/ballroom/.env")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Errore di caricamento del file '.env': %v", err)
		panic(err)
	}

	// set timezone
	loc, err := time.LoadLocation(os.Getenv("TIMEZONE"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Errore di impostazione TIMEZONE: %v", err)
		panic(err)
	}
	time.Local = loc

	logger.NewLog(os.Stdout, "Debug")

	var databaseUrl = os.Getenv("DBCONNECT")
	var ctxBg = context.Background()
	conn, err := NewDbConnection(ctxBg, databaseUrl)
	if err != nil {
		logger.LogError(fmt.Sprintf("Unable to connect to database: %v - err: %v", databaseUrl, err))
		os.Exit(1)
	}
	database = conn
	logger.LogInfo(fmt.Sprintf("Connesso al db %v", databaseUrl))

	repository = NewRepository(database)

	nomeutente = "Admin"
	ruoloutente = domain.Admin
}

func (suite *RepoTestSuite) TeardownSuite() {
	logger.LogInfo("Chiusura connessioni db")
	database.Close()
}

func (suite *RepoTestSuite) TestGetScuola() {
	scuola, err := repository.GetScuola()
	suite.Nil(err)
	suite.Equal("Dove si balla", scuola.Nome)
}

func (suite *RepoTestSuite) TestGetSala() {
	sala, err := repository.GetSala(1)
	suite.Nil(err)
	suite.Equal("Sala principale", sala.Nome)
}

func (suite *RepoTestSuite) TestCreateUtente() {

	logger.LogInfo("Crea utente")

	pwd, err := security.ResetPassword()
	if err != nil {
		logger.LogError(fmt.Sprintf("Unable to create password: %v", err))
		os.Exit(1)
	}
	// solo per test
	pwd = "prova123"

	hashpwd, err := security.HashPassword(pwd)
	if err != nil {
		logger.LogError(fmt.Sprintf("Unable to hash password: %v", err))
		os.Exit(1)
	}
	scad := time.Now()
	newUtente := domain.Utente{
		Nome:          nomeutente,
		Password:      hashpwd,
		Temp_Password: hashpwd,
		Scad_Password: &scad,
		Tipo_Ruolo:    ruoloutente,
	}
	utente, err := repository.CreateUtente(newUtente)
	suite.Nil(err)
	suite.Equal(nomeutente, utente.Nome)
	suite.Equal(utente.Password, utente.Temp_Password)

	idutente = utente.Id

	logger.LogInfo(fmt.Sprintf("Utente creato %v", utente))
}

func (suite *RepoTestSuite) TestGetUtente() {

	logger.LogInfo("Get utente " + strconv.Itoa(int(idutente)))

	utente, err := repository.GetUtente(idutente)
	suite.Nil(err)
	suite.Equal(nomeutente, utente.Nome)
	suite.Equal(domain.Admin, utente.Tipo_Ruolo)
}

func (suite *RepoTestSuite) TestGetUtenteByName() {

	logger.LogInfo("Get utente " + nomeutente)

	utente, err := repository.GetUtenteByName(nomeutente)
	suite.Nil(err)
	suite.Equal(nomeutente, utente.Nome)
	suite.Equal(domain.Admin, utente.Tipo_Ruolo)
}

func (suite *RepoTestSuite) TestListUtenteByRuolo() {

	logger.LogInfo("List utente " + ruoloutente.String())

	utente, err := repository.ListUtenteByRuolo(ruoloutente)
	suite.Nil(err)
	suite.Greater(len(utente), 0)
}

/*
func (suite *RepoTestSuite) TestDeleteUtente() {

		logger.LogInfo(fmt.Sprintf("Delete utente"))

		id, err := repository.DeleteUtente(idutente)
		suite.Nil(err)
		suite.Equal(id, idutente)
	}
*/
func TestRepoTestSuite(t *testing.T) {
	suite.Run(t, new(RepoTestSuite))
}
