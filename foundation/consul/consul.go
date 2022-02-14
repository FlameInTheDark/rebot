package consul

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
)

// Discovery is a consul discovery client
type Discovery struct {
	client *api.Client
}

// NewConsulClient create new Consul consul client
func NewConsulClient(address string) (*Discovery, error) {
	cfg := api.DefaultConfig()
	cfg.Address = address

	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Discovery{client: client}, nil
}

// Register new service
func (d *Discovery) Register(id, name string, host string, port int, meta map[string]string) error {
	register := new(api.AgentServiceRegistration)
	register.Name = name
	register.ID = id

	if host == "" {
		var err error
		host, err = os.Hostname()
		if err != nil {
			return err
		}
	}
	register.Address = host

	register.Meta = meta

	register.Port = port
	register.Check = new(api.AgentServiceCheck)
	register.Check.DeregisterCriticalServiceAfter = "10s"
	register.Check.HTTP = fmt.Sprintf("http://%s:%d/healthz", host, port)
	register.Check.Interval = "5s"
	register.Check.Timeout = "3s"

	return d.client.Agent().ServiceRegister(register)
}

// Deregister service by ID
func (d *Discovery) Deregister(id string) error {
	return d.client.Agent().ServiceDeregister(id)
}

// DiscoverByName discover service by name
func (d *Discovery) DiscoverByName(name string) ([]Service, error) {
	return d.DiscoverFilter(fmt.Sprintf("Service == \"%s\"", name))
}

// DiscoverByTag discover services by tag
func (d *Discovery) DiscoverByTag(tag string) ([]Service, error) {
	return d.DiscoverFilter(fmt.Sprintf("Tags contains \"%s\"", tag))
}

// DiscoverByMeta discover services by specified meta
func (d *Discovery) DiscoverByMeta(key, value string) ([]Service, error) {
	return d.DiscoverFilter(fmt.Sprintf("Meta.%s == \"%s\"", key, value))
}

// DiscoverFilter discover service by specified filter
func (d *Discovery) DiscoverFilter(filter string) ([]Service, error) {
	data, err := d.client.Agent().ServicesWithFilter(filter)
	if err != nil {
		return nil, err
	}
	var services = make([]Service, 0, len(data))
	for _, s := range data {
		id, err := uuid.Parse(s.ID)
		if err != nil {
			continue
		}
		services = append(services, Service{
			ID:   id,
			Name: s.Service,
			Tags: s.Tags,
			Meta: s.Meta,
		})
	}
	return services, nil
}
