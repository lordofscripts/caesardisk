/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							goCaesarDisk
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package caesardisk

import (
	"errors"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	No TriState = iota
	Yes
	Unknown
)

var (
	ErrCharacterNotFound        error = errors.New("character not found in alphabet")
	ErrCharacterIndexOutOfRange error = errors.New("character index greater than alphabet")
)

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

// a tri-state value (yes, no, maybe)
type TriState uint8

type AlphabetModel struct {
	Name        string
	alphabet    []rune
	upperCased  TriState
	symbolsOnly bool
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// (ctor) an instance of an alphabet made of the given string. It
// assumes all characters are unique and no trimming is done!
// However, the alphabet is converted to uppercase.
func NewAlphabetModelCased(alphabet string) *AlphabetModel {
	return &AlphabetModel{

		alphabet:    []rune(strings.ToUpper(alphabet)),
		upperCased:  Yes,
		symbolsOnly: false,
	}
}

// (ctor) alphabet model without case conversion
func NewAlphabetModel(alphabet string) *AlphabetModel {
	return &AlphabetModel{
		alphabet:    []rune(alphabet),
		upperCased:  Unknown,
		symbolsOnly: false,
	}
}

// (ctor) A symbols/punctuation-only alphabet without letters
// that can be upper/lowercased.
func NewAlphabetModelForSymbols(alphabet string) *AlphabetModel {
	return &AlphabetModel{
		alphabet:    []rune(alphabet),
		upperCased:  No,
		symbolsOnly: true,
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer and returns the alphabet string
func (a *AlphabetModel) String() string {
	return string(a.alphabet)
}

// the length or size of the alphabet in Unicode runes (not bytes)
func (a *AlphabetModel) Length() int {
	return len(a.alphabet)
}

// whether the alphabet has multi-byte characters
func (a *AlphabetModel) IsMultiByte() bool {
	return len(a.alphabet) != utf8.RuneCountInString(string(a.alphabet))
}

// return the character at the given zero-based index
func (a *AlphabetModel) Character(index int) (rune, error) {
	if index > a.Length()-1 {
		return rune(0), ErrCharacterIndexOutOfRange
	}

	return a.alphabet[index], nil
}

// finds the index of the exact character (case-sensitive).
// return -1 if not found
func (a *AlphabetModel) FindExact(char rune) int {
	return slices.Index(a.alphabet, char)
}

// finds the index of the character (case-insensitive search),
// or -1 if not present.
func (a *AlphabetModel) Find(char rune) int {
	if a.symbolsOnly {
		return a.FindExact(char)
	}

	if unicode.IsLower(char) && a.upperCased == Yes {
		uc := unicode.ToUpper(char)
		return a.FindExact(uc)
	}

	if unicode.IsUpper(char) && a.upperCased == No {
		lc := unicode.ToLower(char)
		return a.FindExact(lc)
	}

	return slices.Index(a.alphabet, char)
}

func (a *AlphabetModel) FirstChar() rune {
	return a.alphabet[0]
}

func (a *AlphabetModel) IsUpper() bool {
	return a.upperCased == Yes
}

func (a *AlphabetModel) IsLower() bool {
	return a.upperCased == No
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/
