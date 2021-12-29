package wgen

import (
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
	"image/png"
	"time"
)

// Generator is a weather picture generator
type Generator struct {
	iconsFile string
	fontFile  string

	iconsCache map[float64]font.Face
	fontCache  map[float64]font.Face
	bindings   *UnicodeBindings

	dc *gg.Context
}

type ForecastData struct {
	Location string
	Forecast []ForecastRow
}

type ForecastRow struct {
	IconCode    string
	Temperature float64
	Humidity    int
	Clouds      int
	Time        time.Time
}

func NewGenerator(fontPath, iconsPath, bindings string) (*Generator, error) {
	binds, err := LoadBindings(bindings)
	if err != nil {
		return nil, err
	}

	return &Generator{
		iconsFile:  iconsPath,
		fontFile:   fontPath,
		bindings:   binds,
		fontCache:  make(map[float64]font.Face),
		iconsCache: make(map[float64]font.Face),
	}, nil
}

func (g *Generator) loadFont(points float64) error {
	if _, ok := g.fontCache[points]; !ok {
		ff, err := gg.LoadFontFace(g.fontFile, points)
		if err != nil {
			return err
		}
		g.fontCache[points] = ff
	}
	g.dc.SetFontFace(g.fontCache[points])
	return nil
}

func (g *Generator) loadIcons(points float64) error {
	if _, ok := g.iconsCache[points]; !ok {
		ff, err := gg.LoadFontFace(g.iconsFile, points)
		if err != nil {
			return err
		}
		g.iconsCache[points] = ff
	}
	g.dc.SetFontFace(g.iconsCache[points])
	return nil
}

func (g *Generator) Generate(data *ForecastData) (*bytes.Buffer, error) {
	g.dc = gg.NewContext(400, 650)
	defer func() {
		g.dc.Clear()
		g.dc = nil
	}()

	g.prepareImage()
	g.drawWeatherLines()
	err := g.drawHeader(data.Location, data.Forecast[0])
	if err != nil {
		return nil, err
	}
	err = g.drawTime(data.Forecast[1:])
	if err != nil {
		return nil, err
	}
	err = g.drawClouds(data.Forecast[1:])
	if err != nil {
		return nil, err
	}
	err = g.drawHumidity(data.Forecast[1:])
	if err != nil {
		return nil, err
	}
	err = g.drawWeather(data.Forecast[1:])
	if err != nil {
		return nil, err
	}
	err = g.drawIcons(data.Forecast[1:])
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, g.dc.Image())
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (g *Generator) prepareImage() {
	if g.dc == nil {
		return
	}
	g.dc.SetRGBA(0, 0, 0, 0)
	g.dc.Clear()

	// Template
	g.dc.SetRGB255(242, 97, 73)
	g.dc.DrawRoundedRectangle(0, 0, 400, 650, 10)
	g.dc.Fill()
}

func (g *Generator) drawWeatherLines() {
	g.dc.SetRGB255(234, 89, 65)
	g.dc.DrawRectangle(0, 250, 400, 100)
	g.dc.DrawRectangle(0, 450, 400, 100)
	g.dc.Fill()

	g.dc.SetLineWidth(2)
	g.dc.SetRGBA(0, 0, 0, 0.05)
	g.dc.DrawLine(0, 250, 400, 250)
	g.dc.DrawLine(0, 349, 400, 348)
	g.dc.DrawLine(0, 450, 400, 450)
	g.dc.DrawLine(0, 549, 400, 548)
	g.dc.Stroke()
}

func (g *Generator) drawHeader(location string, forecast ForecastRow) error {
	// Header (place and date)
	err := g.loadFont(20)
	if err != nil {
		return err
	}

	g.dc.SetRGBA(1, 1, 1, 0.7)
	g.dc.DrawStringAnchored(location, 10, 15, 0, 0.5)
	g.dc.SetRGBA(1, 1, 1, 0.4)
	g.dc.DrawStringAnchored(time.Now().Format("Jan 2, 2006"), 270, 15, 0, 0.5)

	// First weather data
	g.dc.SetRGBA(1, 1, 1, 0.5)
	err = g.loadFont(30)
	if err != nil {
		return err
	}

	g.dc.DrawStringAnchored(fmt.Sprintf("%.2d:00", forecast.Time.Hour()), 50, 200, 0.5, 0.5)
	g.dc.DrawStringAnchored(fmt.Sprintf("H:%d%%", forecast.Humidity), 200, 200, 0.5, 0.5)
	g.dc.DrawStringAnchored(fmt.Sprintf("C:%d%%", forecast.Clouds), 350, 200, 0.5, 0.5)

	g.dc.SetRGBA(1, 1, 1, 1)
	err = g.loadFont(90)
	if err != nil {
		return err
	}

	g.dc.DrawStringAnchored(fmt.Sprintf("%d°", int(forecast.Temperature)), 100, 120, 0.5, 0.5)

	err = g.loadIcons(70)
	if err != nil {
		return err
	}

	g.dc.DrawStringAnchored(g.bindings.Get(forecast.IconCode), 250, 120, 0, 0.7)
	return nil
}

