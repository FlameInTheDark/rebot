package wgen

import (
	"bytes"
	"fmt"
	"image/png"
	"time"

	"github.com/fogleman/gg"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
)

var (
	errNilContext = errors.New("graphic context is nil")
)

// Generator is a weather picture generator
type Generator struct {
	iconsFile string
	fontFile  string

	iconsCache map[string]font.Face
	fontCache  map[string]font.Face
	bindings   *UnicodeBindings
}

// ForecastData contains forecast information
type ForecastData struct {
	Location string
	Forecast []ForecastRow
}

// ForecastRow forecast row
type ForecastRow struct {
	IconCode    string
	Temperature float64
	Humidity    int
	Clouds      int
	Max         float64
	Min         float64
	Time        time.Time
}

// NewGenerator creates a new generator
func NewGenerator(fontPath, iconsPath, bindings string) (*Generator, error) {
	binds, err := LoadBindings(bindings)
	if err != nil {
		return nil, err
	}

	return &Generator{
		iconsFile:  iconsPath,
		fontFile:   fontPath,
		bindings:   binds,
		fontCache:  make(map[string]font.Face),
		iconsCache: make(map[string]font.Face),
	}, nil
}

func (g *Generator) loadFont(dc *gg.Context, points float64) error {
	if _, ok := g.fontCache[fmt.Sprintf("%f", points)]; !ok {
		ff, err := gg.LoadFontFace(g.fontFile, points)
		if err != nil {
			return err
		}
		g.fontCache[fmt.Sprintf("%f", points)] = ff
	}
	dc.SetFontFace(g.fontCache[fmt.Sprintf("%f", points)])
	return nil
}

func (g *Generator) loadIcons(dc *gg.Context, points float64) error {
	if _, ok := g.iconsCache[fmt.Sprintf("%f", points)]; !ok {
		ff, err := gg.LoadFontFace(g.iconsFile, points)
		if err != nil {
			return err
		}
		g.iconsCache[fmt.Sprintf("%f", points)] = ff
	}
	dc.SetFontFace(g.iconsCache[fmt.Sprintf("%f", points)])
	return nil
}

// Generate creates a weather forecast picture
//TODO: Hardcoded generator isn't good. Need to create a template engine. It will also allow users to modify the template
func (g *Generator) Generate(data *ForecastData) (*bytes.Buffer, error) {
	dc := gg.NewContext(400, 650)

	err := g.prepareImage(dc)
	if err != nil {
		return nil, err
	}
	err = g.drawWeatherLines(dc)
	if err != nil {
		return nil, err
	}
	err = g.drawHeader(dc, data.Location, data.Forecast[0])
	if err != nil {
		return nil, err
	}
	err = g.drawTime(dc, data.Forecast[1:])
	if err != nil {
		return nil, err
	}
	err = g.drawClouds(dc, data.Forecast[1:])
	if err != nil {
		return nil, err
	}
	err = g.drawHumidity(dc, data.Forecast[1:])
	if err != nil {
		return nil, err
	}
	err = g.drawWeather(dc, data.Forecast[1:])
	if err != nil {
		return nil, err
	}
	err = g.drawIcons(dc, data.Forecast[1:])
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, dc.Image())
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (g *Generator) prepareImage(dc *gg.Context) error {
	if dc == nil {
		return errNilContext
	}

	dc.SetRGBA(0, 0, 0, 0)
	dc.Clear()

	// Template
	dc.SetRGB255(242, 97, 73)
	dc.DrawRoundedRectangle(0, 0, 400, 650, 10)
	dc.Fill()
	return nil
}

func (g *Generator) drawWeatherLines(dc *gg.Context) error {
	if dc == nil {
		return errNilContext
	}

	dc.SetRGB255(234, 89, 65)
	dc.DrawRectangle(0, 250, 400, 100)
	dc.DrawRectangle(0, 450, 400, 100)
	dc.Fill()

	dc.SetLineWidth(2)
	dc.SetRGBA(0, 0, 0, 0.05)
	dc.DrawLine(0, 250, 400, 250)
	dc.DrawLine(0, 349, 400, 348)
	dc.DrawLine(0, 450, 400, 450)
	dc.DrawLine(0, 549, 400, 548)
	dc.Stroke()
	return nil
}

