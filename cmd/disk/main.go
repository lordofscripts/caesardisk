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

const (
	LANG_CZ string = "CZ" // Czech
	LANG_DE string = "DE" // Deutsch (German)
	LANG_EN string = "EN" // English
	LANG_ES string = "ES" // Español (Spanish)
	LANG_GR string = "GR" // Greek (Ελληνικά)
	LANG_IT string = "IT" // Italiano
	LANG_PT string = "PT" // Português
	LANG_RU string = "RU" // Russian (русский) (кириллица)
)

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
	var flgHelp, flgES, flgRU, flgPT, flgDE, flgGR, flgIT, flgCZ, flgPunct bool
	var flgTextFontPath, flgDigitFontPath, flgAlphabet, flgTitle string
	flag.Usage = Usage
	flag.BoolVar(&flgHelp, "help", false, "This help")
	// 1.1 flags for supported preset languages
	flag.BoolVar(&flgES, LANG_ES, false, "Spanish alphabet (overrides -alpha)")
	flag.BoolVar(&flgIT, LANG_IT, false, "Italian alphabet (overrides -alpha)")
	flag.BoolVar(&flgPT, LANG_PT, false, "Portuguese alphabet (overrides -alpha)")
	flag.BoolVar(&flgDE, LANG_DE, false, "German alphabet (overrides -alpha)")
	flag.BoolVar(&flgRU, LANG_RU, false, "Cyrillic alphabet (overrides -alpha)")
	flag.BoolVar(&flgGR, LANG_GR, false, "Greek alphabet (overrides -alpha)")
	flag.BoolVar(&flgCZ, LANG_CZ, false, "Czech alphabet (overrides -alpha)")
	// 1.2 flags for preset full punctuation, symbols, numbers and space (auxillary disk)
	flag.BoolVar(&flgPunct, "PU", false, "Punctuation and numerical alphabet (overrides -alpha)")
	// 1.3 flag for custom alphabet
	flag.StringVar(&flgAlphabet, "alpha", "", "Alphabet defaults to English ASCII alphabet")
	// 1.4 flags for output formatting
	flag.StringVar(&flgTitle, "title", "", "Title (usually disk language or ID)")
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
		untitled := true
		var langSuffix string
		if len(flgTitle) != 0 {
			untitled = false
		}

		var alphabet string = caesardisk.Alpha_EN
		switch {
		case flgPunct:
			alphabet = caesardisk.Alpha_PU
			langSuffix = "PU"

		case flgES:
			alphabet = caesardisk.Alpha_ES
			if untitled {
				flgTitle = "Español"
			}
			langSuffix = LANG_ES

		case flgIT:
			alphabet = caesardisk.Alpha_IT
			if untitled {
				flgTitle = "Italiano"
			}
			langSuffix = LANG_IT

		case flgPT:
			alphabet = caesardisk.Alpha_PT
			if untitled {
				flgTitle = "Português"
			}
			langSuffix = LANG_PT

		case flgDE:
			alphabet = caesardisk.Alpha_DE
			if untitled {
				flgTitle = "Deutsch"
			}
			langSuffix = LANG_DE

		case flgCZ:
			alphabet = caesardisk.Alpha_CZ
			if untitled {
				flgTitle = "Czech"
			}
			langSuffix = LANG_CZ

		case flgGR:
			alphabet = caesardisk.Alpha_GR
			if untitled {
				flgTitle = "Greek"
			}
			langSuffix = LANG_GR

		case flgRU:
			alphabet = caesardisk.Alpha_RU
			if untitled {
				flgTitle = "Cyrillic"
			}
			langSuffix = LANG_RU

		case len(flgAlphabet) != 0:
			alphabet = flgAlphabet // custom alphabet
			langSuffix = ""

		default:
			alphabet = caesardisk.Alpha_EN
			langSuffix = LANG_EN
		}

		if len(flgTitle) != 0 {
			flgTitle := strings.Join(strings.Split(flgTitle, ""), " ")
			Options.Title = flgTitle
		}

		// III. Execute
		generateFilename := func(basename string, isInner bool, suffix string) string {
			if len(suffix) != 0 {
				basename = basename + "_" + suffix
			}
			if isInner {
				basename = basename + "_inner.png"
			} else {
				basename = basename + "_outer.png"
			}
			return basename
		}

		generateWheel := func(alpha, filename string, inner bool, options caesardisk.CaesarWheelOptions) {
			filename = generateFilename(filename, inner, langSuffix)
			if err := caesardisk.GenerateCaesarWheel(alpha, filename, inner, options); err != nil {
				fmt.Println(err)
			}
		}

		// the filename is in base form
		generateWheel(alphabet, "caesar_disk", false, Options)
		generateWheel(alphabet, "caesar_disk", true, Options)
		fmt.Print(Options)
		fmt.Printf("%15s: %s\n", "Alphabet", alphabet)
	}

	caesardisk.BuyMeCoffee()
}
