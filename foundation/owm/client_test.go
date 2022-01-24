package owm

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetHourlyWeather(t *testing.T) {
	var lat, lon = float64(42.65258), float64(-73.75623) //float64(51.50853), float64(-0.12574) // London coordinates
	client := Client{
		apiKey:   os.Getenv("OWM_TEST_API_KEY"),
		units:    UnitsMetric,
		language: LanguageEnglish,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	fcast, err := client.GetHourlyForecast(ctx, lat, lon)
	assert.NoError(t, err)
	assert.Greater(t, len(fcast), 0)

	clientErr := Client{
		apiKey:   "",
		units:    "",
		language: "",
	}

	fcerr, err := clientErr.GetHourlyForecast(ctx, lat, lon)
	assert.Error(t, err)
	assert.Nil(t, fcerr)
}
