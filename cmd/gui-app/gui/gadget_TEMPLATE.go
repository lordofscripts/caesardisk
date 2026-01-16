//go:build exclude

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
	"github.com/lordofscripts/caesardisk/crypto"
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

type NAMEGadget struct {
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

func newNAMEGadget() *NAMEGadget {
	return &NAMEGadget{}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

func (g *NAMEGadget) With(session *crypto.SessionController) *NAMEGadget {

}

// Instantiates all its widgets. But does not set initial values
// that would trigger OnChange.
func (g *NAMEGadget) Define() IGadget {

}

// Lets the gadget bind data to the widgets
func (g *NAMEGadget) Bind() IGadget {

}

// After the widgets and gadgets are defined and rendered in the
// application window, but prior to run, we call PostRender() to
// set widget values that may/will trigger onChange cascade events.
func (g *NAMEGadget) PostRender() IGadget {

}

// Get the container that would be composited into another or a window
func (g *NAMEGadget) Container() *fyne.Container {

}

// Update should be called if there are unbound data model values
// that would require widget state to change.
func (g *NAMEGadget) Update() {

}

// Hide gadget
func (g *NAMEGadget) Hide() {

}

// Show gadget
func (g *NAMEGadget) Show() {

}

// Enable gadget
func (g *NAMEGadget) Enable() {

}

// Disable gadget
func (g *NAMEGadget) Disable() {

}

// Clears all fields of a gadget
func (g *NAMEGadget) Clear() {

}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/
