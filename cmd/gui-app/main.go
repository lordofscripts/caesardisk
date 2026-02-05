/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package main

import (
	"flag"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"github.com/lordofscripts/caesardisk"

	"github.com/lordofscripts/caesardisk/cmd/gui-app/gui"
	"github.com/lordofscripts/goapp/app/logx"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

var (
	APP_NAME                                             = "caesardisk-gui"
	APP_GROUP                                            = "coralys"
	MyAlphabets     map[string]*caesardisk.AlphabetModel = make(map[string]*caesardisk.AlphabetModel)
	DefaultAlphabet *caesardisk.AlphabetModel
	WheelOptions    *caesardisk.CaesarWheelOptions
)

/* ----------------------------------------------------------------
 *						L o c a l s
 *-----------------------------------------------------------------*/

const (
	defaultAlphabetKey = "English"
)

var ()

/* ----------------------------------------------------------------
 *				I n i t i a l i z e r
 *-----------------------------------------------------------------*/
func init() {
	logx.SingLogGate = logx.GetLogGateInstance(APP_NAME, APP_GROUP)
	logx.SingLogGate.LoadFilters()
	logx.SingLogGate.WithCallTree("/tmp/caesardisk-calltree.txt")

	MyAlphabets["English"] = caesardisk.NewAlphabetModelCased(caesardisk.Alpha_EN)
	MyAlphabets["Español con acentos"] = caesardisk.NewAlphabetModelCased(caesardisk.Alpha_ES)
	MyAlphabets["Czech"] = caesardisk.NewAlphabetModelCased(caesardisk.Alpha_CZ)
	MyAlphabets["Español"] = caesardisk.NewAlphabetModelCased(caesardisk.Alpha_ES_DUAL)
	MyAlphabets["Deutsch"] = caesardisk.NewAlphabetModelCased(caesardisk.Alpha_DE)
	MyAlphabets["Italiano"] = caesardisk.NewAlphabetModelCased(caesardisk.Alpha_IT)
	MyAlphabets["Português"] = caesardisk.NewAlphabetModelCased(caesardisk.Alpha_PT)
	MyAlphabets["Russian"] = caesardisk.NewAlphabetModelCased(caesardisk.Alpha_RU)
	MyAlphabets["Greek"] = caesardisk.NewAlphabetModelCased(caesardisk.Alpha_GR)
	MyAlphabets["Punctuation (all)"] = caesardisk.NewAlphabetModelCased(caesardisk.Alpha_PU)
	MyAlphabets["Puntuacion para Español"] = caesardisk.NewAlphabetModelCased(caesardisk.Alpha_PU_DUAL_ES)
	MyAlphabets["Punctuation for English"] = caesardisk.NewAlphabetModelCased(caesardisk.Alpha_PU_DUAL_EN)
	MyAlphabets["Runes"] = caesardisk.NewAlphabetModelCased(caesardisk.Alpha_RUNES)

	DefaultAlphabet = MyAlphabets[defaultAlphabetKey]

	const OUTER_ALPHABET_COLOR = "#000000"
	const INNER_ALPHABET_COLOR = "#ff4538"
	const DIGITS_COLOR = "#0000f4"
	WheelOptions = &caesardisk.DefaultCaesarWheelOptions
	//WheelOptions.SetLetterColors("#ff388e", "#ff4538").SetDigitColor("#0000f4")
	WheelOptions.SetLetterColors(OUTER_ALPHABET_COLOR, INNER_ALPHABET_COLOR).SetDigitColor(DIGITS_COLOR)
	WheelOptions.DigitsSize += 2
	//WheelOptions.Orthogonal = false
}

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

func RuneString(latin string) string {
	const (
		LOOKUP_STD string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		RUNES      string = "\u16ab\u16d2\u16b3\u16de\u16d6\u16a0\u16b7\u16bb\u16c1\u16c3\u16f1\u16da\u16d7\u16be\u16a9\u16c8\u16e9\u16b1\u16cb\u16cf\u16a2\u16a1\u16b9\u16ea\u16e6\u16ce"
	)

	chars := []rune(strings.ToUpper(latin))
	runesLookup := []rune(RUNES)
	result := make([]rune, len(chars))

	for index, char := range chars {
		if strings.ContainsRune(LOOKUP_STD, char) {
			at := strings.IndexRune(LOOKUP_STD, char)
			result[index] = runesLookup[at]
		} else {
			result[index] = char
		}
	}

	return string(result)
}

/* ----------------------------------------------------------------
 *					M A I N    |     D E M O
 *-----------------------------------------------------------------*/

func main() {
	var flgLogMode string
	flag.StringVar(&flgLogMode, "log", "", "Log output to: none|term|DIRECTORY")
	flag.Parse()
	logx.SingLogGate.Setup(flgLogMode)
	defer logx.SingLogGate.Close()

	if strings.ToLower(flgLogMode) != "none" {
		logx.Print(caesardisk.Version)
		logx.Print(caesardisk.BuildInfo())
	}

	caesardisk.Copyright(caesardisk.CO1)

	catalogPtr := (*gui.AlphabetList)(&MyAlphabets)
	myGui := gui.NewGUI(catalogPtr, WheelOptions).Define().Bind()
	myGui.Run()

	caesardisk.BuyMeCoffee("lostinwriting")
}