func (g *Generator) drawHeader(dc *gg.Context, location string, forecast ForecastRow) error {
	if dc == nil {
		return errNilContext
	}
	// Header (place and date)
	err := g.loadFont(dc, 20)
	if err != nil {
		return err
	}

	dc.SetRGBA(1, 1, 1, 0.7)
	dc.DrawStringAnchored(location, 10, 15, 0, 0.5)
	dc.SetRGBA(1, 1, 1, 0.4)
	dc.DrawStringAnchored(time.Now().Format("Jan 2, 2006"), 270, 15, 0, 0.5)

	// First weather data
	dc.SetRGBA(1, 1, 1, 0.5)
	err = g.loadFont(dc, 30)
	if err != nil {
		return err
	}

	dc.DrawStringAnchored(fmt.Sprintf("%.2d:00", forecast.Time.Hour()), 50, 200, 0.5, 0.5)
	dc.DrawStringAnchored(fmt.Sprintf("H:%d%%", forecast.Humidity), 200, 200, 0.5, 0.5)
	dc.DrawStringAnchored(fmt.Sprintf("C:%d%%", forecast.Clouds), 350, 200, 0.5, 0.5)

	dc.SetRGBA(1, 1, 1, 1)
	err = g.loadFont(dc, 90)
	if err != nil {
		return err
	}

	dc.DrawStringAnchored(fmt.Sprintf("%d°", int(forecast.Temperature)), 100, 100, 0.5, 0.5)

	err = g.loadIcons(dc, 70)
	if err != nil {
		return err
	}

	dc.DrawStringAnchored(g.bindings.Get(forecast.IconCode), 250, 100, 0, 0.7)
	return nil
}

func (g *Generator) drawHeaderDaily(dc *gg.Context, location string, forecast ForecastRow) error {
	if dc == nil {
		return errNilContext
	}

	// Header (place and date)
	err := g.loadFont(dc, 20)
	if err != nil {
		return err
	}

	dc.SetRGBA(1, 1, 1, 0.7)
	dc.DrawStringAnchored(location, 10, 15, 0, 0.5)
	dc.SetRGBA(1, 1, 1, 0.4)
	dc.DrawStringAnchored(time.Now().Format("Jan 2, 2006"), 270, 15, 0, 0.5)

	// First weather data
	dc.SetRGBA(1, 1, 1, 0.5)
	err = g.loadFont(dc, 30)
	if err != nil {
		return err
	}

	dc.DrawStringAnchored(forecast.Time.Weekday().String(), 80, 200, 0.5, 0.5)
	dc.DrawStringAnchored(fmt.Sprintf("H:%d%%", forecast.Humidity), 200, 200, 0.5, 0.5)
	dc.DrawStringAnchored(fmt.Sprintf("C:%d%%", forecast.Clouds), 350, 200, 0.5, 0.5)

	dc.SetRGBA(1, 1, 1, 1)
	err = g.loadFont(dc, 90)
	if err != nil {
		return err
	}

	dc.DrawStringAnchored(fmt.Sprintf("%d°", int(forecast.Temperature)), 100, 100, 0.5, 0.5)

	err = g.loadIcons(dc, 70)
	if err != nil {
		return err
	}

	dc.DrawStringAnchored(g.bindings.Get(forecast.IconCode), 250, 100, 0, 0.7)
	return nil
}

func (g *Generator) drawTime(dc *gg.Context, data []ForecastRow) error {
	if dc == nil {
		return errNilContext
	}

	if len(data) < 4 {
		return errors.New("not enough data")
	}
	err := g.loadFont(dc, 30)
	if err != nil {
		return err
	}

	// Time
	dc.DrawStringAnchored(fmt.Sprintf("%.2v:00", data[0].Time.Hour()), 100, 285, 0, 0.5)
	dc.DrawStringAnchored(fmt.Sprintf("%.2v:00", data[1].Time.Hour()), 100, 385, 0, 0.5)
	dc.DrawStringAnchored(fmt.Sprintf("%.2v:00", data[2].Time.Hour()), 100, 485, 0, 0.5)
	dc.DrawStringAnchored(fmt.Sprintf("%.2v:00", data[3].Time.Hour()), 100, 585, 0, 0.5)

	return nil
}

func (g *Generator) drawDays(dc *gg.Context, data []ForecastRow) error {
	if dc == nil {
		return errNilContext
	}

	if len(data) < 4 {
		return errors.New("not enough data")
	}
	err := g.loadFont(dc, 30)
	if err != nil {
		return err
	}

	// Days
	dc.DrawStringAnchored(data[0].Time.Weekday().String(), 100, 285, 0, 0.5)
	dc.DrawStringAnchored(data[1].Time.Weekday().String(), 100, 385, 0, 0.5)
	dc.DrawStringAnchored(data[2].Time.Weekday().String(), 100, 485, 0, 0.5)
	dc.DrawStringAnchored(data[3].Time.Weekday().String(), 100, 585, 0, 0.5)

	return nil
}

func (g *Generator) drawHumidity(dc *gg.Context, data []ForecastRow) error {
	if dc == nil {
		return errNilContext
	}

	if len(data) < 4 {
		return errors.New("not enough data")
	}
	err := g.loadFont(dc, 20)
	if err != nil {
		return err
	}

	dc.SetRGBA(1, 1, 1, 0.5)

	dc.DrawStringAnchored(fmt.Sprintf("H:%d%%", data[0].Humidity), 100, 315, 0, 0.5)
	dc.DrawStringAnchored(fmt.Sprintf("H:%d%%", data[1].Humidity), 100, 415, 0, 0.5)
	dc.DrawStringAnchored(fmt.Sprintf("H:%d%%", data[2].Humidity), 100, 515, 0, 0.5)
	dc.DrawStringAnchored(fmt.Sprintf("H:%d%%", data[3].Humidity), 100, 615, 0, 0.5)
	return nil
}

