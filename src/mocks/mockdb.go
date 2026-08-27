package mock

import (
	"time"

	"ricantares.com/ballroom/internal/domain"
)

// Deve implementare tutti i metodi dell'interfaccia db.DB (cfr. repo.go)
type MockDb struct {
}

func (m *MockDb) GetScuola() (domain.Scuola, error) {
	return scuolaMock()
}
func (m *MockDb) GetSala(id domain.Uuid) (sala domain.Sala, err error) { return sala, nil }
func (m *MockDb) ListSala() (sale []domain.Sala, err error)            { return sale, nil }
func (m *MockDb) CreateSala(sala domain.Sala) (result domain.Sala, err error) {
	return *new(domain.Sala), nil
}
func (m *MockDb) UpdateSala(sala domain.Sala) (result domain.Sala, err error) {
	return *new(domain.Sala), nil
}
func (m *MockDb) DeleteSala(id domain.Uuid) (del domain.Uuid, err error) {
	return 0, nil
}

func scuolaMock() (domain.Scuola, error) {
	t := time.Now()
	m := domain.Model{
		Id:         1,
		Created_at: &t,
		Updated_at: &t,
	}
	s := domain.Scuola{
		Model: m,
		Nome:  "Scuola Mock",
	}

	return s, nil
}
