package main

type handling string

func (h handling) Handle(value interface{}, params []interface{}) interface{} {
	return params[0]
}

var Handler handling
