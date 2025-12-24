/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							     goCaesarDisk
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package caesardisk

import (
	"fmt"
	"image/color"
)

var Yellow = NewRGBA[uint8](0xff, 0xfd, 0x01, 0xFF)
var Gray = NewRGBA[uint8](0xd3, 0xd3, 0xd3, 0xFF)
var MidGray = NewRGBA[uint8](0x66, 0x6a, 0x6d, 0xFF)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

// a custom RGBA implementation for the Caesar Wheel
type RGB[T float64 | uint8] struct {
	Red   T
	Green T
	Blue  T
	Alpha T
}

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// RGB with maximum alpha (no opacity)
func NewRGB[T float64 | uint8](r, g, b T) RGB[T] {
	var alpha T
	switch any(alpha).(type) {
	case float64:
		alpha = 1.0
	case uint8:
		alpha = 0xff
	}

	return RGB[T]{
		Red:   r,
		Green: g,
		Blue:  b,
		Alpha: alpha,
	}
}

// RGB with specified alpha
func NewRGBA[T float64 | uint8](r, g, b, a T) RGB[T] {
	return RGB[T]{
		Red:   r,
		Green: g,
		Blue:  b,
		Alpha: a,
	}
}

func NewRGBFromString(colorHex string) RGB[uint8] {
	if c, err := ParseHexColor(colorHex); err != nil {
		panic(err)
	} else {
		return c
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

func (r RGB[T]) ToHexColor() string {
	var paint string
	switch any(r.Alpha).(type) {
	case float64:
		red, _ := any(r.Red).(float64)
		green := any(r.Green).(float64)
		blue := any(r.Blue).(float64)
		paint = fmt.Sprintf("#%02x%02x%02x", uint8(red*255), uint8(green*255), uint8(blue*255))
	case uint8:
		red, _ := any(r.Red).(uint8)
		green := any(r.Green).(uint8)
		blue := any(r.Blue).(uint8)
		paint = fmt.Sprintf("#%02x%02x%02x", red, green, blue)
	}
	return paint
}

// convert to a Go color from the standard library
func (r RGB[T]) ToColor() color.Color {
	var paint color.Color
	switch any(r.Alpha).(type) {
	case float64:
		red, _ := any(r.Red).(float64)
		green := any(r.Green).(float64)
		blue := any(r.Blue).(float64)
		alpha := any(r.Alpha).(float64)
		paint = color.RGBA{R: uint8(red * 255), G: uint8(green * 255), B: uint8(blue * 255), A: uint8(alpha * 255)}
	case uint8:
		red, _ := any(r.Red).(uint8)
		green := any(r.Green).(uint8)
		blue := any(r.Blue).(uint8)
		alpha := any(r.Alpha).(uint8)
		paint = color.RGBA{R: red, G: green, B: blue, A: alpha}
	}
	return paint
}

// implements fmt.Stringer
func (r RGB[T]) String() string {
	var me string
	switch any(r.Alpha).(type) {
	case float64:
		red, _ := any(r.Red).(float64)
		green := any(r.Green).(float64)
		blue := any(r.Blue).(float64)
		alpha := any(r.Alpha).(float64)
		me = fmt.Sprintf("RGBA(%f,%f,%f,%f)", red, green, blue, alpha)
	case uint8:
		red, _ := any(r.Red).(uint8)
		green := any(r.Green).(uint8)
		blue := any(r.Blue).(uint8)
		alpha := any(r.Alpha).(uint8)
		me = fmt.Sprintf("RGBA(#%2X%2X%2X%2X)", red, green, blue, alpha)
	default:
		me = ""
	}
	return me
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

// Parses a string that describes a color in Hex format. using
// one or two digits for each RGB part. The string must start
// with a # and no alpha part. Examples: #a5d or #ab24ef
func ParseHexColor(s string) (c RGB[uint8], err error) {
	c.Alpha = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.Red, &c.Green, &c.Blue)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.Red, &c.Green, &c.Blue)
		// Double the hex digits:
		c.Red *= 17
		c.Green *= 17
		c.Blue *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")

	}
	return
}
