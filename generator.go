/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   goCaesarDisk
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package caesardisk

import (
	"embed"
	"errors"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	// Keep this in mind if you want an dual letter & symbol
	// disks that you can use with the same encoding key
	// (same length for both alphabets)
	// 0        1         2         3
	// 1234567890123456789012345678901
	// -------------------------------
	// ABCDEFGHIJKLMNÑOPQRSTUVWXYZÁÉÍÓÚ
	// !"#$%&'()*+,-./ 0123456789:;<=>?
	Alpha_EN string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Alpha_ES string = "ABCDEFGHIJKLMNÑOPQRSTUVWXYZÁÉÍÓÚ"
	Alpha_CZ string = "ABCČDĎEFGHIJKLMNŇOPQRŘSŠTŤUVWXYÝZŽÁÉÍÓÚĚŮ"
	Alpha_DE string = "ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÜẞ"
	Alpha_IT string = "ABCDEFGHILMNOPQRSTUVZÉÓÀÈÌÒÙ"
	Alpha_PT string = "ABCÇDEFGHIJKLMNOPQRSTUVWXYZÁÉÍÓÚÀÂÊÔÃÕ"
	Alpha_GR string = "ΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩ"
	Alpha_RU string = "АБВГДЕËЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"
	Alpha_PU string = `!"#$%&'()*+,-./ 0123456789:;<=>?`

	Alpha_ES_DUAL    = "ABCDEFGHIJKLMNÑOPQRSTUVWXYZ"
	Alpha_PU_DUAL_ES = `!"#$%&()*+,-./ 0123456789=?` // to match length of ES_DUAL
	Alpha_PU_DUAL_EN = `!"#$%&()*+,-./ 0123456789?`  // to match length of EN

	SubTitle string = "C a e s a r  D i s k"
)

// Default Free embedded font(s)
//
//go:embed ubuntu.regular.ttf
var fontUbuntuRegular embed.FS

// Default Free embedded font(s)
//
//go:embed ubuntu.bold.ttf
var fontUbuntuBold embed.FS

// Default rendering options for the Caesar Wheel
var DefaultCaesarWheelOptions = CaesarWheelOptions{
	Title:           "",
	Size:            image.Rect(0, 0, 800, 800),
	Radius:          360.0, // Must be less than half the size (less than the diameter/2)
	Orthogonal:      true,
	RadialsColor:    NewRGB[uint8](0xd3, 0xd3, 0xd3), // #d3d3d3
	LettersFontPath: "ubuntu.bold.ttf",
	LettersSize:     30.0,
	LetterColor:     NewRGB[uint8](0, 0, 0),
	DigitsFontPath:  "ubuntu.regular.ttf",
	DigitsSize:      18.0,
	DigitsColor:     NewRGB[uint8](0x00, 0x00, 0xff), //NewRGB[uint8](0xff, 0xa5, 0x00),
}

var (
	ErrUnmatchedDual error = errors.New("for dual alphabet mode both alphabets must have equal lengths")
)

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
	DigitsFontPath  string
	DigitsSize      float64
	DigitsColor     RGB[uint8]
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

