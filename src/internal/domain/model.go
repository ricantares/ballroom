package domain

import (
	"strconv"
	"time"
)

type Uuid int32

// Base model dal quale ereditano tutte le entita'
type Model struct {
	Id         Uuid       `db:"id" json:"Id"`
	Created_at *time.Time `db:"created_at" json:"Created_at"`
	Updated_at *time.Time `db:"updated_at" json:"Updated_at"`
	Deleted    bool       `db:"deleted" json:"Deleted"`
}

// Converte una stringa contenente un Uuid in Uuid
func ToUuid(s string) (id Uuid, err error) {
	ii, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	id = Uuid(ii)

	return id, nil

}
