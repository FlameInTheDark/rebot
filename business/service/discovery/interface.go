package discovery

import "github.com/FlameInTheDark/rebot/foundation/consul"

type Discovery interface {
	Discover(name string) []consul.Service
}
