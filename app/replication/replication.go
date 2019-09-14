package replication

import (
	"fmt"
	"horgh2-replicator/app/configs"
	"horgh2-replicator/app/connectors"
)

type Replication struct {
	CurrentLogFile     string
	CurrentLogPosition int
	Storage            Storage
	QueryChannel       chan Query
	Connection         connectors.SlaveConnection
}

func New(config configs.ReplicationConfig, connection connectors.SlaveConnection, storage Storage) Replication {
	// TODO calc channel size
	qc := make(chan Query, 9999)

	return Replication{
		QueryChannel: qc,
		Connection:   connection,
		Storage:      storage,
	}
}

func (r Replication) HandleQueries() {
	query := <-r.QueryChannel

	fmt.Println(query)
}
