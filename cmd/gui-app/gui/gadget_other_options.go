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
	"fyne.io/fyne/v2/widget"
	"github.com/lordofscripts/caesardisk"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ IGadget = (*MiscOptionsGadget)(nil)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type MiscOptionsGadget struct {
	parent IGadgetParent
	// Params/Options tab
	checkOrtho *widget.Check
	checkPDU   *widget.Check
	card       *widget.Card

	wheelOpts *caesardisk.CaesarWheelOptions
}

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func newMiscOptionsGadget(parent IGadgetParent, wheelOpts *caesardisk.CaesarWheelOptions) *MiscOptionsGadget {
	return &MiscOptionsGadget{
		parent:    parent,
		wheelOpts: wheelOpts,
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// Instantiates all its widgets. But does not set initial values
// that would trigger OnChange.
func (g *MiscOptionsGadget) Define() IGadget {
	g.checkOrtho = widget.NewCheck("Orthogonal", nil)
	g.checkOrtho.SetChecked(g.wheelOpts.Orthogonal) // we don't want to trigger just yet
	g.checkOrtho.OnChanged = func(b bool) {
		g.wheelOpts.Orthogonal = g.checkOrtho.Checked
	}

	g.checkPDU = widget.NewCheckWithData("Use PDU format", BoundOptionUsePDU)

	miscCardContent := container.NewVBox(
		g.checkOrtho,
		g.checkPDU,
	)

	g.card = widget.NewCard(
		"Misc. Options",
		"Rendering options",
		miscCardContent,
	)

	return g
}

// Lets the gadget bind data to the widgets
func (g *MiscOptionsGadget) Bind() IGadget {
	return g
}

// After the widgets and gadgets are defined and rendered in the
// application window, but prior to run, we call PostRender() to
// set widget values that may/will trigger onChange cascade events.
func (g *MiscOptionsGadget) PostRender() IGadget {

	return g
}

// When using a Card...
func (g *MiscOptionsGadget) Canvas() fyne.CanvasObject {
	return g.card
}

// Get the container that would be composited into another or a window
func (g *MiscOptionsGadget) Container() *fyne.Container {
	return nil
}

// Update should be called if there are unbound data model values
// that would require widget state to change.
func (g *MiscOptionsGadget) Update() {
}

// Hide gadget
func (g *MiscOptionsGadget) Hide() {
	g.card.Hide()
}

// Show gadget
func (g *MiscOptionsGadget) Show() {
	g.card.Show()
}

// Enable gadget
func (g *MiscOptionsGadget) Enable() {
	g.checkOrtho.Enable()
	g.checkPDU.Enable()
}

// Disable gadget
func (g *MiscOptionsGadget) Disable() {
	g.checkOrtho.Disable()
	g.checkPDU.Disable()
}

// Clears all fields of a gadget
func (g *MiscOptionsGadget) Clear() {
	g.checkOrtho.SetChecked(false)
	g.checkPDU.SetChecked(false)
}

func (g *MiscOptionsGadget) GetRenderOrthogonality() bool {
	return g.checkOrtho.Checked
}

func (g *MiscOptionsGadget) GetOperUsePDU() bool {
	return g.checkPDU.Checked
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/
