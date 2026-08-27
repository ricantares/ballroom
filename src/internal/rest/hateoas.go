package rest

/*
	Facade delle funzionalità del package go2hal
*/
import (
	"github.com/pmoule/go2hal/hal"
)

type Link struct {
	Rel  string
	Href string
}

type BaseResource struct {
	Links []Link
	Data  any
}

const (
	SELF     = "self"
	FIRST    = "first"
	LAST     = "last"
	PREVIOUS = "prev"
	NEXT     = "next"
)

// dati di paginazione in caso di esposizione di elementi multipli
type Page struct {
	Items  int  `json:"items,omitempty"`  // totale elementi nella pagina corrente
	Size   int  `json:"size,omitempty"`   // totale elementi trovati
	Number int  `json:"number,omitempty"` // numero pagina corrente
	Pages  int  `json:"pages,omitempty"`  // totale pagine riempite
	First  bool `json:"first,omitempty"`  // se prima pagina
	Last   bool `json:"last,omitempty"`   // se ultima pagina
}

type Paging struct {
	Page Page `json:"page,omitempty"`
}

type Hateoas interface {
	Resource() hal.Resource
	AddLink(Link) error
	AddLinks([]Link) error
	AddData(any)
	AddEmbedded(string, []BaseResource)
	AddPage(Page)
	ToJSON() ([]byte, error)
}

type hateoasObject struct {
	origin   string
	resource hal.Resource
	encoder  hal.Encoder
}

// nuova risorsa in cui il parametro origin
// viene utilizzato come suffisso per la composizione degli "href"
func NewHateoas(origin string) Hateoas {
	obj := hateoasObject{
		origin:   origin,
		resource: hal.NewResourceObject(),
		encoder:  hal.NewEncoder(),
	}

	return &obj
}

func (h *hateoasObject) Resource() hal.Resource {
	return h.resource
}

func (h *hateoasObject) AddLink(l Link) error {
	link := &hal.LinkObject{Href: h.origin + l.Href}

	relation, err := hal.NewLinkRelation(l.Rel)
	if err != nil {
		return err
	}

	relation.SetLink(link)
	h.resource.AddLink(relation)

	return nil
}

func (h *hateoasObject) AddLinks(links []Link) error {

	for _, link := range links {
		h.AddLink(link) // TODO gestire l'errore
	}

	return nil
}

func (h *hateoasObject) AddData(data any) {
	h.resource.AddData(data)
}

// crea la parte "_embedded" con la radice "name"
func (h *hateoasObject) AddEmbedded(name string, e []BaseResource) {
	var embedded []hal.Resource
	for _, res := range e {
		he := NewHateoas(h.origin)
		he.AddLinks(res.Links)
		he.AddData(res.Data)
		embedded = append(embedded, he.Resource())
	}

	n, _ := hal.NewResourceRelation(name)
	n.SetResources(embedded)
	h.addResource(n)

}

// crea i parametri di paginazione
func (h *hateoasObject) AddPage(p Page) {
	h.resource.AddData(Paging{Page: p})
}

func (h *hateoasObject) ToJSON() ([]byte, error) {
	bytes, err := h.encoder.ToJSON(h.resource)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (h *hateoasObject) addResource(res hal.ResourceRelation) {
	h.resource.AddResource(res)
}