func (g *Generator) drawClouds(dc *gg.Context, data []ForecastRow) error {
	if len(data) < 4 {
		return errors.New("not enough data")
	}
	err := g.loadFont(dc, 20)
	if err != nil {
		return err
	}

	dc.SetRGBA(1, 1, 1, 0.5)

	dc.DrawStringAnchored(fmt.Sprintf("C:%d%%", data[0].Clouds), 170, 315, 0, 0.5)
	dc.DrawStringAnchored(fmt.Sprintf("C:%d%%", data[1].Clouds), 170, 415, 0, 0.5)
	dc.DrawStringAnchored(fmt.Sprintf("C:%d%%", data[2].Clouds), 170, 515, 0, 0.5)
	dc.DrawStringAnchored(fmt.Sprintf("C:%d%%", data[3].Clouds), 170, 615, 0, 0.5)

	return nil
}

func (g *Generator) drawWeather(dc *gg.Context, data []ForecastRow) error {
	if len(data) < 4 {
		return errors.New("not enough data")
	}
	err := g.loadFont(dc, 50)
	if err != nil {
		return err
	}

	dc.SetRGBA(1, 1, 1, 1)

	dc.DrawStringAnchored(fmt.Sprintf("%d°", int(data[0].Temperature)), 320, 300, 0.5, 0.5)
	dc.DrawStringAnchored(fmt.Sprintf("%d°", int(data[1].Temperature)), 320, 400, 0.5, 0.5)
	dc.DrawStringAnchored(fmt.Sprintf("%d°", int(data[2].Temperature)), 320, 500, 0.5, 0.5)
	dc.DrawStringAnchored(fmt.Sprintf("%d°", int(data[3].Temperature)), 320, 600, 0.5, 0.5)
	return nil
}

func (g *Generator) drawWeatherDaily(dc *gg.Context, data []ForecastRow) error {
	if len(data) < 4 {
		return errors.New("not enough data")
	}
	err := g.loadFont(dc, 35)
	if err != nil {
		return err
	}

	dc.SetRGBA(1, 1, 1, 1)
	dc.DrawStringAnchored(fmt.Sprintf("%v°/%v°", int(data[0].Max), int(data[0].Min)), 320, 300, 0.5, 0.5)
	dc.DrawStringAnchored(fmt.Sprintf("%v°/%v°", int(data[1].Max), int(data[1].Min)), 320, 400, 0.5, 0.5)
	dc.DrawStringAnchored(fmt.Sprintf("%v°/%v°", int(data[2].Max), int(data[2].Min)), 320, 500, 0.5, 0.5)
	dc.DrawStringAnchored(fmt.Sprintf("%v°/%v°", int(data[3].Max), int(data[3].Min)), 320, 600, 0.5, 0.5)
	return nil
}

func (g *Generator) drawIcons(dc *gg.Context, data []ForecastRow) error {
	err := g.loadIcons(dc, 50)
	if err != nil {
		return err
	}

	dc.SetRGBA(1, 1, 1, 1)

	dc.DrawStringAnchored(g.bindings.Get(data[0].IconCode), 20, 280, 0, 0.7)
	dc.DrawStringAnchored(g.bindings.Get(data[1].IconCode), 20, 380, 0, 0.7)
	dc.DrawStringAnchored(g.bindings.Get(data[2].IconCode), 20, 480, 0, 0.7)
	dc.DrawStringAnchored(g.bindings.Get(data[3].IconCode), 20, 580, 0, 0.7)
	return nil
}

// GenerateDaily generates a weather image
func (g *Generator) GenerateDaily(data *ForecastData) (*bytes.Buffer, error) {
	dc := gg.NewContext(400, 650)

	err := g.prepareImage(dc)
	if err != nil {
		return nil, err
	}
	err = g.drawWeatherLines(dc)
	if err != nil {
		return nil, err
	}
	err = g.drawHeaderDaily(dc, data.Location, data.Forecast[0])
	if err != nil {
		return nil, err
	}
	err = g.drawDays(dc, data.Forecast[1:])
	if err != nil {
		return nil, err
	}
	err = g.drawClouds(dc, data.Forecast[1:])
	if err != nil {
		return nil, err
	}
	err = g.drawHumidity(dc, data.Forecast[1:])
	if err != nil {
		return nil, err
	}
	err = g.drawWeatherDaily(dc, data.Forecast[1:])
	if err != nil {
		return nil, err
	}
	err = g.drawIcons(dc, data.Forecast[1:])
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, dc.Image())
	if err != nil {
		return nil, err
	}
	return buf, nil
}
