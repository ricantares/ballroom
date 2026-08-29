package api

import (
	"encoding/json"
	"os"

	"github.com/gin-gonic/gin"
	"ricantares.com/ballroom/src/internal/db"
	"ricantares.com/ballroom/src/internal/domain"
	"ricantares.com/ballroom/src/internal/rest"
)

// Interfacce che il route handler deve implementare
type Router interface {
	RouterScuola
	GetUtente(*gin.Context) (domain.Utente, error)
	GetUtenteByName(*gin.Context) (domain.Utente, error)
	ListUtente(*gin.Context) ([]domain.Utente, error)
	ListUtenteByRuolo(*gin.Context) ([]domain.Utente, error)
	CreateUtente(*gin.Context) (result domain.Utente, err error)
	UpdateUtente(*gin.Context) (domain.Utente, error)
	DeleteUtente(*gin.Context) (domain.Uuid, error)
}

type RouteHandler struct {
	repo db.Repository
}

// NewRouteHandler creates a new RouteHandler with the given repository.
// The repository is used to interact with the data layer and perform
// CRUD operations for different entities.

func NewRouteHandler(r db.Repository) *RouteHandler {
	return &RouteHandler{repo: r}
}

type ApiResponse struct {
	Code    uint   `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type Embedded struct {
	Name string
	Data []rest.BaseResource
}

type ResourceWithEmbedds struct {
	Base    rest.BaseResource
	Embedds Embedded
}

// ritorna una risposta con una singola risorsa in formato HAL
func SimpleResponse(res rest.BaseResource) (map[string]any, error) {
	hal := rest.NewHateoas(os.Getenv("HTTP_ORIGIN"))
	hal.AddLinks(res.Links)
	hal.AddData(res.Data)
	jsonData, err := hal.ToJSON()
	if err != nil {
		return nil, err
	}

	// Convert the JSON response to a map
	var response map[string]any
	err = json.Unmarshal(jsonData, &response)

	return response, err
}

// ritorna una risposta con risorse multiple in formato HAL
func ResponseWithEmbedds(res ResourceWithEmbedds) (map[string]any, error) {
	hal := rest.NewHateoas(os.Getenv("HTTP_ORIGIN"))
	hal.AddLinks(res.Base.Links)
	hal.AddData(res.Base.Data)
	hal.AddEmbedded("sale", res.Embedds.Data)

	jsonData, err := hal.ToJSON()
	if err != nil {
		return nil, err
	}

	// Convert the JSON response to a map
	var response map[string]any
	err = json.Unmarshal(jsonData, &response)

	return response, err
}
