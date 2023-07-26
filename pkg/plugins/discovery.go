package plugins

import (
	"log/slog"
	"path/filepath"

	"github.com/rancher/opni/pkg/logger"
	"github.com/rancher/opni/pkg/plugins/meta"
	"github.com/spf13/afero"
)

type Filter = func(meta.PluginMeta) bool

type DiscoveryConfig struct {
	// Directory to search for plugins
	Dir string
	// Optional filesystem (defaults to os filesystem)
	Fs afero.Fs
	// Optional filters to allow excluding plugins from discovery
	Filters []Filter
	// If true, query the plugin for its supported modes. The available modes
	// will be stored in the plugin's ExtendedMetadata and can be used in filters.
	QueryModes bool
	// Optional logger, defaults to no logging
	Logger *slog.Logger
}

func (dc DiscoveryConfig) Discover() []meta.PluginMeta {
	if dc.Fs == nil {
		dc.Fs = afero.NewOsFs()
	}
	paths, err := afero.Glob(dc.Fs, filepath.Join(dc.Dir, DefaultPluginGlob))
	if err != nil {
		panic(err)
	}
	var result []meta.PluginMeta
PLUGINS:
	for _, path := range paths {
		f, err := dc.Fs.Open(path)
		if err != nil {
			if dc.Logger != nil {
				dc.Logger.Error("failed to open plugin for reading", "plugin", path, logger.Err(err))
			}
			continue
		}
		md, err := meta.ReadFile(f)
		if err != nil {
			if dc.Logger != nil {
				dc.Logger.Error("failed to read plugin metadata", "plugin", path, logger.Err(err))
			}
			f.Close()
			continue
		}
		f.Close()
		if dc.QueryModes {
			if dc.Fs != nil {
				if _, ok := dc.Fs.(*afero.OsFs); !ok {
					panic("cannot query plugin modes with custom filesystem")
				}
			}
			modes, err := meta.QueryPluginModes(md.BinaryPath)
			if err != nil {
				if dc.Logger != nil {
					dc.Logger.Error("failed to query plugin modes", "plugin", path, logger.Err(err))
				}
				continue
			}
			if md.ExtendedMetadata == nil {
				md.ExtendedMetadata = &meta.ExtendedPluginMeta{}
			}
			md.ExtendedMetadata.ModeList = modes
		}

		for i, filter := range dc.Filters {
			if !filter(md) {
				if dc.Logger != nil {
					dc.Logger.Debug("plugin ignored due to filter", "plugin", path, "filter", i)
				}
				continue PLUGINS
			}
		}
		result = append(result, md)
	}
	return result
}
