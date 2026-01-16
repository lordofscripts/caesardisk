/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Here we define a Caesar parameter storage that is directly bound
 * via data bindings to their respective control widgets in an IGadget.
 *-----------------------------------------------------------------*/
package gui

import (
	"fmt"

	"fyne.io/fyne/v2/data/binding"
	"github.com/lordofscripts/caesardisk"
	"github.com/lordofscripts/caesardisk/crypto"
	"github.com/lordofscripts/caesardisk/internal/cipher"
	"github.com/lordofscripts/goapp/app/logx"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

var (
	// This is a global data binding struct whose members are
	// data-bound to certain widgets of which gadgets are composed.
	DataBindings *dataBinder = &dataBinder{}
	// The data-bound Caesar Key Shift updated by the Key Shift Slider
	// widget found in the KeyShiftGadget. See CipherParams.GetSessionModel()
	BoundKeyShift binding.ExternalFloat = binding.BindFloat(&DataBindings.keyShift)
	// The data-bound Caesar Key Offset updated by the Key Offset Slider
	// widget found in the KeyOffsetGadget. See CipherParams.GetSessionModel()
	BoundKeyOffset binding.ExternalFloat = binding.BindFloat(&DataBindings.keyOffset)
	// The data-bound Caesar Alphabet name updated by the Alphabet Select
	// widget found in the AlphabettGadget. See CipherParams.GetSessionModel()
	BoundAlphaName binding.ExternalString = binding.BindString(&DataBindings.alphaName)
	// The data-bound Caesar Cipher Mode updated by the CipherMode Select
	// widget found in the CipherModeGadget. See CipherParams.GetSessionModel()
	BoundCipherModeName binding.ExternalString = binding.BindString(&DataBindings.modeName)
	// the Use PDU checkbox
	BoundOptionUsePDU binding.ExternalBool = binding.BindBool(&DataBindings.optUsePDU)
)

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

// a raw container of parameters related to Caesar encoding that
// are used for GUI data binding. See GetSessionModel() to retrieve
// the bound values in their natural form.
type dataBinder struct {
	// only for ExternalBind
	alphaName string
	keyShift  float64
	keyOffset float64
	modeName  string

	optUsePDU bool
}

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

func (cp *dataBinder) String() string {
	return fmt.Sprintf("%s (%s) K:%02d O:%02d", cp.alphaName, cp.modeName, int(cp.keyShift), int(cp.keyOffset))
}

// Initializes the values and calls Reload() to refresh the binding
// to reflect these initial values.
func (cp *dataBinder) Init() {
	// cryptographic parameters
	cp.alphaName = InitialAlphabetName
	cp.keyShift = 0
	cp.keyOffset = 0
	cp.modeName = crypto.CaesarMode.String()
	// application options
	cp.optUsePDU = false
}

// Binds the global bound data to listeners
func (cp *dataBinder) Bind() {
	logx.Visit()

	listenKS := binding.NewDataListener(cp.onChangedKeyShift)
	BoundKeyShift.AddListener(listenKS)

	listenKO := binding.NewDataListener(cp.onChangedKeyOffset)
	BoundKeyOffset.AddListener(listenKO)

	listenCM := binding.NewDataListener(cp.onChangedCipherMode)
	BoundCipherModeName.AddListener(listenCM)

	listenAL := binding.NewDataListener(cp.onChangedAlphabet)
	BoundAlphaName.AddListener(listenAL)
}

// Reload should be called if your code has directly changed
// the underlying members instead of using the bounded variable's
// Set() method.
func (cp *dataBinder) Reload() {
	BoundKeyShift.Reload()
	BoundKeyOffset.Reload()
	BoundAlphaName.Reload()
	BoundCipherModeName.Reload()
	BoundOptionUsePDU.Reload()
}

// get a session model based on the bound data that is directly modified
// by the widgets via data binding interfaces.
func (cp *dataBinder) GetSessionModel() crypto.SessionModel {
	logx.Visit()

	// get the full alphabet object
	α := caesardisk.AlphabetFactory(cp.alphaName)

	// get the parameter corrector
	norm := cipher.NewCaesarParameters(α)
	keyShift := norm.SetKey(int(cp.keyShift))
	altKeyShift := norm.SetAltKeyOffset(int(cp.keyOffset))

	// prepare the key descriptors
	var mk, ak crypto.CaesarKey
	mk.Letter, _ = α.Character(int(keyShift))
	mk.Shift = keyShift
	ak.Letter, _ = α.Character(int(altKeyShift))
	ak.Shift = altKeyShift

	// process the cipher mode
	mode, err := crypto.ParseCipherMode(cp.modeName)
	if err != nil {
		logx.Printf("using default Caesar mode due to error: %v", err)
	}

	model := crypto.SessionModel{
		Mode:    mode,
		Alpha:   *α,
		MainKey: mk,
		AltKey:  ak,
		Offset:  int(cp.keyOffset),
	}

	return model
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

func (cp *dataBinder) onChangedKeyShift() {
	logx.OnChanged(int(cp.keyShift))
	logx.Result("bind-key-shift: %02d", int(cp.keyShift))
}

func (cp *dataBinder) onChangedKeyOffset() {
	logx.OnChanged(int(cp.keyOffset))
	logx.Result("bind-key-ofset: %02d", int(cp.keyOffset))
}

func (cp *dataBinder) onChangedAlphabet() {
	logx.OnChanged(cp.alphaName)
	logx.Result("bind-alpha: %s", cp.alphaName)
}

func (cp *dataBinder) onChangedCipherMode() {
	logx.OnChanged(cp.modeName)
	logx.Result("bind-cipher-mode: %s", cp.modeName)
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/