type ringPositions struct {
	// The radius for the edge of the circle
	Edge float64
	// The disk's outermost alphabet (letters) mid-position
	First float64
	// The disk's innermost alphabet (symbols) mid-position
	Second float64
	// The disk's numerical shift value (key) OUTER only
	Shift float64
	// The disk's arc text radius. OUTER ONLY
	ArcLabelPos float64
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer for the CaesarWheelOptions
func (w CaesarWheelOptions) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%15s: %d x %d pixels\n", "Size", w.Size.Dx(), w.Size.Dy())
	fmt.Fprintf(&sb, "%15s: %d pixels\n", "Radius", int(w.Radius))
	fmt.Fprintf(&sb, "%15s: %s @ %.3f\n", "Text font", w.LettersFontPath, w.LettersSize)
	fmt.Fprintf(&sb, "%15s: %s @ %.3f\n", "Digit font", w.DigitsFontPath, w.DigitsSize)
	fmt.Fprintf(&sb, "%15s: %t\n", "Orthogonal", w.Orthogonal)
	return sb.String()
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

// load a font given its path and set its size
func loadFont(fontPath string, size float64) font.Face {
	if face, err := gg.LoadFontFace(fontPath, size); err != nil {
		log.Printf("could not load font. %v. Using default font.", err)
		if face, err = gg.LoadFontFace("sans-serif", size); err != nil {
			log.Printf("could not load generic font. %v.", err)
			panic(err)
		}
		return face
	} else {
		return face
	}

}

// Draw a text in a semicircle (arc)
func drawArcText(arcText string, fgColor color.Color, fontSize, width, height, radius float64, onLeft bool, fontPath string, dc *gg.Context) {
	// Calculate the arc text position
	arcLength := float64(len(arcText)) * fontSize * 0.5 // Estimate arc length based on font size
	startAngle := -arcLength / (2 * radius)             // Start angle for the text arc

	arcFace := loadFont(fontPath, fontSize)

	for i, r := range arcText {
		// Calculate position for each character
		angle := startAngle + float64(i)*fontSize/(2*radius)
		if onLeft {
			angle += math.Pi
		}

		x := width/2 + (radius-50)*math.Cos(angle) // Offset by 50 pixels for positioning
		y := height/2 + (radius-50)*math.Sin(angle)

		// Rotate and draw each character
		dc.Push()
		dc.Translate(x, y)
		dc.Rotate(angle + math.Pi/2) // Rotate to align character
		//dc.SetRGB(0, 0, 0)           // black color for the text
		dc.SetColor(fgColor)
		dc.SetFontFace(arcFace)
		dc.DrawStringAnchored(string(r), 0, 0, 0.5, 0.5) // Draw character
		dc.Pop()
	}
}

func GenerateDualCaesarWheel(letters, symbols string, filename string, inner bool, opts CaesarWheelOptions) error {
	sizeLet := utf8.RuneCountInString(letters)
	sizeSym := utf8.RuneCountInString(symbols)
	if sizeLet != sizeSym {
		println("LET", sizeLet, letters)
		println("SYM", sizeSym, symbols)
		return ErrUnmatchedDual
	}

	// I. Setup
	// for outer disk, contains 2 alphabets (letters & symbols) and shift number (key)
	ringDualOut := ringPositions{
		Edge:        opts.Radius,
		First:       opts.Radius * 0.95, // at 90% (0.90) place the character
		Second:      opts.Radius * 0.85,
		Shift:       opts.Radius * 0.55, // place digits 70% of the way to the edge
		ArcLabelPos: opts.Radius * 0.55,
	}
	// for inner disk, contains 2 cipher alphabets but only cut-out for digits
	ringDualIn := ringPositions{
		Edge:   opts.Radius * 0.80,
		First:  opts.Radius * 0.75, // at 90% (0.90) place the character
		Second: opts.Radius * 0.65,
	}

	// Load a TrueType font file, i.e. Arial.ttf
	textFace := loadFont(opts.LettersFontPath, opts.LettersSize)
	shiftFace := loadFont(opts.DigitsFontPath, opts.DigitsSize)

	// center point
	x := float64(opts.Size.Dx() / 2)
	y := float64(opts.Size.Dy() / 2)

	// II. Draw Circle for chosen disk
	dc := gg.NewContext(opts.Size.Dx(), opts.Size.Dy())
	dc.Clear()               // start with transparent background
	dc.SetColor(color.White) // Fill & Stroke

	// Draw the main circle outline
	dc.SetLineWidth(2)
	if inner { // Inner disk with outer transparency
		dc.DrawCircle(x, y, ringDualIn.Edge) // was Radius - LetterSize - 8
		dc.Fill()

		dc.SetLineWidth(2)
		dc.SetRGB255(0xd3, 0xd3, 0xd3)
		dc.DrawCircle(x, y, ringDualIn.Edge-2)
		dc.Stroke()
	} else { // Disk with outer transparency
		dc.DrawCircle(x, y, ringDualOut.Edge)
		dc.Fill() // fill with white

		dc.SetColor(color.Black)
		dc.DrawCircle(x, y, ringDualOut.Edge)
		dc.Stroke()
	}

	// III. Draw the N dividing lines and characters
	letterLabel := []rune(letters) // each letter MAY be a multi-byte rune
	symbolLabel := []rune(symbols)
	for i := range sizeLet { // the length in Unicode chars rather than bytes
		// Calculate the start and end angles for the segment
		startAngle := (float64(i) / float64(sizeLet)) * 2 * math.Pi
		//endAngle := (float64(i+1) / float64(sizeLet)) * 2 * math.Pi
		// Calculate the middle angle for text  placement
		midAngle := startAngle + (math.Pi / float64(sizeLet)) // use pi/N for half a segment angle
		// Calculate the end point of the line on the circle's edge
		endX := x + ringDualOut.Edge*math.Cos(startAngle)
		endY := y + ringDualOut.Edge*math.Sin(startAngle)

		// Draw the dividingline (radials)
		dc.MoveTo(x, y)
		dc.SetLineWidth(1)
		dc.SetRGB255(int(opts.RadialsColor.Red), int(opts.RadialsColor.Green), int(opts.RadialsColor.Blue)) // mid gray 0x666a6d
		dc.LineTo(endX, endY)
		dc.Stroke()

		// -- Calculate the letter's baseline angle so that it is perpendicular to
		// the radius. First calculation is letter's side edge parallel to the circle's tangent
		angle := math.Pi*2*float64(i)/float64(sizeLet) + math.Pi/float64(sizeLet)
		if opts.Orthogonal {
			angle = angle + math.Pi/2 // read at XII
		} // else read at III

		// Determine the character labelLetter (A, B, C, ...) of the chosen alphabet
		labelLetter := letterLabel[i]
		labelSymbol := symbolLabel[i]
		dc.SetFontFace(textFace)

		if inner {
			drawInnerChar := func(textAltRadius float64, labelChar rune, hexColor string) {
				textAltX := x + textAltRadius*math.Cos(midAngle)
				textAltY := y + textAltRadius*math.Sin(midAngle)
				//dc.SetColor(opts.LetterColor.ToColor())
				dc.SetHexColor(hexColor)
				dc.Push()
				dc.Translate(textAltX, textAltY)
				dc.Rotate(angle)
				dc.DrawStringAnchored(string(labelChar), 0, 0, 0.5, 0.5)
				dc.Pop()
			}

			// The Inner disk contains the CIPHERed character(s)
			drawInnerChar(ringDualIn.First, labelLetter, opts.LetterColor.ToHexColor())
			drawInnerChar(ringDualIn.Second, labelSymbol, MidGray.ToHexColor())

			// Draw the cut-out that would allow to see-through to
			// the key/index printed in the disc underneath (outer)
			endAngle := (float64(1) / float64(sizeLet)) * 2 * math.Pi
			dc.SetLineWidth(0.75)
			dc.SetRGB(0.827, 0.827, 0.827)
			dc.DrawArc(x, y, ringDualOut.Shift+opts.DigitsSize-2, 0, endAngle)
			dc.Stroke()
			dc.DrawArc(x, y, ringDualOut.Shift-opts.DigitsSize-2, 0, endAngle)
			dc.Stroke()
		} else {
			drawOuterChar := func(textRadius float64, labelChar rune, hexColor string) {
				// -- Place the Character label near the edge
				// Calculate text position (mid-radius, mid-angle)
				textX := x + textRadius*math.Cos(midAngle)
				textY := y + textRadius*math.Sin(midAngle)

				//dc.SetColor(opts.LetterColor.ToColor())
				dc.SetHexColor(hexColor)
				dc.Push()
				dc.Translate(textX, textY)
				dc.Rotate(angle)
				// DrawStringAnchored aligns the text's center point ot the calculated (textX, textY)
				// 0.5 are the anchor points meaning 50% horizontal offset and 50% vertical offset,
				// thus centered.
				dc.DrawStringAnchored(string(labelChar), 0, 0, 0.5, 0.5)
				dc.Pop()
			}

			drawOuterChar(ringDualOut.First, labelLetter, opts.LetterColor.ToHexColor())
			drawOuterChar(ringDualOut.Second, labelSymbol, MidGray.ToHexColor())

			// -- Place the Index label below (must be seen through cut-out in inner disk)

			// Calculate index position (mid-radius, mid-angle)
			digitsX := x + ringDualOut.Shift*math.Cos(midAngle)
			digitsY := y + ringDualOut.Shift*math.Sin(midAngle)

			labelShiftValue := fmt.Sprintf("%02d", i)
			dc.SetFontFace(shiftFace) // set different font face or size
			dc.SetRGB255(int(opts.DigitsColor.Red), int(opts.DigitsColor.Green), int(opts.DigitsColor.Blue))

			dc.Push()
			dc.Translate(digitsX, digitsY)
			dc.Rotate(angle)
			dc.DrawStringAnchored(labelShiftValue, 0, 0, 0.5, 0.5)
			dc.Pop()
		}
	}

	// IV. Epilogue
	// 4.1 Draw a black dot in the middle to aid in making the pinhole
	const DOT_RADIUS = 3 // Middle dot radius in pixels
	dc.SetRGB(0, 0, 0)   // Set the color to black
	dc.DrawCircle(float64(opts.Size.Dx())/2, float64(opts.Size.Dy())/2, DOT_RADIUS)
	dc.Fill() // Fill the circle to make it look like a solid dot

	// 4.2 Arc Text
	// · Arc Text: "Caesar Disk"
	// · Optional Arc Text: Title
	drawArcText(SubTitle, color.Black, opts.LettersSize, float64(opts.Size.Dx()), float64(opts.Size.Dy()), ringDualOut.ArcLabelPos, true, "ubuntu.bold.ttf", dc)
	if len(opts.Title) != 0 {
		drawArcText(opts.Title, color.Black, opts.LettersSize, float64(opts.Size.Dx()), float64(opts.Size.Dy()), ringDualOut.ArcLabelPos, false, "ubuntu.regular.ttf", dc)
	}

	// 4.3 Disk Assembly information
	// · Copyright
	// · Inner OR Outer assembly
	copyrightFace := loadFont(opts.LettersFontPath, 10.0)
	dc.SetFontFace(copyrightFace)
	dc.SetColor(Yellow.ToColor())
	textInfo := "Outer Disk"
	if inner {
		textInfo = "Inner Disk"
	}
	dc.DrawStringAnchored(textInfo, float64(10), float64(opts.Size.Dy()-20), 0, 0.5)

	// 4.4 Donations
	// · Copyright & Donations text
	copyrightText := "https://buymeacoffee.com/lostinwriting"
	dc.DrawStringAnchored(copyrightText, float64(opts.Size.Dx()/2), float64(opts.Size.Dy()-20), 0.5, 0.5)

	// · Text font used in rendering
	basename := filepath.Base(opts.LettersFontPath)
	dc.DrawStringAnchored(basename, float64(opts.Size.Dx())*0.80, float64(opts.Size.Dy()-20), 0.5, 0.5)

	if err := dc.SavePNG(filename); err != nil { // or use SaveJPG
		return fmt.Errorf("failed to save image: %w", err)
	}

	fmt.Printf("successfully generated '%s' with %d segments\n", filename, sizeLet)
	return nil
}

// generate an image with a Caesar encoder wheel that could be printed
func GenerateCaesarWheel(letters string, filename string, inner bool, opts CaesarWheelOptions) error {
	// at 75% (0.75) place the 2-digit offset
	var indexRadius = opts.Radius * 0.70 // place digits 70% of the way to the edge
	// at 90% (0.90) place the character, thus nearest the edge of the circle
	var textRadius = opts.Radius * 0.95

	// center point
	x := float64(opts.Size.Dx() / 2)
	y := float64(opts.Size.Dy() / 2)

	dc := gg.NewContext(opts.Size.Dx(), opts.Size.Dy())
	dc.Clear()         // start with transparent background
	dc.SetRGB(1, 1, 1) // White background

	// Load a TrueType font file, i.e. Arial.ttf
	textFace := loadFont(opts.LettersFontPath, opts.LettersSize)
	counterFace := loadFont(opts.DigitsFontPath, opts.DigitsSize)

	// Draw the main circle outline
	dc.SetLineWidth(2)
	if !inner { // Disk with outer transparency
		dc.DrawCircle(x, y, opts.Radius)
		dc.Fill() // fill with white

		dc.SetColor(color.Black)
		dc.DrawCircle(x, y, opts.Radius)
		dc.Stroke()
	} else { // Inner disk with outer transparency
		dc.DrawCircle(x, y, opts.Radius-opts.LettersSize-8)
		dc.Fill()

		dc.SetLineWidth(2)
		dc.SetRGB255(0xd3, 0xd3, 0xd3)
		dc.DrawCircle(x, y, opts.Radius-opts.LettersSize-10)
		dc.Stroke()
	}

	// Draw the N dividing lines and characters
	letterLabel := []rune(letters) // each letter MAY be a multi-byte rune
	n := len(letterLabel)          // the length in Unicode chars rather than bytes
	for i := range n {
		// Calculate the start and end angles for the segment
		startAngle := (float64(i) / float64(n)) * 2 * math.Pi
		//endAngle := (float64(i+1) / float64(n)) * 2 * math.Pi
		// Calculate the middle angle for text  placement
		midAngle := startAngle + (math.Pi / float64(n)) // use pi/N for half a segment angle
		// Calculate the end point of the line on the circle's edge
		endX := x + opts.Radius*math.Cos(startAngle)
		endY := y + opts.Radius*math.Sin(startAngle)

		// Draw the dividingline (radials)
		dc.MoveTo(x, y)
		dc.SetLineWidth(1)
		dc.SetRGB255(int(opts.RadialsColor.Red), int(opts.RadialsColor.Green), int(opts.RadialsColor.Blue)) // mid gray 0x666a6d
		dc.LineTo(endX, endY)
		dc.Stroke()

		// -- Calculate the letter's baseline angle so that it is perpendicular to
		// the radius. First calculation is letter's side edge parallel to the circle's tangent
		angle := math.Pi*2*float64(i)/float64(n) + math.Pi/float64(n)
		if opts.Orthogonal {
			angle = angle + math.Pi/2 // read at XII
		} // else read at III

		// Determine the character label (A, B, C, ...) of the chosen alphabet
		label := string(letterLabel[i])
		dc.SetFontFace(textFace)

		if !inner {
			// -- Place the Character label near the edge
			// Calculate text position (mid-radius, mid-angle)
			textX := x + textRadius*math.Cos(midAngle)
			textY := y + textRadius*math.Sin(midAngle)

			dc.SetColor(opts.LetterColor.ToColor())
			dc.Push()
			dc.Translate(textX, textY)
			dc.Rotate(angle)
			// DrawStringAnchored aligns the text's center point ot the calculated (textX, textY)
			// 0.5 are the anchor points meaning 50% horizontal offset and 50% vertical offset,
			// thus centered.
			dc.DrawStringAnchored(label, 0, 0, 0.5, 0.5)
			dc.Pop()

			// -- Place the Index label below

			// Calculate index position (mid-radius, mid-angle)
			digitsX := x + indexRadius*math.Cos(midAngle)
			digitsY := y + indexRadius*math.Sin(midAngle)

			indexLabel := fmt.Sprintf("%02d", i)
			dc.SetFontFace(counterFace) // set different font face or size
			dc.SetRGB255(int(opts.DigitsColor.Red), int(opts.DigitsColor.Green), int(opts.DigitsColor.Blue))

			dc.Push()
			dc.Translate(digitsX, digitsY)
			dc.Rotate(angle)
			dc.DrawStringAnchored(indexLabel, 0, 0, 0.5, 0.5)
			dc.Pop()
		} else {
			var textAltRadius = opts.Radius * 0.85

			textAltX := x + textAltRadius*math.Cos(midAngle)
			textAltY := y + textAltRadius*math.Sin(midAngle)
			dc.SetColor(opts.LetterColor.ToColor())
			dc.Push()
			dc.Translate(textAltX, textAltY)
			dc.Rotate(angle)
			dc.DrawStringAnchored(label, 0, 0, 0.5, 0.5)
			dc.Pop()

			// Draw the cut-out that would allow to see-through to
			// the key/index printed in the disc underneath (outer)
			endAngle := (float64(1) / float64(n)) * 2 * math.Pi
			dc.SetLineWidth(0.75)
			dc.SetRGB(0.827, 0.827, 0.827)
			dc.DrawArc(x, y, indexRadius+opts.DigitsSize-2, 0, endAngle)
			dc.Stroke()
			dc.DrawArc(x, y, indexRadius-opts.DigitsSize-2, 0, endAngle)
			dc.Stroke()
		}
	}

	// Draw a black dot in the middle to aid in making the pinhole
	const DOT_RADIUS = 3 // Middle dot radius in pixels
	dc.SetRGB(0, 0, 0)   // Set the color to black
	dc.DrawCircle(float64(opts.Size.Dx())/2, float64(opts.Size.Dy())/2, DOT_RADIUS)
	dc.Fill() // Fill the circle to make it look like a solid dot

	// -- Epilogue
	// · Arc Text: "Caesar Disk"
	// · Optional Arc Text: Title
	drawArcText(SubTitle, color.Black, opts.LettersSize, float64(opts.Size.Dx()), float64(opts.Size.Dy()), opts.Radius*0.65, true, "ubuntu.bold.ttf", dc)
	if len(opts.Title) != 0 {
		drawArcText(opts.Title, color.Black, opts.LettersSize, float64(opts.Size.Dx()), float64(opts.Size.Dy()), opts.Radius*0.65, false, "ubuntu.regular.ttf", dc)
	}

	// · Disk set assembly information (inner OR outer disk)
	copyrightFace := loadFont(opts.LettersFontPath, 10.0)
	dc.SetFontFace(copyrightFace)
	dc.SetColor(Yellow.ToColor())
	textInfo := "Outer Disk"
	if inner {
		textInfo = "Inner Disk"
	}
	dc.DrawStringAnchored(textInfo, float64(10), float64(opts.Size.Dy()-20), 0, 0.5)

	// · Copyright & Donations text
	copyrightText := "https://buymeacoffee.com/lostinwriting"
	dc.DrawStringAnchored(copyrightText, float64(opts.Size.Dx()/2), float64(opts.Size.Dy()-20), 0.5, 0.5)

	// · Text font used in rendering
	basename := filepath.Base(opts.LettersFontPath)
	dc.DrawStringAnchored(basename, float64(opts.Size.Dx())*0.80, float64(opts.Size.Dy()-20), 0.5, 0.5)

	if err := dc.SavePNG(filename); err != nil { // or use SaveJPG
		return fmt.Errorf("failed to save image: %w", err)
	}

	fmt.Printf("successfully generated '%s' with %d segments\n", filename, n)
	return nil
}
