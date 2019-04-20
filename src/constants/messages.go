package constants

const (
	MessagePosFrom            = "Get pos from %s. Pos: %s; Name: %s"
	MessageDeleted            = "[time: %v][model: %s][pos: %v] delete row"
	MessageInserted           = "[time: %v][model: %s][pos: %v] insert row"
	MessageUpdated            = "[time: %v][model: %s][pos: %v] update row"
	MessageIgnoreInsert       = "[time: %v][model: %s][pos: %v] ignore insert row"
	MessageIgnoreDelete       = "[time: %v][model: %s][pos: %v] ignore delete row"
	MessageIgnoreUpdate       = "[time: %v][model: %s][pos: %v] ignore update row"
	MessageRetryConnect       = "Retry to connect to \"%s\" after %s seconds..."
	MessageLogFileChanged     = "Log file changed [model: %s] to \"%s\""
	MessageStopHandlingBinlog = "Stopping binlog handling..."
	MessageStopHandlingSave   = "Stopping replication handling..."
	MessageWait               = "Waiting %s %s..."
)
