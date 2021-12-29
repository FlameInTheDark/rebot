package wgen

import (
	"encoding/json"
	"os"
)

type UnicodeBindings struct {
	binds map[string]string
}

func LoadBindings(path string) (*UnicodeBindings, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var bindings UnicodeBindings

	err = json.NewDecoder(file).Decode(&bindings.binds)
	if err != nil {
		return nil, err
	}
	return &bindings, nil
}

func (b *UnicodeBindings) Get(s string) string {
	if u, ok := b.binds[s]; ok {
		return u
	}
	return ""
}
