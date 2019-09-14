package replication

import (
	"horgh2-replicator/app/connectors/mysql"
	"horgh2-replicator/app/constants"
)

type Storage interface {
	SetPosition(interface{}) error
	GetPosition() (interface{}, error)
}

func NewStorage(masterType string, entity string) Storage {
	switch masterType {
	case constants.TypeMYSQL:
		return mysql.New(entity)
	}

	return nil
}
