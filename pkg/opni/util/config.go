package cliutil

import (
	"errors"
	"log/slog"

	"github.com/rancher/opni/pkg/config"
	"github.com/rancher/opni/pkg/config/meta"
)

func LoadConfigObjectsOrDie(
	configLocation string,
	lg *slog.Logger,
) meta.ObjectList {
	if configLocation == "" {
		// find config file
		path, err := config.FindConfig()
		if err != nil {
			if errors.Is(err, config.ErrConfigNotFound) {
				panic(`could not find a config file in current directory or ["/etc/opni"], and --config was not given`)
			}
			panic("an error occurred while searching for a config file")
		}
		lg.Info("using config file", "path", path)
		configLocation = path
	}
	objects, err := config.LoadObjectsFromFile(configLocation)
	if err != nil {
		panic("failed to load config")
	}
	return objects
}
