package configs

import (
	"errors"
	"github.com/joho/godotenv"
	"horgh2-replicator/app/constants"
	"os"
)

type ConnectionConfig struct {
	DSN   string
	Type  string
	Table string
	Role  string
}

type ReplicationConfig struct {
	SlaveId       string
	LogFilePrefix string
}

type Config struct {
	Master      ConnectionConfig
	Slave       ConnectionConfig
	Replication ReplicationConfig
}

func New() (Config, error) {
	config := Config{}
	err := godotenv.Load()
	if err != nil {
		return config, errors.New(constants.ErrorLoadEnv)
	}

	config.Master = ConnectionConfig{
		DSN:   os.Getenv("MASTER_DSN"),
		Type:  os.Getenv("MASTER_TABLE"),
		Table: os.Getenv("MASTER_TABLE"),
		Role:  constants.RoleMaster,
	}

	config.Slave = ConnectionConfig{
		DSN:   os.Getenv("SLAVE_DSN"),
		Type:  os.Getenv("SLAVE_TYPE"),
		Table: os.Getenv("SLAVE_TABLE"),
		Role:  constants.RoleSlave,
	}

	config.Replication = ReplicationConfig{
		SlaveId:       os.Getenv("REPLICATION_SLAVE_ID"),
		LogFilePrefix: os.Getenv("REPLICATION_LOG_FILE_PREFIX"),
	}

	return config, nil
}
