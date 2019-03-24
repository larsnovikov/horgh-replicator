package models

type AbstractModel interface {
	TableName() string
	SchemaName() string
	Insert() bool
	Update() bool
	Delete() bool
	ParseKey([]interface{})
}

func GetModel(name string) interface{ AbstractModel } {
	var model func() interface{ AbstractModel }
	switch name {
	case "user":
		model = func() interface{ AbstractModel } {
			return &User{}
		}
	case "post":
		model = func() interface{ AbstractModel } {
			return &Post{}
		}
	}

	output := model()

	return output
}
