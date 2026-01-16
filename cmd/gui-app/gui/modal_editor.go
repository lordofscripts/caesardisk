/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *						goCaesarDisk GUI
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A modal window where the user can edit text in a larger multi-line
 * entry widget and pass back changes when done.
 *-----------------------------------------------------------------*/
package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
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

// ModalEditor creates a manages a custom modal window for editing text.
type ModalEditor struct {
	isReadOnly  bool
	isDismissed bool
	modal       *dialog.CustomDialog
	// Callback when Accept button is pressed on an editable modal
	OnAccept func(string)
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// Create a modal editor but don't show it yet.
func NewModalEditor(w fyne.Window, title, text string, readOnly bool) *ModalEditor {
	me := &ModalEditor{
		isReadOnly:  readOnly,
		isDismissed: false,
	}
	me.build(w, title, text)
	return me
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

func (me *ModalEditor) Resize(size fyne.Size) {
	me.modal.Resize(size)
}

// Show the modal editor window if it has not been dismissed.
func (me *ModalEditor) Show() {
	if !me.isDismissed {
		me.modal.Show()
	}
}

// Hide the modal editor window if it has not been dismissed.
func (me *ModalEditor) Hide() {
	if !me.isDismissed {
		me.modal.Hide()
	}
}

// Destroy the modal editor window. The instance can no longer be shown.
// It is ignored if it has already been dismissed.
func (me *ModalEditor) Dismiss() {
	if !me.isDismissed {
		me.modal.Dismiss()
	}
}

// Check whether the modal editor has already been dismissed/destroyed.
// If true, don't call Show()
func (me *ModalEditor) IsDismissed() bool {
	return me.isDismissed
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

// build up the modal window containing the multi-line text box and
// the Accept & Cancel buttons at the bottom.
func (me *ModalEditor) build(w fyne.Window, title, initialText string) {
	var dataWidget fyne.CanvasObject
	var rtb *widget.RichText = nil
	var mle *widget.Entry = nil

	if me.isReadOnly {
		rtb = widget.NewRichTextWithText(initialText)
		dataWidget = rtb
	} else {
		mle = widget.NewMultiLineEntry()
		mle.MultiLine = true
		mle.PlaceHolder = "Edit your text here..."
		mle.Text = initialText
		mle.Wrapping = fyne.TextWrapOff
		dataWidget = mle
	}

	// · The "Accept" button notifies the parent window of the
	//   new text value provided OnAccept has been set. This only
	//   happens if the modal is not read-only. The modal is dismissed.
	var modal *dialog.CustomDialog
	buttonAccept := widget.NewButton("Accept", func() {
		if !me.isReadOnly {
			println("Accepted", mle.Text)
			// pass back the new text
			if me.OnAccept != nil {
				me.OnAccept(mle.Text)
			}
		} else {
			println("Accepted", rtb.String())
		}
		me.isDismissed = true
		modal.Dismiss()
	})
	// · The "Cancel" button indicates the user wants to discard any
	//	 text modification. The modal window is dismissed.
	buttonCancel := widget.NewButton("Cancel", func() {
		me.isDismissed = true
		modal.Dismiss()
	})

	editorContainer := container.NewBorder(nil, nil, nil, nil, dataWidget)

	modal = dialog.NewCustom(title, "", editorContainer, w)
	modal.SetButtons([]fyne.CanvasObject{buttonAccept, buttonCancel})
	me.modal = modal
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

/*
func demo() {
	modal := NewModalEditor(g.w, "My title", "initial text", false)
	modal.OnAccept = func(s string) {
		println("Accepted", s)
	}
	modal.Resize(fyne.NewSize(WINDOW_WIDTH, WINDOW_HEIGHT))
	modal.Show()
}
*/
