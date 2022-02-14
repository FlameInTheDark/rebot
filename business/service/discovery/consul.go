package discovery

import (
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/FlameInTheDark/rebot/foundation/consul"
)

// ConsulDiscovery is a discovery service
type ConsulDiscovery struct {
	client   *consul.Discovery
	handlers []HandlerFunc

	stop chan struct{}

	logger *zap.Logger
	mu     sync.RWMutex
}

// HandlerFunc discovery handler. Called when new service is discovered
type HandlerFunc func(discovery consul.Service)

// NewConsulDiscoveryService returns a new ConsulDiscoveryService with logger
func NewConsulDiscoveryService(client *consul.Discovery, logger *zap.Logger) *ConsulDiscovery {
	return &ConsulDiscovery{client: client, logger: logger}
}

// Discover discover services by name
func (d *ConsulDiscovery) Discover(name string) []consul.Service {
	services, err := d.client.DiscoverByName(name)
	if err != nil {
		return nil
	}
	return services
}

// AddHandler ...
func (d *ConsulDiscovery) AddHandler(f HandlerFunc) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers = append(d.handlers, f)
}

// Start starts a service discovery for the given service name
func (d *ConsulDiscovery) Start(name string) {
	d.discover(name) // first discover
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:
			d.discover(name)
		case <-d.stop:
			return
		}
	}
}

func (d *ConsulDiscovery) discover(name string) {
	services, err := d.client.DiscoverByName(name)
	if err != nil {
		d.logger.Error("Commands discovery error", zap.Error(err))
		return
	}
	d.logger.Debug("Found consul services", zap.Int("services-count", len(services)))
	for i := range services {
		d.mu.RLock()
		for _, h := range d.handlers {
			h(services[i])
		}
		d.mu.RUnlock()
	}
}

func (d *ConsulDiscovery) Stop() {
	close(d.stop)
}
