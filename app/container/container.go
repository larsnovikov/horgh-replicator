package container

import (
	"horgh2-replicator/app/configs"
	"horgh2-replicator/app/connectors"
	"horgh2-replicator/app/replication"
)

var BaseContainer Container

type Container struct {
	Config           configs.Config
	MasterConnection connectors.MasterConnection
	SlaveConnection  connectors.SlaveConnection
	Replication      replication.Replication
}

func Make(config configs.Config) {
	slave := connectors.NewSlave(config.Slave)
	master := connectors.NewMaster(config.Master)
	storage := replication.NewStorage(config.Master.Type, config.Master.Table)

	BaseContainer = Container{
		Config:           config,
		MasterConnection: master,
		SlaveConnection:  slave,
		Replication:      replication.New(config.Replication, slave, storage),
	}
}
