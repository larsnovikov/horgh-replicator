package helpers

type Storage interface {
	Ping() bool
	Exec(params Query) bool
}

type QueryAction struct {
	Method     func() bool
	StopMethod func() bool
}

type Query struct {
	Query  string
	Params []interface{}
}
