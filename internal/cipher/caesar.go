/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *						goCaesarDisk GUI
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Implements basic Caesar cipher encoding/decoding using Unicode
 * (foreign) alphabets (not just ASCII).
 *-----------------------------------------------------------------*/
package cipher

import (
	"log"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/lordofscripts/caesardisk/internal/hash"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

var hashSeed uint64 = 0xDEADBEA7

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type Caesar struct {
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// (ctor) new instance of plain Caesar cipher
func NewCaesarCipher() *Caesar {
	return &Caesar{}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// encode a message and package it in a "standard" form. The
// standard PDU format is {TIMESTAMP}{CHECKSUM}{PAYLOAD} where
// Payload is the encrypted message string, Checksum is the checksum
// over the payload, and Timestamp is the the form YYYYMMDDTHHMMSS
func (c *Caesar) EncodeMessage(plain string, params *CaesarParameters) (string, error) {
	payload := c.Encode(plain, params)
	msgPDU := NewCaesarMessage(hash.NewXXH64(hashSeed))
	msgPDU.AddMessage(payload)

	return msgPDU.String(), nil
}

// decode a message that is in "standard" format
func (c *Caesar) DecodeMessage(cipheredMessage string, params *CaesarParameters) (string, error) {
	check := hash.NewXXH64(hashSeed)
	if payload, err := VerifyCaesarMessage(check, cipheredMessage); err != nil {
		return "", err
	} else {
		return c.Decode(payload, params), nil
	}
}

func (c *Caesar) Encode(plain string, params *CaesarParameters) string {
	var result strings.Builder
	var alphabet string = params.Alphabet.String()
	var shift int
	tokens := []rune(alphabet)
	shift = params.KeyValue

	if len(tokens) < params.KeyValue-1 {
		log.Printf("wrapping key '%d' of alphabet length %d", params.KeyValue, len(tokens))
		shift %= len(tokens)
	}

	tabulaRaw := RotateStringLeft(alphabet, shift)
	tabulaOut := []rune(tabulaRaw)

	for _, plainRune := range []rune(plain) {
		// the reference alphabet is Uppercase, input may be lowercase
		isLower := unicode.IsLower(plainRune)
		if isLower {
			plainRune = unicode.ToUpper(plainRune)
		}

		if at := slices.Index(tokens, plainRune); at != -1 {
			ciphered := tabulaOut[at]
			// maintain case
			if isLower {
				ciphered = unicode.ToLower(ciphered)
			}

			result.WriteRune(ciphered)
		} else {
			result.WriteRune(plainRune)
		}
	}

	return result.String()
}

func (c *Caesar) Decode(ciphered string, params *CaesarParameters) string {
	var result strings.Builder
	var alphabet string = params.Alphabet.String()
	var shift int

	tokens := []rune(alphabet)
	shift = params.KeyValue
	if len(tokens) < params.KeyValue-1 {
		log.Printf("wrapping key '%d' of alphabet length %d", params.KeyValue, len(tokens))
		shift %= len(tokens)
	}

	tabulaRaw := RotateStringLeft(alphabet, shift)
	tabulaOut := []rune(tabulaRaw)

	for _, cipherRune := range []rune(ciphered) {
		// the reference alphabet is Uppercase, input may be lowercase
		isLower := unicode.IsLower(cipherRune)
		if isLower {
			cipherRune = unicode.ToUpper(cipherRune)
		}

		if at := slices.Index(tabulaOut, cipherRune); at != -1 {
			plain := tokens[at]
			// maintain case
			if isLower {
				plain = unicode.ToLower(plain)
			}

			result.WriteRune(plain)
		} else {
			result.WriteRune(cipherRune)
		}
	}

	return result.String()
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

func RotateStringLeft(s string, shift int) string {
	complementShift := utf8.RuneCountInString(s) - shift
	return RotateStringRight(s, complementShift)
}

func RotateStringRight(s string, shift int) string {
	alphaSize := utf8.RuneCountInString(s)
	// only rotate right
	if shift < 0 {
		shift = shift * -1
	}
	// bigger than alphabet? then wrap it
	if shift > alphaSize {
		shift = shift % alphaSize
	}
	// no rotation?
	if shift == alphaSize || shift == 0 {
		return s
	}
	// rotate
	//return s[alphaSize-shift:] + s[0:alphaSize-shift]
	runic := []rune(s)
	result := string(runic[alphaSize-shift:])
	result += string(runic[0 : alphaSize-shift])
	return result
}
