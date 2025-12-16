/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   goCaesarDisk
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/lordofscripts/caesardisk"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/
const ()

var ()

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

func Usage() {
	fmt.Println("Usage:")
	fmt.Println("\tcaesardisk [options]")
	fmt.Println("\tcaesardisk [options] -text-font FONT.ttf")
	fmt.Println("\tcaesardisk [options] -text-font FONT.ttf -digit-font FONT.ttf")
	fmt.Println("Options:")
	flag.PrintDefaults()
	fmt.Println("Note: Fonts must be TrueType (*.ttf)")
}

/* ----------------------------------------------------------------
 *					M A I N    |     D E M O
 *-----------------------------------------------------------------*/

func main() {
	// The Caesar's wheel alphabet determines how many slices
	// · When true (orthogonal) the letter's mid axis is aligned with the
	// circle's radius as in most common Caesar wheels out there. You
	// read the letter at the XII o'clock position.
	// · When false the letter's mid axis is parallel to the circle's
	// tangent, thus the letter appears parallel to the edge. You read
	// the letter at the III o'clock position.

	// I. Command-line flag definition and parsing
	var flgHelp, flgES, flgPunct bool
	var flgTextFontPath, flgDigitFontPath, flgAlphabet, flgTitle string
	flag.Usage = Usage
	flag.BoolVar(&flgHelp, "help", false, "This help")
	flag.BoolVar(&flgES, "ES", false, "Spanish alphabet (overrides -alpha)")
	flag.BoolVar(&flgPunct, "PU", false, "Punctuation and numerical alphabet (overrides -alpha)")
	flag.StringVar(&flgTitle, "title", "", "Title (usually disk language or ID)")
	flag.StringVar(&flgAlphabet, "alpha", "", "Alphabet defaults to English ASCII alphabet")
	flag.StringVar(&flgTextFontPath, "text-font", "", "Text font path")
	flag.StringVar(&flgDigitFontPath, "digit-font", "", "Digit font path (else use same text font)")
	flag.Parse()

	// II. Command-line flag validation and processing
	caesardisk.Copyright(caesardisk.CO1, true)

	if flgHelp {
		flag.Usage()
	} else {
		var Options caesardisk.CaesarWheelOptions = caesardisk.DefaultCaesarWheelOptions
		if len(flgDigitFontPath) != 0 {
			Options.DigitsFontPath = flgDigitFontPath
		}
		if len(flgTextFontPath) != 0 {
			Options.LettersFontPath = flgTextFontPath
			if len(flgDigitFontPath) == 0 {
				Options.DigitsFontPath = flgTextFontPath
			}
		}
		if len(flgTitle) != 0 {
			flgTitle := strings.Join(strings.Split(flgTitle, ""), " ")
			Options.Title = flgTitle
		}

		var alphabet string = caesardisk.Alpha_EN
		if flgPunct {
			alphabet = caesardisk.Alpha_PU
		} else if flgES {
			alphabet = caesardisk.Alpha_ES
		}

		// III. Execute
		generateWheel := func(alpha, filename string, inner bool, options caesardisk.CaesarWheelOptions) {
			if err := caesardisk.GenerateCaesarWheel(alpha, filename, inner, options); err != nil {
				fmt.Println(err)
			}
		}

		generateWheel(alphabet, "caesar_disk_outer.png", false, Options)
		generateWheel(alphabet, "caesar_disk_inner.png", true, Options)
		fmt.Print(Options)
		fmt.Printf("%15s: %s\n", "Alphabet", alphabet)
	}

	caesardisk.BuyMeCoffee()
}
