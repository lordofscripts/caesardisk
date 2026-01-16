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

// service for gadgets that need to retrieve the current cipher mode
type ICipherModeService interface {
	Update()
	GetCipherMode() crypto.CaesarCipherMode
}

var _ ICipherModeService = (*CipherModeGadget)(nil)
var _ IGadget = (*CipherModeGadget)(nil)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type CipherModeGadget struct {
	parent IGadgetParent
	// Params/Options tab
	cipherLabel     *widget.Label
	cipherSelect    *widget.Select
	cipherContainer *fyne.Container
	card            *widget.Card
	previous        int
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

func newCipherModeGadget(parent IGadgetParent) *CipherModeGadget {
	return &CipherModeGadget{
		parent:   parent,
		previous: -1,
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// Instantiates all its widgets. But does not set initial values
// that would trigger OnChange.
func (g *CipherModeGadget) Define() IGadget {
	// · Options tab layout. See GetOptionContainer()
	g.cipherLabel = widget.NewLabel("Cipher mode")

	//g.cipherSelect = widget.NewSelect([]string{
	g.cipherSelect = widget.NewSelectWithData([]string{
		crypto.CaesarMode.String(),    // 0
		crypto.DidimusMode.String(),   // 1
		crypto.FibonacciMode.String(), // 2
		crypto.PrimusMode.String()},   // 3
		BoundCipherModeName)

	g.cipherContainer = container.New(layout.NewBorderLayout(nil, nil, g.cipherLabel, nil),
		g.cipherLabel,
		g.cipherSelect,
	)

	g.card = widget.NewCard(
		"Caesar Cipher Mode",
		"encoding algorithm variant",
		g.cipherContainer,
	)

	return g
}

// Lets the gadget bind data to the widgets
func (g *CipherModeGadget) Bind() IGadget {
	return g
}

// After the widgets and gadgets are defined and rendered in the
// application window, but prior to run, we call PostRender() to
// set widget values that may/will trigger onChange cascade events.
func (g *CipherModeGadget) PostRender() IGadget {
	g.cipherSelect.SetSelectedIndex(int(crypto.CaesarMode))
	g.cipherSelect.OnChanged = g.onChangeEnded

	return g
}

// When using a Card...
func (g *CipherModeGadget) Canvas() fyne.CanvasObject {
	return g.card
}

// Get the container that would be composited into another or a window
func (g *CipherModeGadget) Container() *fyne.Container {
	return g.cipherContainer
}

// Update should be called if there are unbound data model values
// that would require widget state to change.
func (g *CipherModeGadget) Update() {
	logx.OnUpdate()

	var sm crypto.SessionModel
	sm = DataBindings.GetSessionModel()

	g.cipherSelect.SetSelectedIndex(int(sm.Mode))
}

// Hide gadget
func (g *CipherModeGadget) Hide() {
	g.cipherContainer.Hide()
}

// Show gadget
func (g *CipherModeGadget) Show() {
	g.cipherContainer.Show()
}

// Enable gadget
func (g *CipherModeGadget) Enable() {
	g.cipherSelect.Enable()
}

// Disable gadget
func (g *CipherModeGadget) Disable() {
	g.cipherSelect.Disable()
}

// Clears all fields of a gadget
func (g *CipherModeGadget) Clear() {

}

// implement ICipherModeService
func (g *CipherModeGadget) GetCipherMode() crypto.CaesarCipherMode {
	var value crypto.CaesarCipherMode
	selected := g.cipherSelect.SelectedIndex()
	if selected > -1 {
		value = crypto.CaesarCipherMode(selected)
	} else {
		logx.Print("invalid selectedIndex at CipherModeGadget. using default")
		value = crypto.CaesarMode
	}

	return value
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

func (g *CipherModeGadget) onChangeEnded(s string) {
	logx.OnChanged()
	newMode := g.GetCipherMode()
	if g.previous != int(newMode) {
		g.previous = int(newMode)
		BoundCipherModeName.Set(s)
	}

	g.parent.Cascade(GadgetCipherMode, newMode)
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/
