package consul

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

// CommandMeta is a command service meta data
type CommandMeta struct {
	Command string `json:"command" validate:"required"`
	Queue   string `json:"queue" validate:"required"`
}

// CommandMetaInfo contains meta information about a commands
type CommandMetaInfo []CommandMeta

// ParseCommandMeta returns an array of metadata converted from JSON. If it's invalid, returns error
func ParseCommandMeta(data []byte) (*CommandMetaInfo, error) {
	var meta CommandMetaInfo
	err := json.Unmarshal(data, &meta)
	if err != nil {
		return nil, err
	}
	v := validator.New()
	for _, s := range meta {
		err = v.Struct(s)
		if err != nil {
			return nil, err
		}
	}
	return &meta, nil
}

// MarshalCommandMeta converts a command metadata to the JSON format
func MarshalCommandMeta(data []CommandMeta) (string, error) {
	v := validator.New()
	for _, s := range data {
		err := v.Struct(s)
		if err != nil {
			return "", err
		}
	}

	meta, err := json.Marshal(CommandMetaInfo(data))
	return string(meta), err
}
