package main

import (
	"errors"
	"fmt"
	"image"
	"regexp"
	"strconv"
)

// Color represents a 24-bit color, 8 bits each of red, green, and blue.
type Color uint32

var ErrInvalidColorString = errors.New("color string must match pattern #RRGGBB")

// NewColor creates a Color from a #RRGGBB string.
func NewColor(c string) (Color, error) {
	colorRegex := regexp.MustCompile("(?i)^#[0-9a-f]{6}$")
	if !colorRegex.MatchString(c) {
		return Color(0), ErrInvalidColorString
	}

	// Errors can be ignored because the regex has guaranteed valid digits.
	r, _ := strconv.ParseUint(c[1:3], 16, 8)
	g, _ := strconv.ParseUint(c[3:5], 16, 8)
	b, _ := strconv.ParseUint(c[5:7], 16, 8)

	color := (r << 16) | (g << 8) | b
	return Color(color), nil
}

// R returns the red component of a Color.
func (c Color) R() uint8 {
	return uint8((c & (0xFF << 16)) >> 16)
}

// G returns the green component of a Color.
func (c Color) G() uint8 {
	return uint8((c & (0xFF << 8)) >> 8)
}

// B returns the blue component of a Color.
func (c Color) B() uint8 {
	return uint8(c & 0xFF)
}

// String renders a Color as a #RRGGBB string.
func (c Color) String() string {
	return fmt.Sprintf("#%02X%02X%02X", c.R(), c.G(), c.B())
}

// TODO ColorHistogram
type ColorHistogram map[Color]uint

func (ch ColorHistogram) TopN(n int) {
}

func NewColorHistogram(img *image.Image) ColorHistogram {
	return nil
}
