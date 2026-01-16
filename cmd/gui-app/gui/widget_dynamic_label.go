/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   photoQ
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A custom Fyne widget that extends widget.Label{} to deal with Fyne
 * quirkness. This dynamic label sports a Validation callback. It
 * serves a purpose when you need to display read-only data that looks
 * good, contrary to Fyne's disabled widget.Entry{} which is difficult
 * to view in that state. Fyne chose not implement Read-Only functionality
 * even though it is a very valid use case.
 *-----------------------------------------------------------------*/
package gui

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

/**
 *   In all GUI frameworks you can change the state of a disabled widget. In Fyne you cannot.
 * In Fyne is is difficult, if not impossible, to alter the visual colors of a disabled widget.
 * In Fyne there is no such thing as a Read-Only widget, quite absurd, thus making us go through
 * complicated workarounds.
 *  	I needed to display read-only checkboxes or indicators that changed according to data,
 * but that the user could not interact with, just look at them to get to know whether it is
 * set or not.
 *	  This widget is like that except it displays Green for checked and Red for unchecked in
 * the form of a filled circle.
 * Usage:
 *		stateIndicator := NewStateIndicator().SetChecked(true)
 *		stateIndicator.SetState(gui.IndicatorUnset)
 */
type DynamicLabel struct {
	widget.Label
	OnChanged func(text string)
	OnTapped  func()
	locker    sync.Mutex
}

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewDynamicLabel(text string, onChanged func(string)) *DynamicLabel {
	lbl := &DynamicLabel{}
	lbl.OnChanged = onChanged
	lbl.OnTapped = nil
	lbl.Selectable = false
	lbl.Label.SetText(text)
	lbl.locker = sync.Mutex{}
	lbl.ExtendBaseWidget(lbl)
	return lbl
}

func NewDynamicLabelWithStyle(text string, alignment fyne.TextAlign, style fyne.TextStyle, onChanged func(string)) *DynamicLabel {
	lbl := &DynamicLabel{}
	lbl.OnChanged = onChanged
	lbl.Selectable = false
	lbl.Label.SetText(text)
	lbl.Label.TextStyle = style
	lbl.Label.Alignment = alignment
	lbl.locker = sync.Mutex{}
	lbl.ExtendBaseWidget(lbl)
	return lbl
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

/**
 * Sets the label text and triggers the OnChange callback (if any)
 */
func (d *DynamicLabel) SetText(text string) {
	d.locker.Lock()
	defer d.locker.Unlock()

	d.Label.SetText(text)
	if d.OnChanged != nil {
		d.OnChanged(text)
	}
}

// implements Tappable interface
func (d *DynamicLabel) Tapped(evt *fyne.PointEvent) {
	if d.OnTapped != nil {
		d.OnTapped()
	}
}
