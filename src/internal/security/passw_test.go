package security

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type HashTestSuite struct {
	suite.Suite
}

var password string
var hash string

func (suite *HashTestSuite) SetupSuite() {
	// Caricamento environment
	err := godotenv.Load("/home/ubuntu/go/ballroom/.env")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Errore di caricamento del file '.env': %v", err)
		panic(err)
	}
	//password = "p4ssw0rd.t3st"
	password = "12345678"
}

func (suite *HashTestSuite) TeardownSuite() {
}

func (suite *HashTestSuite) TestSetPassword() {
	h, err := HashPassword(password)
	suite.Nil(err)
	hash = h
	fmt.Printf("Hash: %v\n", hash)
}

func (suite *HashTestSuite) TestVerifyPassword() {
	fmt.Printf("Password %v/Hash: %v\n", password, hash)
	ok := VerifyPassword(password, hash)
	suite.Equal(true, ok)
}

func (suite *HashTestSuite) TestResetPassword() {
	pwd, err := ResetPassword()
	suite.Nil(err)
	fmt.Printf("Nuova password %v\n", pwd)
}

func TestHashTestSuite(t *testing.T) {
	suite.Run(t, new(HashTestSuite))
}
