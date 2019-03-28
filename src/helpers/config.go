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
	Type          string
	RetryTimeout  int
	RetryAttempts int
}

type CredentialsDB struct {
	Credentials
	Host   string
	Port   int
	User   string
	Pass   string
	DBname string
}

type CredentialsAMQP struct {
	Credentials
	Host string
	Port string
	User string
	Pass string
}

var master CredentialsDB
var slave interface{}
var replicator CredentialsDB

var tables []string

func MakeCredentials() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var timeout, attempts int

	timeout, _ = strconv.Atoi(os.Getenv("MASTER_RETRY_TIMEOUT"))
	attempts, _ = strconv.Atoi(os.Getenv("MASTER_RETRY_ATTEMPTS"))

	masterPort, _ := strconv.Atoi(os.Getenv("MASTER_PORT"))
	master = CredentialsDB{
		Credentials{
			os.Getenv("MASTER_TYPE"),
			timeout,
			attempts,
		},
		os.Getenv("MASTER_HOST"),
		masterPort,
		os.Getenv("MASTER_USER"),
		os.Getenv("MASTER_PASS"),
		os.Getenv("MASTER_DBNAME"),
	}

	timeout, _ = strconv.Atoi(os.Getenv("SLAVE_RETRY_TIMEOUT"))
	attempts, _ = strconv.Atoi(os.Getenv("SLAVE_RETRY_ATTEMPTS"))
	slavePort, _ := strconv.Atoi(os.Getenv("SLAVE_PORT"))
	slave = CredentialsDB{
		Credentials{
			os.Getenv("SLAVE_TYPE"),
			timeout,
			attempts,
		},
		os.Getenv("SLAVE_HOST"),
		slavePort,
		os.Getenv("SLAVE_USER"),
		os.Getenv("SLAVE_PASS"),
		os.Getenv("SLAVE_DBNAME"),
	}

	timeout, _ = strconv.Atoi(os.Getenv("REPLICATOR_RETRY_TIMEOUT"))
	attempts, _ = strconv.Atoi(os.Getenv("REPLICATOR_RETRY_ATTEMPTS"))
	replicationPort, _ := strconv.Atoi(os.Getenv("REPLICATOR_PORT"))
	replicator = CredentialsDB{
		Credentials{
			"mysql",
			timeout,
			attempts,
		},
		os.Getenv("REPLICATOR_HOST"),
		replicationPort,
		os.Getenv("REPLICATOR_USER"),
		os.Getenv("REPLICATOR_PASS"),
		os.Getenv("REPLICATOR_DBNAME"),
	}

	tables = strings.Split(os.Getenv("ALLOWED_TABLES"), ",")
}

func GetCredentials(storageType string) interface{} {
	switch storageType {
	case constants.DBMaster:
		return getMaster()
	case constants.DBSlave:
		return getSlave()
	case constants.DBReplicator:
		return getReplicator()
	default:
		return Credentials{}
	}
}

func getMaster() CredentialsDB {
	return master
}

func getSlave() interface{} {
	return slave
}

func getReplicator() CredentialsDB {
	return replicator
}

func GetTables() []string {
	return tables
}
