package app

import (
	"horgh2-replicator/app/configs"
	"horgh2-replicator/app/container"
)

func New() (container.Container, error) {
	config, err := configs.New()
	if err != nil {
		return container.Container{}, err
	}

	cont, err := container.New(config)
	if err != nil {
		return container.Container{}, err
	}

	return cont, err
}
