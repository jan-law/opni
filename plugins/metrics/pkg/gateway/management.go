package gateway

import (
	"github.com/rancher/opni/pkg/logger"
	"github.com/rancher/opni/plugins/metrics/pkg/gateway/drivers"
)

func (p *Plugin) configureCortexManagement() {
	driverName := p.config.Get().Spec.Cortex.Management.ClusterDriver
	if driverName == "" {
		p.logger.Warn("no cluster driver configured")
	}

	builder, ok := drivers.ClusterDrivers.Get(driverName)
	if !ok {
		p.logger.Error("unknown cluster driver, using fallback noop driver", "driver", driverName)

		builder, ok = drivers.ClusterDrivers.Get("noop")
		if !ok {
			panic("bug: noop cluster driver not found")
		}
	}

	driver, err := builder(p.ctx)
	if err != nil {
		p.logger.Error("failed to initialize cluster driver", logger.Err(err), "driver", driverName)
		return
	}

	p.clusterDriver.Set(driver)
}
