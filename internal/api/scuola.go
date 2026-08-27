package api

/*
Handler scuola
*/
import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"ricantares.com/ballroom/internal/domain"
	"ricantares.com/ballroom/internal/rest"
)

type RouterScuola interface {
	GetScuola(*gin.Context)
	GetSala(*gin.Context) (domain.Sala, error)
	ListSala(*gin.Context) ([]domain.Sala, error)
	CreateSala(*gin.Context) (result domain.Sala, err error)
	UpdateSala(*gin.Context) (domain.Sala, error)
	DeleteSala(*gin.Context) (domain.Uuid, error)
}

func (h RouteHandler) GetScuola(c *gin.Context) {
	scuola, err := h.repo.GetScuola()
	if err != nil {
		c.JSON(http.StatusNotFound, ApiResponse{Code: http.StatusNotFound, Message: err.Error(), Data: ""})
		return
	}

	res := rest.BaseResource{Links: []rest.Link{{Rel: rest.SELF, Href: "/scuola"}}, Data: scuola}
	jsonData, _ := SimpleResponse(res) //TODO gestire l'errore

	c.JSON(http.StatusOK, ApiResponse{Code: http.StatusOK, Message: "", Data: jsonData})
}

func (h RouteHandler) GetSala(c *gin.Context) {
	id, err := domain.ToUuid(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ApiResponse{Code: http.StatusBadRequest, Message: err.Error(), Data: ""})
		return
	}
	sala, err := h.repo.GetSala(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	res := rest.BaseResource{Links: []rest.Link{{Rel: rest.SELF, Href: fmt.Sprintf("/scuola/sala/%v", id)}}, Data: sala}
	jsonData, _ := SimpleResponse(res) //TODO gestire l'errore

	c.JSON(http.StatusOK, ApiResponse{Code: http.StatusOK, Message: "", Data: jsonData})
}

func (h RouteHandler) GetSale(c *gin.Context) {
	sale, err := h.repo.ListSala()
	if err != nil {
		c.JSON(http.StatusNotFound, ApiResponse{Code: http.StatusNotFound, Message: err.Error(), Data: ""})
		return
	}

	res := ResourceWithEmbedds{}
	res.Base = rest.BaseResource{Links: []rest.Link{{Rel: rest.SELF, Href: "/scuola/sale"}}, Data: ""}
	emb := Embedded{}
	emb.Name = "sale"
	for _, sala := range sale {
		embedded := rest.BaseResource{Links: []rest.Link{{Rel: rest.SELF, Href: fmt.Sprintf("/scuola/sala/%v", sala.Id)}}, Data: sala}
		emb.Data = append(emb.Data, embedded)
	}
	res.Embedds = emb
	jsonData, _ := ResponseWithEmbedds(res) //TODO gestire l'errore

	c.JSON(http.StatusOK, ApiResponse{Code: http.StatusOK, Message: "", Data: jsonData})
}

func (h RouteHandler) PostSala(c *gin.Context) {
	// Get request body and convert it to domain.Sala
	var sala domain.Sala
	if err := c.ShouldBindJSON(&sala); err != nil {
		c.JSON(http.StatusBadRequest, ApiResponse{Code: http.StatusBadRequest, Message: err.Error(), Data: ""})
		return
	}

	// Add to the store
	newsala, err := h.repo.CreateSala(sala)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponse{Code: http.StatusInternalServerError, Message: err.Error(), Data: ""})
		return
	}

	// Return success payload
	c.JSON(http.StatusOK, ApiResponse{Code: http.StatusCreated, Message: "", Data: newsala})
}

func (h RouteHandler) PutSala(c *gin.Context) {
	// Get request body and convert it to domain.Sala
	var sala domain.Sala
	if err := c.ShouldBindJSON(&sala); err != nil {
		c.JSON(http.StatusBadRequest, ApiResponse{Code: http.StatusBadRequest, Message: err.Error(), Data: ""})
		return
	}

	// Modify the store
	updsala, err := h.repo.UpdateSala(sala)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponse{Code: http.StatusInternalServerError, Message: err.Error(), Data: ""})
		return
	}

	// Return success payload
	c.JSON(http.StatusOK, ApiResponse{Code: http.StatusOK, Message: "Updated", Data: updsala})
}

func (h RouteHandler) DeleteSala(c *gin.Context) {
	id, err := domain.ToUuid(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ApiResponse{Code: http.StatusBadRequest, Message: err.Error(), Data: ""})
		return
	}

	// Delete the store
	delid, err := h.repo.DeleteSala(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ApiResponse{Code: http.StatusNotFound, Message: err.Error(), Data: ""})
		return
	}

	// Return success payload
	type deleted struct {
		Del domain.Uuid `json:"Id"`
	}
	c.JSON(http.StatusOK, ApiResponse{Code: http.StatusOK, Message: "Deleted", Data: deleted{Del: delid}})
}
