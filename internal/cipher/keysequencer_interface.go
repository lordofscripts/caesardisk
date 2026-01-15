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
)

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

// Allows for key schedules in a message. Plain caesar is monosyllabic
// and thus has no schedule (same key always). Didimus & Fibonacci are
// polisyllabic and have their own key schedules over a single message.
type IKeySequencer interface {
	// implements fmt.Stringer and returns the name of the Cipher
	fmt.Stringer
	// Validates the key parameters. Must be called after the
	// constructor.
	Validate() error
	// Query the key value range (min,max) for main or alt key
	KeyRange() (min, max int)
	// Retrieve cipher parameters
	GetParams() *CaesarParameters
	// gets the key to use for the current character. Must only be
	// called if the character to encode/decode is part of the
	// encoding alphabet.
	NextKey() int
	// whether the sequencer is polialphabetic or not
	IsPolyalphabetic() bool
	// whether the Offset parameter is used in key sequencing
	IsOffsetRequired() bool
}