func (g *Generator) drawTime(data []ForecastRow) error {
	if len(data) < 4 {
		return errors.New("not enough data")
	}
	err := g.loadFont(30)
	if err != nil {
		return err
	}

	// Time
	g.dc.DrawStringAnchored(fmt.Sprintf("%.2v:00", data[0].Time.Hour()), 100, 285, 0, 0.5)
	g.dc.DrawStringAnchored(fmt.Sprintf("%.2v:00", data[1].Time.Hour()), 100, 385, 0, 0.5)
	g.dc.DrawStringAnchored(fmt.Sprintf("%.2v:00", data[2].Time.Hour()), 100, 485, 0, 0.5)
	g.dc.DrawStringAnchored(fmt.Sprintf("%.2v:00", data[3].Time.Hour()), 100, 585, 0, 0.5)

	return nil
}

func (g *Generator) drawHumidity(data []ForecastRow) error {
	if len(data) < 4 {
		return errors.New("not enough data")
	}
	err := g.loadFont(20)
	if err != nil {
		return err
	}

	g.dc.SetRGBA(1, 1, 1, 0.5)

	g.dc.DrawStringAnchored(fmt.Sprintf("H:%d%%", data[0].Humidity), 100, 315, 0, 0.5)
	g.dc.DrawStringAnchored(fmt.Sprintf("H:%d%%", data[1].Humidity), 100, 415, 0, 0.5)
	g.dc.DrawStringAnchored(fmt.Sprintf("H:%d%%", data[2].Humidity), 100, 515, 0, 0.5)
	g.dc.DrawStringAnchored(fmt.Sprintf("H:%d%%", data[3].Humidity), 100, 615, 0, 0.5)
	return nil
}

func (g *Generator) drawClouds(data []ForecastRow) error {
	if len(data) < 4 {
		return errors.New("not enough data")
	}
	err := g.loadFont(20)
	if err != nil {
		return err
	}

	g.dc.SetRGBA(1, 1, 1, 0.5)

	g.dc.DrawStringAnchored(fmt.Sprintf("C:%d%%", data[0].Clouds), 170, 315, 0, 0.5)
	g.dc.DrawStringAnchored(fmt.Sprintf("C:%d%%", data[1].Clouds), 170, 415, 0, 0.5)
	g.dc.DrawStringAnchored(fmt.Sprintf("C:%d%%", data[2].Clouds), 170, 515, 0, 0.5)
	g.dc.DrawStringAnchored(fmt.Sprintf("C:%d%%", data[3].Clouds), 170, 615, 0, 0.5)

	return nil
}

func (g *Generator) drawWeather(data []ForecastRow) error {
	if len(data) < 4 {
		return errors.New("not enough data")
	}
	err := g.loadFont(50)
	if err != nil {
		return err
	}

	g.dc.SetRGBA(1, 1, 1, 1)

	g.dc.DrawStringAnchored(fmt.Sprintf("%d°", int(data[0].Temperature)), 320, 300, 0.5, 0.5)
	g.dc.DrawStringAnchored(fmt.Sprintf("%d°", int(data[1].Temperature)), 320, 400, 0.5, 0.5)
	g.dc.DrawStringAnchored(fmt.Sprintf("%d°", int(data[2].Temperature)), 320, 500, 0.5, 0.5)
	g.dc.DrawStringAnchored(fmt.Sprintf("%d°", int(data[3].Temperature)), 320, 600, 0.5, 0.5)
	return nil
}

func (g *Generator) drawIcons(data []ForecastRow) error {
	err := g.loadIcons(50)
	if err != nil {
		return err
	}

	g.dc.SetRGBA(1, 1, 1, 1)

	g.dc.DrawStringAnchored(g.bindings.Get(data[0].IconCode), 20, 280, 0, 0.7)
	g.dc.DrawStringAnchored(g.bindings.Get(data[1].IconCode), 20, 380, 0, 0.7)
	g.dc.DrawStringAnchored(g.bindings.Get(data[2].IconCode), 20, 480, 0, 0.7)
	g.dc.DrawStringAnchored(g.bindings.Get(data[3].IconCode), 20, 580, 0, 0.7)
	return nil
}
