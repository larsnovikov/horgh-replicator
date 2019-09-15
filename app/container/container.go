package container

import (
	"horgh2-replicator/app/configs"
	"horgh2-replicator/app/connectors"
	"horgh2-replicator/app/queue"
	"horgh2-replicator/app/replication"
)

type Container struct {
	Config           configs.Config
	MasterConnection connectors.MasterConnection
	SlaveConnection  connectors.SlaveConnection
	Replication      replication.Replication
	Queue            queue.Connection
}

func New(config configs.Config) (Container, error) {
	slave := connectors.NewSlave(config.Slave)
	master := connectors.NewMaster(config.Master)
	storage := replication.NewStorage(config.Master.Type, config.Master.Table)
	kafka, err := queue.New(config.Queue)

	if err != nil {
		return Container{}, err
	}
	return Container{
		Config:           config,
		MasterConnection: master,
		SlaveConnection:  slave,
		Replication:      replication.New(config.Replication, slave, storage),
		Queue:            kafka,
	}, nil
}
