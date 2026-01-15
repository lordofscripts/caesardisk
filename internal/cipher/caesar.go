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

// An interface for all text/cipher transforms
type ITranscoder interface {
	Encode(string) string
	Decode(string) string
}

var _ ITranscoder = (*Caesar)(nil)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type Caesar struct {
	sequencer IKeySequencer
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// (ctor) to create a Caesar substitution cipher handler from a
// preset key sequencer that must have previously been validated.
func NewCaesarCipherFromSequencer(seq IKeySequencer) *Caesar {
	return &Caesar{
		sequencer: seq,
	}
}

// (ctor) new instance of plain Caesar cipher.
// Caesar (plain) is a monoalphabetic substitution cipher.
// Rot13 is nothing but Caesar with (main) key=13.
// Note: the config.Offset is not used.
func NewCaesarCipher(config *CaesarParameters) *Caesar {
	return &Caesar{
		sequencer: NewCaesarSequencer(config),
	}
}

// (ctor) new instance of Didimus cipher. Didimus is a
// variant of Caesar that uses an offset over the main Caesar
// key. Being a bi-alphabetic substitution cipher, it alternates
// encoding using the main key for even characters and the
// alternate key (main key + offset % ALPHA_LEN) for odd characters.
// Note: the config.Offset is always used.
func NewDidimusCipher(config *CaesarParameters) *Caesar {
	return &Caesar{
		sequencer: NewDidimusSequencer(config),
	}
}

// (ctor) Fibonacci is a polyalphabetic substitution cipher based
// on Caesar, but it uses a 10-term Fibonacci series to generate
// alternate keys that are used for every other character.
// Note: the config.Offset is not used.
func NewFibonacciCipher(config *CaesarParameters) *Caesar {
	return &Caesar{
		sequencer: NewFibonacciSequencer(config),
	}
}

// (ctor) Primus is a polyalphabetic substitution cipher based
// on Caesar. It is similar to Fibonacci with two differences.
// It uses (up to) the first 10 Prime numbers to generate the
// alternate keys instead of a Fibonacci series. The Offset (0..10)
// is used to determine how many prime numbers are used prior to
// rewinding. In every iteration the first is the main Caesar key,
// and the rest uses the main key plus the prime number modulo
// length of the alphabet.
// Note: the config.Offset is used.
func NewPrimusCipher(config *CaesarParameters) *Caesar {
	return &Caesar{
		sequencer: NewPrimusSequencer(config),
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer and returns algorithm parameters
func (c *Caesar) String() string {
	return c.sequencer.String()
}

// encode a message and package it in a "standard" form. The
// standard PDU format is {TIMESTAMP}{CHECKSUM}{PAYLOAD} where
// Payload is the encrypted message string, Checksum is the checksum
// over the payload, and Timestamp is the the form YYYYMMDDTHHMMSS
func (c *Caesar) EncodeMessage(plain string) (string, error) {
	payload := c.Encode(plain)
	msgPDU := NewCaesarMessage(hash.NewXXH64(hashSeed))
	msgPDU.AddMessage(payload)

	return msgPDU.String(), nil
}

// decode a message that is in "standard" format
func (c *Caesar) DecodeMessage(cipheredMessage string) (string, error) {
	check := hash.NewXXH64(hashSeed)
	if payload, err := VerifyCaesarMessage(check, cipheredMessage); err != nil {
		return "", err
	} else {
		return c.Decode(payload), nil
	}
}

// implements ITranscoder for encoding/encrypting a message using
// the selected Caesar mode/variant (Caesar, Didimus, Fibonacci, Primus)
func (c *Caesar) Encode(plain string) string {
	var result strings.Builder
	var alphabet string = c.sequencer.GetParams().Alphabet.String()

	var tabulaIn []rune  // the plain-text alphabet tabula
	var tabulaRaw string // the raw ciphered alphabet string
	var tabulaOut []rune // the ciphered alphabet tabula
	if !c.sequencer.IsPolyalphabetic() {
		withKey := c.sequencer.GetParams().KeyValue
		tabulaRaw = RotateStringLeft(alphabet, withKey)
		tabulaOut = []rune(tabulaRaw)
	}
	// during encryption we map from tabulaIn to tabulaOut
	tabulaIn = []rune(alphabet)

	// · iterate through each of the plain-text Unicode characters in the input string
	for _, plainRune := range []rune(plain) {
		// · Letter-case preservation
		//   the reference alphabet is Uppercase, input may be lowercase
		isLower := unicode.IsLower(plainRune)
		if isLower {
			plainRune = unicode.ToUpper(plainRune)
		}

		if at := slices.Index(tabulaIn, plainRune); at != -1 {
			// · The Unicode point CAN be encoded (present in alphabet)
			//	 select the appropriate key & tabula for polialphabetic ciphers
			withKey := c.sequencer.NextKey()
			if c.sequencer.IsPolyalphabetic() {
				tabulaRaw = RotateStringLeft(alphabet, withKey)
				tabulaOut = []rune(tabulaRaw)
			}

			// · map from tabulaIn (withKey) to tabulaOut content
			ciphered := tabulaOut[at]
			// · preserve case on output
			if isLower {
				ciphered = unicode.ToLower(ciphered)
			}

			// · write the encrypted Unicode character
			result.WriteRune(ciphered)
		} else {
			// · The Unicode point CANNOT be encoded (not present in alphabet)
			//	 pass it through as-is.
			result.WriteRune(plainRune)
		}
	}

	return result.String()
}

// implements ITranscoder for decoding/decrypting a message using
// the selected Caesar mode/variant (Caesar, Didimus, Fibonacci, Primus)
func (c *Caesar) Decode(ciphered string) string {
	var result strings.Builder
	var alphabet string = c.sequencer.GetParams().Alphabet.String()
	var tabulaIn []rune  // the plain-text alphabet tabula
	var tabulaRaw string // the raw ciphered alphabet string
	var tabulaOut []rune // the ciphered alphabet tabula

	// · slight optimization for monoalphabetic modes
	if !c.sequencer.IsPolyalphabetic() {
		withKey := c.sequencer.GetParams().KeyValue
		tabulaRaw = RotateStringLeft(alphabet, withKey)
		tabulaOut = []rune(tabulaRaw)
	}
	// during decryption we map from tabulaOut to tabulaIn
	tabulaIn = []rune(alphabet)

	for _, cipherRune := range []rune(ciphered) {
		// · Letter-case preservation
		// the reference alphabet is Uppercase, input may be lowercase
		isLower := unicode.IsLower(cipherRune)
		if isLower {
			cipherRune = unicode.ToUpper(cipherRune)
		}

		// we cannot use at just yet because for polyalphabetic
		// we have not yet (re)constructed tabulaOut
		if at := slices.Index(tabulaIn, cipherRune); at != -1 {
			// · The Unicode point CAN be decoded (present in alphabet)
			//	 select the appropriate key & tabula for polialphabetic ciphers
			withKey := c.sequencer.NextKey()
			if c.sequencer.IsPolyalphabetic() {
				tabulaRaw = RotateStringLeft(alphabet, withKey)
				tabulaOut = []rune(tabulaRaw)
			}

			// · Now that we have the ciphered tabula, we can determine index.
			//   This works because both tabulaIn & tabulaOut contain the same
			//	 character set, i.e. not transliteration of text to symbols.
			//	 Therefore, the following value will always be >= 0
			at = slices.Index(tabulaOut, cipherRune)
			plain := tabulaIn[at]
			// · preserve the input letter-case
			if isLower {
				plain = unicode.ToLower(plain)
			}

			// · Write to output decoded string
			result.WriteRune(plain)
		} else {
			// · If not present, pass as-is to the output string
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
