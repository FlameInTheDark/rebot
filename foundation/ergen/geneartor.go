package ergen

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

const (
	rateHeaderHeight = 100
	rateRowHeight    = 100
	rateWidth        = 400
)

var (
	errNilContext = errors.New("draw context is nil")
)

// Generator is an exchange rates image generator
type Generator struct {
	fontFile  string
	fontCache map[string]font.Face
}

// NewGenerator creates a new Generator
func NewGenerator(fontFile string) *Generator {
	return &Generator{fontFile: fontFile}
}

// Rates rates data
type Rates struct {
	Base   string
	Amount float64
	Rates  []RateRow
}

// RateRow contains rate data for the currency
type RateRow struct {
	Currency string
	Exchange float64
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

func (g *Generator) prepareImage(dc *gg.Context, height, width float64) error {
	if dc == nil {
		return errNilContext
	}

	dc.SetRGBA(0, 0, 0, 0)
	dc.Clear()

	// Template
	dc.SetRGB255(242, 97, 73)
	dc.DrawRoundedRectangle(0, 0, width, height, 10)
	dc.Fill()
	return nil
}

func (g *Generator) drawRateLine(dc *gg.Context, name string, amount, position float64) error {
	if dc == nil {
		return errNilContext
	}

	dc.SetRGBA(0, 0, 0, 0)
	dc.Clear()

	dc.SetRGB255(234, 89, 65)
}

// GenerateExchangeRates creates a new exchange rates image
func (g *Generator) GenerateExchangeRates(rates Rates) (*bytes.Buffer, error) {
	pictureHeight := (len(rates.Rates) * rateRowHeight) + rateHeaderHeight
	dc := gg.NewContext(rateWidth, pictureHeight)

	err := g.prepareImage(dc, float64(pictureHeight), rateWidth)

}
