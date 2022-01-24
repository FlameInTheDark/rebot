package geonames

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient_FindOneLocation(t *testing.T) {
	var location = "north york"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	client := Client{username: os.Getenv("GEONAMES_TEST_NAME")}

	geo, err := client.FindOneLocation(ctx, location)

	assert.NoError(t, err)
	assert.NotNil(t, geo)

	client2 := Client{username: ""}

	geo2, err2 := client2.FindOneLocation(ctx, location)

	assert.Error(t, err2)
	assert.Nil(t, geo2)
}
