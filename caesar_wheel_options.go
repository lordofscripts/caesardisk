/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package caesardisk

import (
	"fmt"
	"image"
	"strings"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

// Default rendering options for the Caesar Wheel
var DefaultCaesarWheelOptions = CaesarWheelOptions{
	Title:           "",
	Size:            image.Rect(0, 0, 800, 800),
	Radius:          360.0, // Must be less than half the size (less than the diameter/2)
	Orthogonal:      true,
	RadialsColor:    NewRGB[uint8](0xd3, 0xd3, 0xd3), // #d3d3d3
	LettersFontPath: DEFAULT_FONT_BOLD,
	LettersSize:     30.0,
	LetterColor:     NewRGB[uint8](0, 0, 0),
	LetterColorAlt:  NewRGB[uint8](0, 0, 0),
	DigitsFontPath:  DEFAULT_FONT_REGULAR,
	DigitsSize:      18.0,
	DigitsColor:     NewRGB[uint8](0x00, 0x00, 0xff), // #0000ff OR #ffa500
}

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

// Rendering options for a Caesar encoder/decoder wheel
type CaesarWheelOptions struct {
	Title           string
	Size            image.Rectangle
	Radius          float64
	Orthogonal      bool
	RadialsColor    RGB[uint8]
	LettersFontPath string
	LettersSize     float64
	LetterColor     RGB[uint8]
	LetterColorAlt  RGB[uint8]
	DigitsFontPath  string
	DigitsSize      float64
	DigitsColor     RGB[uint8]
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer for the CaesarWheelOptions
func (w *CaesarWheelOptions) String() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%15s: %d x %d pixels\n", "Size", w.Size.Dx(), w.Size.Dy())
	fmt.Fprintf(&sb, "%15s: %d pixels\n", "Radius", int(w.Radius))
	fmt.Fprintf(&sb, "%15s: %s @ %.3f\n", "Text font", w.LettersFontPath, w.LettersSize)
	fmt.Fprintf(&sb, "%15s: %s @ %.3f\n", "Digit font", w.DigitsFontPath, w.DigitsSize)
	fmt.Fprintf(&sb, "%15s: %t\n", "Orthogonal", w.Orthogonal)

	return sb.String()
}

// set font path and size for letters (alphabets)
func (w *CaesarWheelOptions) SetLetterFont(path string, size float64) *CaesarWheelOptions {
	w.LettersFontPath = path
	w.LettersSize = size

	return w
}

// set font path and size for digits (Caesar key shift)
func (w *CaesarWheelOptions) SetDigitFont(path string, size float64) *CaesarWheelOptions {
	w.DigitsFontPath = path
	w.DigitsSize = size

	return w
}

func (w *CaesarWheelOptions) SetLetterColors(outer, inner string) *CaesarWheelOptions {
	w.LetterColor = NewRGBFromString(outer)
	w.LetterColorAlt = NewRGBFromString(inner)
	return w
}

func (w *CaesarWheelOptions) SetDigitColor(colorHex string) *CaesarWheelOptions {
	w.DigitsColor = NewRGBFromString(colorHex)
	return w
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/
