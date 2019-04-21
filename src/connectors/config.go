package connectors

type ConfigSlave struct {
	Table  string        `json:"table"`
	Fields []ConfigField `json:"fields"`
}

type ConfigBeforeSave struct {
	Handler string        `json:"handler"`
	Params  []interface{} `json:"params"`
}

type ConfigField struct {
	Name       string           `json:"name"`
	Key        bool             `json:"key"`
	Mode       string           `json:"mode"`
	BeforeSave ConfigBeforeSave `json:"beforeSave"`
}
