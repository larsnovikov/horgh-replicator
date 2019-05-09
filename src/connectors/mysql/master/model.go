package master

import (
	"horgh-replicator/src/models/slave"
	"horgh-replicator/src/tools/exit"
)

type Model struct {
}

func (model *Model) Listen() {
	exit.BeforeExitPool = append(exit.BeforeExitPool, Stop)
	exit.BeforeExitPool = append(exit.BeforeExitPool, slave.Stop)
	BinlogListener()
}

func (model *Model) BuildSlave(table string) {
	buildModel(table)
}
