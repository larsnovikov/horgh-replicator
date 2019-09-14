package mysql

import (
	"errors"
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"horgh2-replicator/app/constants"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	PositionMask = "%s:%s"
)

type Storage struct {
	entity string
}

func New(entity string) Storage {
	return Storage{entity: entity}
}

func (s Storage) SetPosition(position interface{}) error {
	content := fmt.Sprintf(PositionMask, position.(mysql.Position).Name, strconv.Itoa(int(position.(mysql.Position).Pos)))
	fileName := fmt.Sprintf(constants.PositionsPath, s.entity)
	err := ioutil.WriteFile(fileName, []byte(content), 0644)

	return err
}

func (s Storage) GetPosition() (interface{}, error) {
	fileName := fmt.Sprintf(constants.PositionsPath, s.entity)

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return mysql.Position{}, err
	}

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return mysql.Position{}, err
	}
	data := strings.Split(string(content[:]), ":")
	if len(data) != 2 {
		return mysql.Position{}, errors.New(constants.ErrorParsePosition)
	}
	position, err := strconv.Atoi(data[1])
	if err != nil {
		return mysql.Position{}, err
	}
	pos := mysql.Position{
		Name: data[0],
		Pos:  uint32(position),
	}

	return pos, nil
}
