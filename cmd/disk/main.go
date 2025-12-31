/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   goCaesarDisk
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

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

var (
	ErrDualNotSupported error = errors.New("the -dual option is only valid for English & Spanish")
)

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

// Death of an application by outputting a good-bye and setting
// the OS exit code. It is logged as fatal.
func Die(message string, exitCode int) {
	fmt.Println("\n", "\t💀 x 💀 x 💀\n\t", message, "\n\tExit code: ", exitCode)
	os.Exit(exitCode)
}

// display the error and die with an exit code, logging it as Fatal.
func DieWithError(err error, exitCode int) {
	fmt.Println("\n", "\t💀 x 💀 x 💀\n\t", err.Error(), "\n\tExit code: ", exitCode)
	os.Exit(exitCode)
}

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
	var flgHelp, flgES, flgRU, flgPT, flgDE, flgGR, flgIT, flgCZ, flgPunct, flgDual bool
	var flgTextFontPath, flgDigitFontPath, flgAlphabet, flgTitle string
	var flgAssemble int
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
	flag.BoolVar(&flgDual, "dual", false, "Dual alphabet disk (only for -ES and -EN)")
	// 1.2 flags for preset full punctuation, symbols, numbers and space (auxillary disk)
	flag.BoolVar(&flgPunct, "PU", false, "Punctuation and numerical alphabet (overrides -alpha)")
	// 1.3 flag for custom alphabet
	flag.StringVar(&flgAlphabet, "alpha", "", "Alphabet defaults to English ASCII alphabet")
	// 1.4 flags for output formatting
	flag.StringVar(&flgTitle, "title", "", "Title (usually disk language or ID)")
	flag.StringVar(&flgTextFontPath, "text-font", "", "Text font path")
	flag.StringVar(&flgDigitFontPath, "digit-font", "", "Digit font path (else use same text font)")
	flag.IntVar(&flgAssemble, "assemble", -1, "when set also output a final disk with key N")
	flag.Parse()

	// II. Command-line flag validation and processing
	caesardisk.Copyright(caesardisk.CO1, true)

	if flgHelp {
		flag.Usage()
	} else {
		var Options caesardisk.CaesarWheelOptions = caesardisk.DefaultCaesarWheelOptions
		Options.LetterColorAlt = caesardisk.NewRGBFromString("#c13e93") // for inner
		Options.DigitsSize += 2.0

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

		alphabetLet := caesardisk.Alpha_EN
		alphabetPun := "" // only used in dual-disk print-out
		switch {
		case flgPunct:
			alphabetLet = caesardisk.Alpha_PU
			langSuffix = "PU"

		case flgES:
			alphabetLet = caesardisk.Alpha_ES
			if untitled {
				flgTitle = "Español"
			}
			langSuffix = LANG_ES
			if flgDual { // Spanish -ES supports -dual
				alphabetLet = caesardisk.Alpha_ES_DUAL
				alphabetPun = caesardisk.Alpha_PU_DUAL_ES
			}

		case flgIT:
			alphabetLet = caesardisk.Alpha_IT
			if untitled {
				flgTitle = "Italiano"
			}
			langSuffix = LANG_IT

		case flgPT:
			alphabetLet = caesardisk.Alpha_PT
			if untitled {
				flgTitle = "Português"
			}
			langSuffix = LANG_PT

		case flgDE:
			alphabetLet = caesardisk.Alpha_DE
			if untitled {
				flgTitle = "Deutsch"
			}
			langSuffix = LANG_DE

		case flgCZ:
			alphabetLet = caesardisk.Alpha_CZ
			if untitled {
				flgTitle = "Czech"
			}
			langSuffix = LANG_CZ

		case flgGR:
			alphabetLet = caesardisk.Alpha_GR
			if untitled {
				flgTitle = "Greek"
			}
			langSuffix = LANG_GR

		case flgRU:
			alphabetLet = caesardisk.Alpha_RU
			if untitled {
				flgTitle = "Cyrillic"
			}
			langSuffix = LANG_RU

		case len(flgAlphabet) != 0:
			alphabetLet = flgAlphabet // custom alphabet
			langSuffix = ""

		default:
			alphabetLet = caesardisk.Alpha_EN
			langSuffix = LANG_EN
			if flgDual { // English (default) -EN supports -dual
				alphabetPun = caesardisk.Alpha_PU_DUAL_EN
			}
		}

		if flgDual && len(alphabetPun) == 0 {
			DieWithError(ErrDualNotSupported, 1)
		}

		if len(flgTitle) != 0 {
			flgTitle := strings.Join(strings.Split(flgTitle, ""), " ")
			Options.Title = flgTitle
		}

		// III. Execute
		generateFilename := func(basename string, isInner bool, suffix string) string {
			// the ISO language code
			if len(suffix) != 0 {
				basename = basename + "_" + suffix
			}
			// whether it is a Dual version, else single
			if flgDual {
				basename = basename + "_dual"
			}
			// Inner vs. Outer disk ID
			if isInner {
				basename = basename + "_inner.png"
			} else {
				basename = basename + "_outer.png"
			}
			return basename
		}

		generateWheel := func(alpha, filename string, inner, dual bool, options caesardisk.CaesarWheelOptions) {
			filename = generateFilename(filename, inner, langSuffix)
			var err error = nil
			if !dual {
				err = caesardisk.GenerateCaesarWheel(alpha, filename, inner, options)
			} else {
				err = caesardisk.GenerateDualCaesarWheel(alpha, alphabetPun, filename, inner, options)
			}

			if err != nil {
				DieWithError(err, 2)
			}
		}

		// the filename is in base form
		generateWheel(alphabetLet, "caesar_disk", false, flgDual, Options)
		generateWheel(alphabetLet, "caesar_disk", true, flgDual, Options)

		// -assemble
		if flgAssemble > -1 {
			// how many unicode characters in the chosen alphabet
			maxKey := utf8.RuneCountInString(alphabetLet)
			// adjust the shift if necessary, 0..N-1
			shift := flgAssemble % maxKey
			outerDiskFile := generateFilename("caesar_disk", false, langSuffix)
			innerDiskFile := generateFilename("caesar_disk", true, langSuffix)
			finalFile := "caesar_disk_" + langSuffix + "_final.png"
			err := caesardisk.SuperimposeDisksByShift(
				shift, maxKey,
				outerDiskFile,
				innerDiskFile,
				finalFile,
				flgDual,
				Options)
			if err != nil {
				println("WARNING", err)
			} else {
				fmt.Printf("successfully generated '%s' with Caesar shift %d: ", finalFile, shift)
			}
		}

		fmt.Print('\n', Options)
		fmt.Printf("%15s: %s\n", "Alphabet", alphabetLet)
		if flgDual {
			fmt.Printf("%15s: %s\n", "Symbols ", alphabetPun)
		}
	}

	caesardisk.BuyMeCoffee()
}
