package wgen

import (
	"encoding/json"
	"os"
)

// UnicodeBindings is a map of unicode symbols.
// Key is a weather code, values is a font unicode symbol (weather icon from the font file)
type UnicodeBindings struct {
	binds map[string]string
}

// LoadBindings loads unicode bindings
func LoadBindings(path string) (*UnicodeBindings, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	var bindings UnicodeBindings

	err = json.NewDecoder(file).Decode(&bindings.binds)
	if err != nil {
		return nil, err
	}
	return &bindings, nil
}

// Get returns a unicode symbol by weather code. Returns an empty string if binding is not found
func (b *UnicodeBindings) Get(s string) string {
	if u, ok := b.binds[s]; ok {
		return u
	}
	return ""
}
