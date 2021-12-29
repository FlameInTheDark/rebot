package discovery

import (
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/FlameInTheDark/rebot/foundation/consul"
)

type ConsulDiscovery struct {
	client   *consul.ConsulDiscovery
	handlers []HandlerFunc

	stop chan struct{}

	logger *zap.Logger
	mu     sync.RWMutex
}

type HandlerFunc func(discovery consul.Service)

func NewConsulDiscoveryService(client *consul.ConsulDiscovery, logger *zap.Logger) *ConsulDiscovery {
	return &ConsulDiscovery{client: client, logger: logger}
}

func (d *ConsulDiscovery) Discover(name string) []consul.Service {
	services, err := d.client.DiscoverByName(name)
	if err != nil {
		return nil
	}
	return services
}

func (d *ConsulDiscovery) AddHandler(f HandlerFunc) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers = append(d.handlers, f)
}

func (d *ConsulDiscovery) Start(name string) {
	t := time.Tick(time.Minute)
	for {
		select {
		case <-t:
			services, err := d.client.DiscoverByName(name)
			if err != nil {
				d.logger.Error("Commands discovery error", zap.Error(err))
				continue
			}
			d.logger.Debug("Found consul services", zap.Int("services-count", len(services)))
			for i, _ := range services {
				d.mu.RLock()
				for _, h := range d.handlers {
					h(services[i])
				}
				d.mu.RUnlock()
			}
		case <-d.stop:
			return
		}
	}
}

func (d *ConsulDiscovery) Stop() {
	close(d.stop)
}
