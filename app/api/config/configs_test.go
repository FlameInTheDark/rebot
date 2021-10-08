package config

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetConfig_GetTheSameInstance(t *testing.T) {
	conf, err := GetConfig()
	assert.NoError(t, err)

	conf.Database.Host = uuid.New().String()

	nextConf, err := GetConfig()
	assert.NoError(t, err)

	assert.Equal(t, conf.Database.Host, nextConf.Database.Host)
}
