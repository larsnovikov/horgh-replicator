package helpers

type Storage interface {
	Ping() bool
	Exec(params map[string]interface{}) bool
}
