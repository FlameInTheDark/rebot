package discovery

import "github.com/FlameInTheDark/rebot/foundation/consul"

// Discovery is a service discovery interface
type Discovery interface {
	Discover(name string) []consul.Service
}
