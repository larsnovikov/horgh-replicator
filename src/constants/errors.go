package constants

const (
	ErrorMysqlCanal     = "Invalid canal"
	ErrorMysqlPosition  = "Invalid position"
	ErrorSliceCreation  = "InterfaceSlice() given a non-slice type"
	ErrorNoColumn       = "There is no column %s in %s.%s"
	ErrorDBConnect      = "Can't connect to \"%s\""
	ErrorSave           = "Can't %s model in \"%s\" data: %v"
	ErrorExecQuery      = "Can't exec query. Type: \"%s\" error: \"%v\""
	ErrorNoModelFile    = "Model file \"%s\" not exists"
	ErrorParserPosition = "Catch error: \"%s\""
	ErrorGetMinPosition = "Can't get min position. Error: \"%s\""
	ErrorUndefinedSlave = "Can't get slave. Error: \"%s\""
	ErrorCobraStarter   = "Catch cobra error: \"%s\""
)
