package helpers

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

type Credentials struct {
	Host   string
	Port   int
	User   string
	Pass   string
	DBname string
	Type   string
}

var credentials CredentialsPool

type CredentialsPool struct {
	master Credentials
	slave  Credentials
	tables []string
}

func MakeCredentials() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	masterPort, _ := strconv.Atoi(os.Getenv("MASTER_PORT"))
	masterCredentials := Credentials{
		os.Getenv("MASTER_HOST"),
		masterPort,
		os.Getenv("MASTER_USER"),
		os.Getenv("MASTER_PASS"),
		os.Getenv("MASTER_DBNAME"),
		os.Getenv("MASTER_TYPE"),
	}

	slavePort, _ := strconv.Atoi(os.Getenv("SLAVE_PORT"))
	slaveCredentials := Credentials{
		os.Getenv("SLAVE_HOST"),
		slavePort,
		os.Getenv("SLAVE_USER"),
		os.Getenv("SLAVE_PASS"),
		os.Getenv("SLAVE_DBNAME"),
		os.Getenv("SLAVE_TYPE"),
	}

	credentials = CredentialsPool{
		master: masterCredentials,
		slave:  slaveCredentials,
		tables: strings.Split(os.Getenv("ALLOWED_TABLES"), ","),
	}
}

func GetCredentials(dbName string) Credentials {
	switch db := dbName; db {
	case "master":
		return credentials.master
	case "slave":
		return credentials.slave
	default:
		return Credentials{}
	}
}

func GetTables() []string {
	return credentials.tables
}
