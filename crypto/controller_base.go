/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package crypto

import (
	"fmt"
	"sync"

	"github.com/lordofscripts/caesardisk"
	"github.com/lordofscripts/goapp/app/logx"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

// The View that presents the main key must implement this interface.
// The SessionController consumes this service to notify changes to
// the GUI that are cascaded as a response to another change.
type ICaesarKeyService interface {
	// get the current Caesar key as letter (alphabet-dependent) and shift
	GetKey() (rune, int)
	SetKey(rune, int)

	Update()
}

// The View that presents the offset & alternate key must implement this
// interface. SessionController consumes this service to notify of cascaded
// changes in response to a request.
type IKeyOffsetService interface {
	// get the current RAW/uncorrected Caesar key offset
	GetOffset() int
	// get the alternate key (letter & shift) derived from the corrected offset
	GetAltKey() (rune, int)
	SetOffsetRune(rune)
	SetOffset(rune, int)

	Hide()
	Show()
	Reset()
	Update()
}

// When the SessionController responds to changes in the crypto parameters
// that affect the visual wheel, we must request a wheel update.
type IWheelUpdateService interface {
	Update()
	CanShowOffset(bool)
	IsShowingMain() bool
}

// The view should implement this interface to pop up a dialog when
// the controller needs to notify of changes that were made that
// require the attention of the user.
type IViewNotifier interface {
	// the 1st parameter is true for an error, false for informational.
	// the 2nd parameter is the message to be displayed
	Notify(bool, string)
}

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

// All controllers should implement this base object and extend it
type ControllerBase struct {
	mutexRW    sync.Mutex
	viewNotify IViewNotifier
}

type CaesarKey struct {
	Letter rune
	Shift  int
}

// SessionModel contains all the data needed to perform encryption/decryption
// with any of the supported Caesar modes.
type SessionModel struct {
	Mode    CaesarCipherMode // Caesar cipher mode (Caesar|Didimus|Fibonacci|Primus)
	Alpha   caesardisk.AlphabetModel
	MainKey CaesarKey // Main Caesar key (required)
	AltKey  CaesarKey // Alternate Caesar key (for Didimus only)
	Offset  int       // Offset (for Didimus) from which alternate key is derived
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// implement fmt.Stringer. Returns key letter and shift
func (ck CaesarKey) String() string {
	return fmt.Sprintf("('%c',%02d)", ck.Letter, ck.Shift)
}

// implement fmt.Stringer. Returns mode, main key, offset, alternate key
func (sm SessionModel) String() string {
	return fmt.Sprintf("%s M:%s OFS:%02d => A:%s", sm.Mode, sm.MainKey, sm.Offset, sm.AltKey)
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

// if we have a notifier on the view do invoke it
func (cb *ControllerBase) notifyError(format string, args ...any) {
	logx.Enter()
	defer logx.Leave()

	if cb.viewNotify != nil {
		cb.viewNotify.Notify(true, fmt.Sprintf(format, args...))
	}
}

// if we have a notifier on the view do invoke it
func (cb *ControllerBase) notifyInfo(format string, args ...any) {
	logx.Enter()
	defer logx.Leave()

	if cb.viewNotify != nil {
		cb.viewNotify.Notify(false, fmt.Sprintf(format, args...))
	}
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/
