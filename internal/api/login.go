package api

/*
Handler login
*/

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"ricantares.com/ballroom/internal/domain"
	"ricantares.com/ballroom/internal/logger"
	"ricantares.com/ballroom/internal/security"
)

// verifica i dati di login e restituisce un token jwt
func (h RouteHandler) HandleLogin(c *gin.Context) {
	// Get request body and convert it to domain entity
	var utente domain.Utente
	if err := c.ShouldBindJSON(&utente); err != nil {
		c.JSON(http.StatusBadRequest, ApiResponse{Code: http.StatusBadRequest, Message: err.Error(), Data: ""})

		return
	}

	// recupero dati utente
	logger.LogError(fmt.Sprintf("utente: %v", utente))
	urepo, err := h.repo.GetUtenteByName(utente.Nome)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ApiResponse{Code: http.StatusUnauthorized, Message: "Credenziali non valide (U001)", Data: ""})
		return
	}

	// verifica hash password fornita con quella registrata
	verified := security.VerifyPassword(urepo.Password, utente.Password)
	if !verified {
		c.JSON(http.StatusUnauthorized, ApiResponse{Code: http.StatusUnauthorized, Message: "Credenziali non valide (P001)", Data: ""})
		return
	}

	token, err := security.GeneraToken(urepo.Nome, urepo.Tipo_Ruolo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponse{Code: http.StatusInternalServerError, Message: "Errore nella generazione del token jwt", Data: ""})
		return
	}

	// ritorna il token jwt
	c.JSON(http.StatusOK, ApiResponse{Code: http.StatusOK, Message: "", Data: token})
}
