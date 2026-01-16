/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/lordofscripts/caesardisk/crypto"
	"github.com/lordofscripts/goapp/app/logx"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ IGadget = (*SecretDataGadget)(nil)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type SecretDataGadget struct {
	textEntry1   *widget.Entry
	textEntry2   *widget.Entry
	buttonEncode *widget.Button
	buttonDecode *widget.Button
	container    *fyne.Container

	parentWindow fyne.Window
	engine       *crypto.CipherController
	alert        crypto.IViewNotifier
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				I n i t i a l i z e r
 *-----------------------------------------------------------------*/
func init() {}

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func newSecretDataGadget(window fyne.Window, alerter crypto.IViewNotifier) *SecretDataGadget {
	return &SecretDataGadget{
		parentWindow: window,
		engine:       nil,
		alert:        alerter,
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

func (g *SecretDataGadget) With(engine *crypto.CipherController) *SecretDataGadget {
	g.engine = engine

	return g
}

// Instantiates all its widgets. But does not set initial values
// that would trigger OnChange.
func (g *SecretDataGadget) Define() IGadget {
	// · Multi-line text widget for Input data (plain or ciphered)
	g.textEntry1 = widget.NewMultiLineEntry()
	g.textEntry1.SetPlaceHolder("Enter plain text here...")
	g.textEntry1.Wrapping = fyne.TextWrapBreak

	// · Multi-line text widget for output data (encrypted or plain)
	g.textEntry2 = widget.NewMultiLineEntry()
	g.textEntry2.SetPlaceHolder("Output will appear here...")
	g.textEntry2.Wrapping = fyne.TextWrapBreak

	// Action buttons
	// · Exchange : exchange text between Input & Output text entries
	//actionExchangeButton := widget.NewButton(string(rune(0x21d5)), g.onExchangeClicked)
	actionExchangeButton := widget.NewButtonWithIcon("XCHG", theme.SearchReplaceIcon(), g.onExchangeClicked)
	// · Clear : Left click clears both text entries
	actionClearButton := widget.NewButtonWithIcon("CLR", theme.DeleteIcon(), g.onClearClicked)
	// · Edit : Opens a larger modal window for better editing of Input text
	//actionEditButton := widget.NewButton(string(rune(0x1f440)), func() {
	actionEditButton := widget.NewButtonWithIcon("Edit", theme.DocumentCreateIcon(), func() {
		modal := NewModalEditor(g.parentWindow, "Comfortable Editor", g.textEntry1.Text, false)
		modal.OnAccept = func(s string) {
			g.textEntry1.SetText(s)
		}
		modal.Resize(fyne.NewSize(WINDOW_WIDTH, WINDOW_HEIGHT/2))
		modal.Show()
	})
	// · Encode : perform encryption in the selected mode and parameters
	g.buttonEncode = widget.NewButton("Encode", g.onEncodeClicked)
	g.buttonEncode.Disable()
	// · Decode : perform decryption in the selected mode and parameters
	g.buttonDecode = widget.NewButton("Decode", g.onDecodeClicked)
	g.buttonDecode.Disable()

	otherButtonContainer := container.NewHBox(layout.NewSpacer(), actionExchangeButton, actionClearButton, actionEditButton, layout.NewSpacer())
	actionButtonContainer := container.NewBorder(nil, nil, g.buttonEncode, g.buttonDecode,
		otherButtonContainer,
	)

	// OnChange - if there is text in the input entry, enable the buttons,
	// else disable the buttons
	g.textEntry1.OnChanged = g.onInputTextChanged

	g.container = container.NewVBox(
		g.textEntry1,
		actionButtonContainer,
		g.textEntry2,
	)

	return g
}

// Lets the gadget bind data to the widgets
func (g *SecretDataGadget) Bind() IGadget {
	return g
}

// After the widgets and gadgets are defined and rendered in the
// application window, but prior to run, we call PostRender() to
// set widget values that may/will trigger onChange cascade events.
func (g *SecretDataGadget) PostRender() IGadget {
	return g
}

// Get the container that would be composited into another or a window
func (g *SecretDataGadget) Container() *fyne.Container {
	return g.container
}

// Update should be called if there are unbound data model values
// that would require widget state to change.
func (g *SecretDataGadget) Update() {
	logx.OnUpdate()
}

// Hide gadget
func (g *SecretDataGadget) Hide() {
	g.container.Hide()
}

// Show gadget
func (g *SecretDataGadget) Show() {
	g.container.Show()
}

// Enable gadget
func (g *SecretDataGadget) Enable() {
	g.textEntry1.Enable()
	g.textEntry2.Enable()
	// the buttons only if there is input data
	g.updateActionButtonState()
}

// Disable gadget
func (g *SecretDataGadget) Disable() {
	g.textEntry1.Disable()
	g.textEntry2.Disable()
	// we can disregard CLR and XCH
	g.buttonDecode.Disable()
	g.buttonEncode.Disable()
}

// Clears all fields of a gadget
func (g *SecretDataGadget) Clear() {
	g.textEntry1.SetText("")
	g.textEntry2.SetText("")
	g.updateActionButtonState()
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

func (g *SecretDataGadget) updateActionButtonState() {
	if len(g.textEntry1.Text) > 0 {
		g.buttonEncode.Enable()
		g.buttonDecode.Enable()
	} else {
		g.buttonDecode.Disable()
		g.buttonEncode.Disable()
	}
}

// detect when the Input text entry changes so that we can
// update the action button states.
func (g *SecretDataGadget) onInputTextChanged(s string) {
	if len(s) > 0 {
		g.buttonEncode.Enable()
		g.buttonDecode.Enable()
	} else {
		g.buttonEncode.Disable()
		g.buttonDecode.Disable()
	}
}

// (Click) "Exchange" button
func (g *SecretDataGadget) onExchangeClicked() {
	logx.OnClick()

	temp := g.textEntry1.Text
	g.textEntry1.SetText(g.textEntry2.Text)
	g.textEntry2.SetText(temp)
}

// (Click) "Clear" button
func (g *SecretDataGadget) onClearClicked() {
	logx.OnClick()

	g.textEntry1.SetText("")
	g.textEntry2.SetText("")
}

// (Click) "Encode" button
func (g *SecretDataGadget) onEncodeClicked() {
	var err error = nil
	var result string

	logx.OnClick()

	g.textEntry2.SetText("") // clear output space
	var sm crypto.SessionModel
	sm = DataBindings.GetSessionModel()
	logx.Printf("Encode with %s", sm)

	cipherC := g.engine.CloneWith(&sm.Alpha)
	// · Encrypt operation
	if sm.Mode != crypto.DidimusMode {
		result, err = cipherC.Encrypt(sm.Mode, g.textEntry1.Text, sm.MainKey.Shift)
	} else {
		result, err = cipherC.Encrypt(sm.Mode, g.textEntry1.Text, sm.MainKey.Shift, sm.Offset)
	}

	if err != nil {
		logx.Printf("Encrypt Error: %v", err)
		g.textEntry2.SetText(err.Error())

		if g.alert != nil {
			g.alert.Notify(true, err.Error())
		}
	} else {
		// · PDU Packaging if requested by the user
		if usePDU, _ := BoundOptionUsePDU.Get(); usePDU {
			result = g.engine.PackMessage(result)
		}

		g.textEntry2.SetText(result)
	}
}

// (Click) "Decode" button
func (g *SecretDataGadget) onDecodeClicked() {
	var err error = nil
	var result string = ""

	logx.OnClick()

	g.textEntry2.SetText("") // clear output space
	var sm crypto.SessionModel
	sm = DataBindings.GetSessionModel()
	logx.Print(sm.String())
	cipherC := g.engine.CloneWith(&sm.Alpha)

	// · For PDUs we must unpack them first prior to Decrypting
	if usePDU, _ := BoundOptionUsePDU.Get(); usePDU {
		result, err = g.engine.UnpackMessage(g.textEntry1.Text, sm.Mode, sm.MainKey.Shift, sm.Offset)
		if err != nil {
			logx.AttentionAlways("PDU-Unpack", err)
			g.textEntry2.SetText(err.Error())

			if g.alert != nil {
				g.alert.Notify(true, err.Error())
			}

			return
		}
	} else {
		result = g.textEntry1.Text
	}

	// · Decrypt operation
	if sm.Mode != crypto.DidimusMode && sm.Mode != crypto.PrimusMode {
		result, err = cipherC.Decrypt(sm.Mode, result, sm.MainKey.Shift)
	} else {
		// Didimus & Primus use Offset value
		result, err = cipherC.Decrypt(sm.Mode, result, sm.MainKey.Shift, sm.Offset)
	}

	if err != nil {
		logx.Printf("Error: %v", err)
		g.textEntry2.SetText(err.Error())

		if g.alert != nil {
			g.alert.Notify(true, err.Error())
		}
	} else {
		g.textEntry2.SetText(result)
	}
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/
