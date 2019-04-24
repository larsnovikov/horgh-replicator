package plugins

import (
	"fmt"
	"github.com/siddontang/go-log/log"
	"horgh-replicator/src/connectors"
	"horgh-replicator/src/constants"
	"horgh-replicator/src/helpers"
	"plugin"
)

func Handle(beforeSave connectors.ConfigBeforeSave, value interface{}) interface{} {
	mod := fmt.Sprintf(constants.PluginPath, beforeSave.Handler)
	plug, err := plugin.Open(mod)
	if err != nil {
		log.Fatalf(constants.ErrorCachePluginError, err)
	}

	symHandler, err := plug.Lookup("Handler")
	if err != nil {
		log.Fatalf(constants.ErrorCachePluginError, err)
	}

	var handler helpers.Handler
	handler, ok := symHandler.(helpers.Handler)
	if !ok {
		log.Fatalf(constants.ErrorCachePluginError, "unexpected type from module symbol")
	}

	return handler.Handle(value, beforeSave.Params)
}
