package helpers

type Handler interface {
	Handle(value interface{}, params []interface{}) interface{}
}
