package system

import (
	"bufio"
	"fmt"
	"github.com/siddontang/go-log/log"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/spf13/cobra"
	"horgh-replicator/src/constants"
	"horgh-replicator/src/helpers"
	"horgh-replicator/src/models/slave"
	"horgh-replicator/src/parser"
	"horgh-replicator/src/tools"
	helpers2 "horgh-replicator/src/tools/helpers"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	CreateDump     = "mysqldump --extended-insert=FALSE --no-create-info --master-data=1 --port=%s -u%s -p%s -h %s %s %s > %s"
	OutputFile     = "/tmp/dump.sql"
	InsertRegexp   = `VALUES \([A-Za-z0-9,\s,\S]+\)`
	PositionRegexp = `MASTER_LOG_FILE=\'([a-zA-Z\-\.0-9]+)\', MASTER_LOG_POS=([0-9]+)`
)

var CmdBuildTable = &cobra.Command{
	Use:   "build-table",
	Short: "Build slave table from master. Format: [table]",
	Long:  "Build slave table from master. Format: [table]",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tools.BeforeExit = func() bool {
			log.Infof(constants.MessageStopTableBuild)
			return false
		}
		tableName := args[0]
		buildModel(tableName)
	},
}

func buildModel(tableName string) {
	helpers2.Table = tableName
	if canHandle() == true {
		removeOutputFile()
		makeDump()
		parseDump()

		helpers2.Wait()
	} else {

	}
}

func canHandle() bool {
	// todo check if position exists and not 0
	savedPos := parser.GetSavedPos(helpers2.Table)
	if savedPos.Name == "" && savedPos.Pos == 0 {
		return true
	}

	log.Fatalf(constants.ErrorTableBuilt, helpers2.Table, helpers2.Table)
	return false
}

func removeOutputFile() {
	if _, err := os.Stat(OutputFile); os.IsExist(err) {
		err := os.Remove(OutputFile)
		if err != nil {
			log.Fatalf(constants.ErrorRemoveDump, OutputFile, err)
		}
	}
}

func makeDump() {
	log.Infof(constants.MessageStartCreateDump, helpers2.Table)
	cred := helpers.GetCredentials(constants.DBMaster).(helpers.CredentialsDB)

	dumpCmd := fmt.Sprintf(CreateDump, strconv.Itoa(cred.Port), cred.User, cred.Pass, cred.Host, cred.DBname, helpers2.Table, OutputFile)

	_, err := exec.Command("sh", "-c", dumpCmd).Output()
	if err != nil {
		log.Fatalf(constants.ErrorCreateDump, OutputFile, err)
	}

	log.Infof(constants.MessageDumpCreated, helpers2.Table)
}

func parseDump() {
	log.Infof(constants.MessageStartParseDump, helpers2.Table)
	file, err := os.Open(OutputFile)
	defer func() {
		err = file.Close()
	}()

	if err != nil {
		log.Fatalf(constants.ErrorParseDump, helpers2.Table)
	}

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			break
		}

		_ = parseLine(line)
	}
	log.Infof(constants.MessageDumpParsed, helpers2.Table)
}

func parseLine(line string) map[string]interface{} {
	out := make(map[string]interface{})

	re := regexp.MustCompile(InsertRegexp)
	match := re.FindStringSubmatch(line)
	if len(match) > 0 {
		// TODO fix me
		r := strings.NewReplacer("VALUES", "",
			"'", "",
			"(", "",
			")", "")

		params := strings.Split(strings.TrimSpace(r.Replace(match[0])), ",")

		slave.GetSlaveByName(helpers2.Table).ClearParams()

		interfaceParams := make([]interface{}, len(params))
		for i := range params {
			interfaceParams[i] = params[i]
		}
		err := parser.ParseRow(slave.GetSlaveByName(helpers2.Table), interfaceParams)
		if err != nil {
			log.Fatalf(constants.ErrorParseDump, OutputFile, err)
		}

		header, positionSet := helpers2.GetHeader()

		slave.GetSlaveByName(helpers2.Table).Insert(&header, positionSet)
		return out
	}

	// parse position
	re = regexp.MustCompile(PositionRegexp)
	match = re.FindStringSubmatch(line)

	if len(match) > 0 {
		pos, _ := strconv.Atoi(match[2])
		helpers2.Position = mysql.Position{
			Name: match[1],
			Pos:  uint32(pos),
		}

		helpers2.SetPosition()
	}

	return out
}
