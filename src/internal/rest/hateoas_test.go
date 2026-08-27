package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

var hostName string

type localData struct {
	Campo1 string     `json:"campo1,omitempty"`
	Campo2 int        `json:"campo2,omitempty"`
	Campo3 bool       `json:"campo3,omitempty"`
	Campo4 *time.Time `json:"campo4,omitempty"`
}

type HateoasTestSuite struct {
	suite.Suite
}

func (suite *HateoasTestSuite) SetupSuite() {
	protocol := "http"
	host := "localhost"
	port := ":8081"
	hostName = protocol + "//" + host + port
}

func (suite *HateoasTestSuite) TeardownSuite() {
}

func (suite *HateoasTestSuite) TestSimpleHal() {
	h := NewHateoas(hostName)
	err := h.AddLink(Link{SELF, "/scuola"})
	suite.Nil(err)

	t := time.Now()
	d := localData{
		Campo1: "campo1",
		Campo2: 42,
		Campo3: true,
		Campo4: &t,
	}
	h.AddData(d)

	jsonData, err := h.ToJSON()
	suite.Nil(err)

	jsonPrettyPrint(jsonData)
}

func (suite *HateoasTestSuite) TestMultipleLink() {
	h := NewHateoas(hostName)
	err := h.AddLink(Link{SELF, "/scuola"})
	suite.Nil(err)

	err = h.AddLink(Link{"sale", "/scuola/sale"})
	suite.Nil(err)

	t := time.Now()
	d := localData{
		Campo1: "campo1",
		Campo2: 42,
		Campo3: true,
		Campo4: &t,
	}
	h.AddData(d)

	jsonData, err := h.ToJSON()
	suite.Nil(err)

	jsonPrettyPrint(jsonData)
}

func (suite *HateoasTestSuite) TestEmbedded() {
	type Sala struct {
		Id   int    `json:"-"`
		Nome string `json:"nome"`
	}

	h := NewHateoas(hostName)
	err := h.AddLink(Link{SELF, "/scuola/sale"})
	suite.Nil(err)
	err = h.AddLink(Link{FIRST, "/scuola/sale"})
	suite.Nil(err)
	err = h.AddLink(Link{LAST, "/scuola/sale"})
	suite.Nil(err)
	err = h.AddLink(Link{PREVIOUS, "/scuola/sale"})
	suite.Nil(err)
	err = h.AddLink(Link{NEXT, "/scuola/sale"})
	suite.Nil(err)

	t := time.Now()
	d := localData{
		Campo1: "Balla con noi",
		Campo2: 42,
		Campo3: true,
		Campo4: &t,
	}
	h.AddData(d)

	sale := []Sala{
		{1, "sala uno"},
		{2, "sala due"},
		{3, "sala tre"},
	}

	links := [][]Link{
		{{SELF, "/scuola/sala/1"}},
		{{SELF, "/scuola/sala/2"}},
		{{SELF, "/scuola/sala/3"}},
	}

	var embeddeds []BaseResource
	for i, sala := range sale {
		embedded := BaseResource{Links: links[i], Data: sala}
		embeddeds = append(embeddeds, embedded)
	}

	h.AddEmbedded("sale", embeddeds)

	p := Page{
		Size:  1,
		First: true,
		Last:  true,
		Pages: 1,
		Items: 3,
	}
	h.AddPage(p)

	jsonData, err := h.ToJSON()
	suite.Nil(err)

	jsonPrettyPrint(jsonData)
}

func TestHateoasTestSuite(t *testing.T) {
	suite.Run(t, new(HateoasTestSuite))
}

func jsonPrettyPrint(jsonData []byte) {

	var buf bytes.Buffer
	err := json.Indent(&buf, jsonData, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nJSON: \n %v \n", buf.String())

}
