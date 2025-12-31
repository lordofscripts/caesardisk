/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *						goCaesarDisk GUI
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A Caesar Message PDU (Protocol Data Unit) consists of the encrypted
 * payload string, prepended with {TIMESTAMP}{CHECKSUM} where the
 * Timestamp is the standard YYYYMMDDTHHMMSS and the Checksum is the
 * XXHash64 checksum over the entire Payload.
 *-----------------------------------------------------------------*/
package cipher

import (
	"fmt"
	"strings"
	"time"

	"github.com/lordofscripts/caesardisk/internal/hash"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type CaesarMessage struct {
	hasher  hash.IDigest
	payload string
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// (ctor) A ciphered Caesar message that will be formatted with
// date/time and a user-chosen hash or checksum
func NewCaesarMessage(digest hash.IDigest) *CaesarMessage {
	return &CaesarMessage{
		hasher:  digest,
		payload: "",
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer returning the encapsulated Caesar cipher
// message prepended with the date/time YYYYMMDDTHHMMSS and the
// message hash
func (m *CaesarMessage) String() string {
	dateStr := time.Now().Format("2006-01-02T15-04-05")
	dateStr = strings.ReplaceAll(dateStr, "-", "")
	return fmt.Sprintf("%s%s%s", dateStr, m.hasher.String(), m.payload)
}

// used when packaging a new Caesar message by adding new
// cipher data.
func (m *CaesarMessage) AddMessage(ciphered string) {
	data := []byte(ciphered)
	if m.hasher != nil {
		m.hasher.Update(data)
		m.payload += ciphered
	}
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

// Verify that the packaged Caesar message has not been corrupted.
// on success returns the payload (actual cipher message) and nil.
func VerifyCaesarMessage(hasher hash.IDigest, packet string) (string, error) {
	const EMPTY = ""

	if len(packet) == 0 {
		return EMPTY, fmt.Errorf("packet is empty")
	}
	if hasher == nil {
		return EMPTY, fmt.Errorf("no hasher defined for verification")
	}
	const TIMESTAMP_LEN = 15
	sizeH := hasher.Length() * 2 // length as Hex
	if len(packet)-TIMESTAMP_LEN < sizeH {
		return EMPTY, fmt.Errorf("corrupted Caesar packet")
	}
	if len(packet)-TIMESTAMP_LEN == sizeH {
		return EMPTY, fmt.Errorf("corrupted Caesar packet has no message")
	}

	hashStr := packet[TIMESTAMP_LEN : TIMESTAMP_LEN+sizeH]
	payloadStr := packet[TIMESTAMP_LEN+sizeH:]

	hasher.Update([]byte(payloadStr))
	if !strings.EqualFold(hashStr, hasher.String()) {
		return EMPTY, fmt.Errorf("ciphered message is altered %s != %s", hashStr, hasher.String())
	}

	return payloadStr, nil
}
