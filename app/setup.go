package app

import (
	"horgh2-replicator/app/configs"
	"horgh2-replicator/app/container"
)

func Make() error {
	config, err := configs.New()
	if err != nil {
		return err
	}
	container.Make(config)

	return nil
}
