package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"regexp"
	"strconv"
)

// Color represents a 24-bit color, 8 bits each of red, green, and blue.
type Color color.RGBA

func (c Color) RGBA() (uint32, uint32, uint32, uint32) {
	return color.RGBA(c).RGBA()
}

var _ color.Color = Color{}

var ErrInvalidColorString = errors.New("color string must match pattern #RRGGBB")

// NewColorFromString creates a Color from a #RRGGBB string.
func NewColorFromString(c string) (Color, error) {
	colorRegex := regexp.MustCompile("(?i)^#[0-9a-f]{6}$")
	if !colorRegex.MatchString(c) {
		return Color{}, ErrInvalidColorString
	}

	// Errors can be ignored because the regex has guaranteed valid digits.
	r, _ := strconv.ParseUint(c[1:3], 16, 8)
	g, _ := strconv.ParseUint(c[3:5], 16, 8)
	b, _ := strconv.ParseUint(c[5:7], 16, 8)

	return Color{R: uint8(r), G: uint8(g), B: uint8(b)}, nil
}

// NewColor creates a Color from a color.Color instance.
func NewColor(c color.Color) Color {
	// color.Color's R, G, B methods return uint32, using 16 bits thereof, but
	// we want uint8, so shift everything right by 8 bits to capture the most
	// significant bits.
	r, g, b, _ := c.RGBA()
	return Color{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
	}
}

// String renders a Color as a #RRGGBB string.
func (c Color) String() string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}

// ColorHistogram tracks frequency of colors, and can track the
// top N colors in the histogram.
type ColorHistogram struct {
	hist map[color.Color]uint
	topN []color.Color
}

// Add a pixel's color information to the histogram.
func (ch *ColorHistogram) Add(c color.Color) {
	if ch.hist == nil {
		ch.hist = make(map[color.Color]uint)
	}
	ch.hist[c] += 1
	ch.topN = nil
}

// Return the top N colors as the pixels are added to the histogram.  In the
// case where multiple colors have the same incidence, they will all be grouped
// together in the TopN slice (e.g., a 2-way tie for top incidence will be the
// first and second elements of the TopN slice), but their exact order within
// the tied range of elemnts is not guaranteed.
func (ch *ColorHistogram) TopN(n int) []color.Color {
	if ch.topN != nil {
		return ch.topN
	}
	if l := len(ch.hist); l < n {
		n = l
	}
	ch.topN = make([]color.Color, n)
	for k, v := range ch.hist {
		if minTopN := ch.hist[ch.topN[len(ch.topN)-1]]; v <= minTopN {
			// This value is not greater than the Nth-highest, so skip it
			continue
		}
		var i int
		for i = 0; i < len(ch.topN); i++ {
			if v > ch.hist[ch.topN[i]] {
				break
			}
		}
		newTop := make([]color.Color, 0, len(ch.topN))
		newTop = append(newTop, ch.topN[:i]...)
		newTop = append(newTop, k)
		ch.topN = append(newTop, ch.topN[i:len(ch.topN)-1]...)
	}
	return ch.topN
}

// Scan builds a color histogram for the provided image.
func Scan(img image.Image) ColorHistogram {
	var (
		ch ColorHistogram
		m  = img.Bounds()
	)
	for y := m.Min.Y; y < m.Max.Y; y++ {
		for x := m.Min.X; x < m.Max.X; x++ {
			ch.Add(m.At(x, y))
		}
	}
	return ch
}
