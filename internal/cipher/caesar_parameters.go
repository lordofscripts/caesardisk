/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package cipher

import (
	"fmt"

	"github.com/lordofscripts/caesardisk"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type CaesarParameters struct {
	Alphabet *caesardisk.AlphabetModel
	KeyValue int
}

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// (ctor) a new instance of Caesar cipher parameters for the
// chosen alphabet.
func NewCaesarParameters(alphabet *caesardisk.AlphabetModel) *CaesarParameters {
	return &CaesarParameters{
		Alphabet: alphabet,
		KeyValue: 0,
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer and returns the current key in character and
// shift value together with the actual alphabet letters or symbols.
func (c *CaesarParameters) String() string {
	keyLetter, _ := c.Alphabet.Character(c.KeyValue)
	return fmt.Sprintf("(%02d|%c) %s", c.KeyValue, keyLetter, c.Alphabet)
}
