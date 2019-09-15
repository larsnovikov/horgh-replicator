package configs

import (
	"errors"
	"github.com/joho/godotenv"
	"horgh2-replicator/app/constants"
	"os"
	"strconv"
)

type QueueConfig struct {
	Host        string
	SASLEnabled bool
	User        string
	Password    string
	Topic       string
}

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
	Queue       QueueConfig
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

	SASLEnabled, err := strconv.ParseBool(os.Getenv("QUEUE_SASL"))
	if err != nil {
		return config, err
	}
	config.Queue = QueueConfig{
		Host:        os.Getenv("QUEUE_HOST"),
		SASLEnabled: SASLEnabled,
		User:        os.Getenv("QUEUE_USER"),
		Password:    os.Getenv("QUEUE_PASSWORD"),
		Topic:       os.Getenv("QUEUE_TOPIC"),
	}

	return config, nil
}
