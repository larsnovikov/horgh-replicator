package helpers

type Storage interface {
	Ping() bool
	Exec(params map[string]interface{}) bool
}

type QueryAction struct {
	Method     func() bool
	StopMethod func() bool
}
