package master

import (
	"bufio"
	"fmt"
	"github.com/siddontang/go-log/log"
	"github.com/siddontang/go-mysql/mysql"
	"horgh-replicator/src/constants"
	"horgh-replicator/src/helpers"
	"horgh-replicator/src/models/slave"
	"horgh-replicator/src/tools/exit"
	toolsHelper "horgh-replicator/src/tools/helpers"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	CreateDump      = "mysqldump --extended-insert=FALSE --no-create-info --master-data=1 --port=%s -u%s -p%s -h %s %s %s"
	InsertRegexp    = `VALUES \([A-Za-z0-9,\s,\S]+\)`
	PositionRegexp  = `MASTER_LOG_FILE=\'([a-zA-Z\-\.0-9]+)\', MASTER_LOG_POS=([0-9]+)`
	ParseStringSize = 99999999
)

func buildModel(tableName string) {
	toolsHelper.Table = tableName
	if canHandle() == true {
		toolsHelper.ParseStrings = make(chan string, ParseStringSize)
		go parseLine(toolsHelper.ParseStrings)

		readDump()

		toolsHelper.Wait(func() bool {
			return slave.GetSlaveByName(toolsHelper.Table).GetChannelLen() == 0 && len(toolsHelper.ParseStrings) == 0
		})
	}
}

func canHandle() bool {
	savedPos := GetSavedPos(toolsHelper.Table)
	if savedPos.Name == "" && savedPos.Pos == 0 {
		return true
	}

	exit.Fatal(constants.ErrorSlaveBuilt, toolsHelper.Table, toolsHelper.Table)
	return false
}

func readDump() {
	log.Infof(constants.MessageStartReadDump, toolsHelper.Table)
	cred := helpers.GetCredentials(constants.DBMaster).(helpers.CredentialsDB)

	dumpCmd := fmt.Sprintf(CreateDump, strconv.Itoa(cred.Port), cred.User, cred.Pass, cred.Host, cred.DBname, toolsHelper.Table)
	cmdArgs := strings.Fields(dumpCmd)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:len(cmdArgs)]...)
	// create a pipe for the output of the script
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		exit.Fatal(constants.ErrorDumpRead, err)
	}

	scanner := bufio.NewScanner(cmdReader)

	go func() {
		for scanner.Scan() {
			toolsHelper.ParseStrings <- scanner.Text()
		}
	}()

	err = cmd.Start()
	if err != nil {
		exit.Fatal(constants.ErrorDumpRead, err)
	}

	log.Infof(constants.MessageDumpRead, toolsHelper.Table)
}

func parseLine(c chan string) {
	for {
		line := <-c

		// try to parse like insert
		if parseInsert(line) == true {
			continue
		}

		// try to parse like position setter
		if parsePosition(line) == true {
			continue
		}
	}
}

func parseInsert(line string) bool {
	re := regexp.MustCompile(InsertRegexp)
	match := re.FindStringSubmatch(line)
	if len(match) > 0 {
		// TODO fix me
		r := strings.NewReplacer("VALUES", "",
			"'", "",
			"(", "",
			")", "")

		params := strings.Split(strings.TrimSpace(r.Replace(match[0])), ",")

		slave.GetSlaveByName(toolsHelper.Table).ClearParams()

		interfaceParams := make([]interface{}, len(params))
		for i := range params {
			interfaceParams[i] = params[i]
		}
		err := ParseRow(slave.GetSlaveByName(toolsHelper.Table), interfaceParams)
		if err != nil {
			exit.Fatal(constants.ErrorParseLine, line, err)
		}

		header, _ := toolsHelper.GetHeader()

		slave.GetSlaveByName(toolsHelper.Table).Insert(&header)

		return true
	}

	return false
}

func parsePosition(line string) bool {
	re := regexp.MustCompile(PositionRegexp)
	match := re.FindStringSubmatch(line)

	if len(match) > 0 {
		pos, _ := strconv.Atoi(match[2])
		toolsHelper.Position = mysql.Position{
			Name: match[1],
			Pos:  uint32(pos),
		}

		toolsHelper.SetPosition()

		return true
	}

	return false
}
