package system

import (
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"horgh-replicator/src/constants"
	"horgh-replicator/src/tools/exit"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	PositionMask = "%s:%s"
)

func SetPosition(hash string, position mysql.Position) error {
	content := fmt.Sprintf(PositionMask, position.Name, strconv.Itoa(int(position.Pos)))
	fileName := fmt.Sprintf(constants.PositionsPath, hash)
	err := ioutil.WriteFile(fileName, []byte(content), 0644)

	return err
}

func GetPosition(hash string) mysql.Position {
	fileName := fmt.Sprintf(constants.PositionsPath, hash)

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return mysql.Position{}
	}

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		exit.Fatal(constants.ErrorGetPosition, err.Error())
	}
	data := strings.Split(string(content[:]), ":")
	if len(data) != 2 {
		exit.Fatal(constants.ErrorGetPosition, "Can't parse position")
	}
	position, err := strconv.Atoi(data[1])
	if err != nil {
		exit.Fatal(constants.ErrorGetPosition, err.Error())
	}
	pos := mysql.Position{
		Name: data[0],
		Pos:  uint32(position),
	}

	return pos
}
