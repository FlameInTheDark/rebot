package owm

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetHourlyWeather(t *testing.T) {
	var lat, lon = float64(51.50853), float64(-0.12574) // London coordinates
	client := Client{
		apiKey:   os.Getenv("OWM_TEST_API_KEY"),
		units:    UnitsMetric,
		language: LanguageEnglish,
	}

	fcast, err := client.GetHourlyForecast(lat, lon)
	assert.NoError(t, err)
	assert.Greater(t, len(fcast), 0)

	clientErr := Client{
		apiKey:   "",
		units:    "",
		language: "",
	}

	fcerr, err := clientErr.GetHourlyForecast(lat, lon)
	assert.Error(t, err)
	assert.Nil(t, fcerr)
}
