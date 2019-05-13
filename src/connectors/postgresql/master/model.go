package master

import (
	"horgh-replicator/src/models/slave"
	"horgh-replicator/src/tools/exit"
)

type Model struct {
}

func (model *Model) Listen() {
	exit.BeforeExitPool = append(exit.BeforeExitPool, stop)
	exit.BeforeExitPool = append(exit.BeforeExitPool, slave.Stop)
	Listen()
}

func (model *Model) BuildSlave(table string) {
	buildModel(table)
}
