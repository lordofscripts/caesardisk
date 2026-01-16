/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package gui

import "fyne.io/fyne/v2"

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	GadgetAlphabet GadgetId = iota
	GadgetCipherMode
	GadgetMainKey
	GadgetOffset
	GadgetWheel      // receive only
	GadgetSecretData // receive only
	GadgetOtherOpts
)

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

type IGadget interface {
	// Instantiates all its widgets. But does not set initial values
	// that would trigger OnChange.
	Define() IGadget
	// Lets the gadget bind data to the widgets
	Bind() IGadget
	// After the widgets and gadgets are defined and rendered in the
	// application window, but prior to run, we call PostRender() to
	// set widget values that may/will trigger onChange cascade events.
	PostRender() IGadget
	// Get the container that would be composited into another or a window
	Container() *fyne.Container
	// Update should be called if there are unbound data model values
	// that would require widget state to change.
	Update()
	// Hide gadget
	Hide()
	// Show gadget
	Show()
	// Enable gadget
	Enable()
	// Disable gadget
	Disable()
	// Clears all fields of a gadget
	Clear()
}

// ask the parent of all gadgets to relay a value change to
// a gadget's sibling.
type IGadgetParent interface {
	Cascade(serviceId GadgetId, value any)
}

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

// unique identifier type of objects implementing IGadget. The actual
// enumeration set is defined in the user application space.
type GadgetId uint8

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

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/
