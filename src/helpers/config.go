package helpers

import (
	"github.com/joho/godotenv"
	"github.com/siddontang/go-log/log"
	"go-binlog-replication/src/constants"
	"os"
	"strconv"
	"strings"
)

type Credentials struct {
	Host          string
	Port          int
	User          string
	Pass          string
	DBname        string
	Type          string
	RetryTimeout  int
	RetryAttempts int
}

type CredentialsPool struct {
	master     Credentials
	slave      Credentials
	replicator Credentials
	tables     []string
}

var credentials CredentialsPool

func MakeCredentials() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var timeout, attempts int

	timeout, _ = strconv.Atoi(os.Getenv("MASTER_RETRY_TIMEOUT"))
	attempts, _ = strconv.Atoi(os.Getenv("MASTER_RETRY_ATTEMPTS"))

	masterPort, _ := strconv.Atoi(os.Getenv("MASTER_PORT"))
	masterCredentials := Credentials{
		os.Getenv("MASTER_HOST"),
		masterPort,
		os.Getenv("MASTER_USER"),
		os.Getenv("MASTER_PASS"),
		os.Getenv("MASTER_DBNAME"),
		os.Getenv("MASTER_TYPE"),
		timeout,
		attempts,
	}

	timeout, _ = strconv.Atoi(os.Getenv("SLAVE_RETRY_TIMEOUT"))
	attempts, _ = strconv.Atoi(os.Getenv("SLAVE_RETRY_ATTEMPTS"))
	slavePort, _ := strconv.Atoi(os.Getenv("SLAVE_PORT"))
	slaveCredentials := Credentials{
		os.Getenv("SLAVE_HOST"),
		slavePort,
		os.Getenv("SLAVE_USER"),
		os.Getenv("SLAVE_PASS"),
		os.Getenv("SLAVE_DBNAME"),
		os.Getenv("SLAVE_TYPE"),
		timeout,
		attempts,
	}

	timeout, _ = strconv.Atoi(os.Getenv("REPLICATOR_RETRY_TIMEOUT"))
	attempts, _ = strconv.Atoi(os.Getenv("REPLICATOR_RETRY_ATTEMPTS"))
	replicationPort, _ := strconv.Atoi(os.Getenv("REPLICATOR_PORT"))
	replicatorCredentials := Credentials{
		os.Getenv("REPLICATOR_HOST"),
		replicationPort,
		os.Getenv("REPLICATOR_USER"),
		os.Getenv("REPLICATOR_PASS"),
		os.Getenv("REPLICATOR_DBNAME"),
		"mysql",
		timeout,
		attempts,
	}

	credentials = CredentialsPool{
		master:     masterCredentials,
		slave:      slaveCredentials,
		replicator: replicatorCredentials,
		tables:     strings.Split(os.Getenv("ALLOWED_TABLES"), ","),
	}
}

func GetCredentials(dbName string) Credentials {
	switch db := dbName; db {
	case constants.DBMaster:
		return credentials.master
	case constants.DBSlave:
		return credentials.slave
	case constants.DBReplicator:
		return credentials.replicator
	default:
		return Credentials{}
	}
}

func GetTables() []string {
	return credentials.tables
}
