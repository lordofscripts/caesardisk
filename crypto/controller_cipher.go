/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A public controller for Caesar-class substitution cipher encryption
 * and decryption. It supports plain Caesar, Didimus, Fibonacci & Primus.
 *-----------------------------------------------------------------*/
package crypto

import (
	"errors"

	"github.com/lordofscripts/caesardisk"
	"github.com/lordofscripts/caesardisk/internal/cipher"
	"github.com/lordofscripts/caesardisk/internal/hash"
	"github.com/lordofscripts/goapp/app/logx"
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

// The CipherController object holds a reference alphabet and can
// perform repeated, independent Encryption/Decryption operations
// on that alphabet with different parameter values.
type CipherController struct {
	ControllerBase
	alpha *caesardisk.AlphabetModel
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				I n i t i a l i z e r
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewCipherController(alpha *caesardisk.AlphabetModel, vwn IViewNotifier) *CipherController {
	if alpha == nil {
		logx.Ctor()
		logx.Fatalln("nil alphabet to CipherController")
	}

	return &CipherController{
		ControllerBase: ControllerBase{
			viewNotify: vwn,
		},
		alpha: alpha,
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

func (cc *CipherController) CloneWith(newAlpha *caesardisk.AlphabetModel) *CipherController {
	alter := &CipherController{
		ControllerBase: ControllerBase{
			viewNotify: cc.viewNotify,
		},
		alpha: cc.alpha,
	}

	if newAlpha != nil {
		alter.alpha = newAlpha
	}
	return alter
}

// Encrypt a plain string using the selected Caesar-class cipher mode
// with the selected encryption parameters.
func (cc *CipherController) Encrypt(mode CaesarCipherMode, plain string, keyShift int, args ...any) (string, error) {
	var sequencer cipher.IKeySequencer
	var p *cipher.CaesarParameters = cipher.NewCaesarParameters(cc.alpha)
	logx.Enter()
	defer logx.Leave()

	cc.mutexRW.Lock()
	defer cc.mutexRW.Unlock()

	// I. Select parameters and create sequencer
	switch mode {
	case CaesarMode:
		p.SetKey(keyShift)
		sequencer = cipher.NewCaesarSequencer(p)

	case DidimusMode:
		if len(args) != 1 {
			return "", errors.New("missing Offset parameter to Didimus encryptor")
		}
		if ofs, ok := (args[0]).(int); ok {
			p.SetKey(keyShift)
			p.SetAltKeyOffset(ofs)
			sequencer = cipher.NewDidimusSequencer(p)
		} else {
			return "", errors.New("invalid parameter type to Didimus")
		}

	case FibonacciMode:
		p.SetKey(keyShift)
		sequencer = cipher.NewFibonacciSequencer(p)

	case PrimusMode:
		p.SetKey(keyShift)
		sequencer = cipher.NewPrimusSequencer(p)

	default:
		return "", errors.New("invalid cipher mode given to controller")
	}

	// II. Validate parameters via sequencer
	if err := sequencer.Validate(); err != nil {
		return "", err
	}

	// III. Encrypt
	caesarHandler := cipher.NewCaesarCipherFromSequencer(sequencer)

	return caesarHandler.Encode(plain), nil
}

// Decrypts a Caesar-class string using the selected cipher mode and decryption
// parameters.
func (cc *CipherController) Decrypt(mode CaesarCipherMode, ciphered string, keyShift int, args ...any) (string, error) {
	var sequencer cipher.IKeySequencer
	var p *cipher.CaesarParameters = cipher.NewCaesarParameters(cc.alpha)

	logx.Enter()
	defer logx.Leave()

	cc.mutexRW.Lock()
	defer cc.mutexRW.Unlock()

	// I. Select parameters and create sequencer
	switch mode {
	case CaesarMode:
		p.SetKey(keyShift)
		sequencer = cipher.NewCaesarSequencer(p)

	case DidimusMode:
		if len(args) != 1 {
			return "", errors.New("missing Offset parameter to Didimus decryptor")
		}
		if ofs, ok := (args[0]).(int); ok {
			p.SetKey(keyShift)
			p.SetAltKeyOffset(ofs)
			sequencer = cipher.NewDidimusSequencer(p)
		} else {
			return "", errors.New("invalid parameter type to Didimus")
		}

	case FibonacciMode:
		p.SetKey(keyShift)
		sequencer = cipher.NewFibonacciSequencer(p)

	case PrimusMode:
		p.SetKey(keyShift)
		sequencer = cipher.NewPrimusSequencer(p)

	default:
		return "", errors.New("invalid cipher mode given to controller")
	}

	// II. Validate parameters via sequencer
	if err := sequencer.Validate(); err != nil {
		return "", err
	}
	logx.Printf("Decrypt %s sequencer validated", sequencer)
	// III. Decrypt
	caesarHandler := cipher.NewCaesarCipherFromSequencer(sequencer)

	return caesarHandler.Decode(ciphered), nil
}

// takes an already encrypted string and packages it in a PDU that can
// be sent over the communications channel.
func (cc *CipherController) PackMessage(cipherPayload string) string {
	logx.Enter()
	defer logx.Leave()

	msgPDU := cipher.NewCaesarMessage(hash.NewXXH64(hashSeed))
	msgPDU.AddMessage(cipherPayload)

	return msgPDU.String()
}

// takes a PDU that contains an encrypted message from a communications
// channel and unpacks it to verify the metadata and if successful,
// return the decrypted payload.
func (cc *CipherController) UnpackMessage(pdu string, mode CaesarCipherMode, keyShift int, keyOffset int) (string, error) {
	logx.Enter()
	defer logx.Leave()

	check := hash.NewXXH64(hashSeed)
	payload, err := cipher.VerifyCaesarMessage(check, pdu)
	if err == nil {
		var plain string
		if mode != DidimusMode {
			plain, err = cc.Decrypt(mode, payload, keyShift)
		} else {
			plain, err = cc.Decrypt(mode, payload, keyShift, keyOffset)
		}

		return plain, err
	}

	return "", err
}

func (cc *CipherController) CaesarCorrection(keyShift int) (main CaesarKey, warn error) {
	var shf int
	shf, warn = cipher.CaesarCorrection(keyShift, cc.alpha)
	let, _ := cc.alpha.Character(shf)

	return CaesarKey{
		Shift:  shf,
		Letter: let,
	}, warn
}

// The key schedule for a plain Caesar scheduler
func (cc *CipherController) GetCaesarSchedule(keyShift int) (KeySchedule, error) {
	// parameters for the sequencer
	p := cipher.NewCaesarParameters(cc.alpha)
	p.KeyValue = keyShift
	// the sequencer that will provide us the raw sequence of keys
	seq := cipher.NewCaesarSequencer(p)
	if seq.Validate() != nil {
		// obtain corrected parameters.
		p = seq.GetParams()
	}
	// get the raw key schedule
	rawSchedule := seq.GetRawKeySchedule()
	// convert it to a public API object
	qty := len(rawSchedule)
	schedule := make(KeySchedule, qty)
	for i, raw := range rawSchedule {
		shf := raw.KeyShift
		chr, _ := cc.alpha.Character(shf)
		schedule[i] = KeyScheduleItem{
			KeyShift: shf,
			KeyChar:  chr,
			Comment:  raw.Comment,
			Tabula:   cipher.RotateStringLeft(cc.alpha.String(), shf),
		}
	}
	return schedule, nil
}

// The key schedule for a Didimus scheduler
func (cc *CipherController) GetDidimusSchedule(keyShift, keyOffset int) (KeySchedule, error) {
	// parameters for the sequencer
	p := cipher.NewCaesarParameters(cc.alpha)
	p.KeyValue = keyShift
	p.Offset = keyOffset
	// the sequencer that will provide us the raw sequence of keys
	seq := cipher.NewDidimusSequencer(p)
	if seq.Validate() != nil {
		// obtain corrected parameters.
		p = seq.GetParams()
	}
	// get the raw key schedule
	rawSchedule := seq.GetRawKeySchedule()
	// convert it to a public API object
	qty := len(rawSchedule)
	schedule := make(KeySchedule, qty)
	for i, raw := range rawSchedule {
		shf := raw.KeyShift
		chr, _ := cc.alpha.Character(shf)
		schedule[i] = KeyScheduleItem{
			KeyShift: shf,
			KeyChar:  chr,
			Comment:  raw.Comment,
			Tabula:   cipher.RotateStringLeft(cc.alpha.String(), shf),
		}
	}
	return schedule, nil
}

// The key schedule for a Fibonacci scheduler
func (cc *CipherController) GetFibonacciSchedule(keyShift int) (KeySchedule, error) {
	// parameters for the sequencer
	p := cipher.NewCaesarParameters(cc.alpha)
	p.KeyValue = keyShift
	// the sequencer that will provide us the raw sequence of keys
	seq := cipher.NewFibonacciSequencer(p)
	if seq.Validate() != nil {
		// obtain corrected parameters.
		p = seq.GetParams()
	}
	// get the raw key schedule
	rawSchedule := seq.GetRawKeySchedule()
	// convert it to a public API object
	qty := len(rawSchedule)
	schedule := make(KeySchedule, qty)
	for i, raw := range rawSchedule {
		shf := raw.KeyShift
		chr, _ := cc.alpha.Character(shf)
		schedule[i] = KeyScheduleItem{
			KeyShift: shf,
			KeyChar:  chr,
			Comment:  raw.Comment,
			Tabula:   cipher.RotateStringLeft(cc.alpha.String(), shf),
		}
	}
	return schedule, nil
}

// The key schedule for a Primus scheduler
func (cc *CipherController) GetPrimusSchedule(keyShift, keyOffset int) (KeySchedule, error) {
	// parameters for the sequencer
	p := cipher.NewCaesarParameters(cc.alpha)
	p.KeyValue = keyShift
	p.Offset = keyOffset
	// the sequencer that will provide us the raw sequence of keys
	seq := cipher.NewPrimusSequencer(p)
	if seq.Validate() != nil {
		// obtain corrected parameters.
		p = seq.GetParams()
	}
	// get the raw key schedule
	rawSchedule := seq.GetRawKeySchedule()
	// convert it to a public API object
	qty := len(rawSchedule)
	schedule := make(KeySchedule, qty)
	for i, raw := range rawSchedule {
		shf := raw.KeyShift
		chr, _ := cc.alpha.Character(shf)
		schedule[i] = KeyScheduleItem{
			KeyShift: shf,
			KeyChar:  chr,
			Comment:  raw.Comment,
			Tabula:   cipher.RotateStringLeft(cc.alpha.String(), shf),
		}
	}
	return schedule, nil
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/
