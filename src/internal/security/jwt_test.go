package security

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"ricantares.com/ballroom/src/internal/domain"
	"ricantares.com/ballroom/src/internal/logger"
)

var tokenString string
var claims JwtCustomClaims
var logfile *os.File

type JwtTestSuite struct {
	suite.Suite
}

func (suite *JwtTestSuite) SetupSuite() {
	// Caricamento environment
	err := godotenv.Load("/home/ubuntu/go/ballroom/.env")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Errore di caricamento del file '.env': %v", err)
		panic(err)
	}

	loc, err := time.LoadLocation(os.Getenv("TIMEZONE"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Errore di impostazione TIMEZONE: %v", err)
		panic(err)
	}
	time.Local = loc

	outfile, err := os.Create(os.Getenv("LOGFILE"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Errore di caricamento del file di log: %v", err)
		panic(err)
	}

	logger.NewLog(outfile, os.Getenv("LOGLEVEL"))
	logfile = outfile
}

func (suite *JwtTestSuite) TeardownSuite() {
	logfile.Close()
}

func (suite *JwtTestSuite) TestGeneraToken() {
	t, err := GeneraToken("Admin", domain.Admin)
	suite.Nil(err)
	tokenString = t.AccessToken
	fmt.Printf("Token: %v\n", tokenString)
}

func (suite *JwtTestSuite) TestTokenClaims() {
	c, err := GetTokenClaims(tokenString)
	suite.Nil(err)
	claims = *c
	fmt.Printf("Claims: %v\n", claims)

	exp := TokenScaduto(&claims)
	suite.False(exp)
	fmt.Printf("Token valido\n")

	r := GetTockenRole(&claims)
	suite.Equal(domain.Admin, r)
	fmt.Printf("Ruolo: %v\n", r)
}

func TestJwtTestSuite(t *testing.T) {
	suite.Run(t, new(JwtTestSuite))
}
