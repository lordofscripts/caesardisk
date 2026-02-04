/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *						goCaesarDisk GUI
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A modal window where the key schedule is displayed as per the
 * selected cipher mode. It is however cipher-agnostic becaues it
 * displays them as a table.
 *-----------------------------------------------------------------*/
package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/lordofscripts/caesardisk/crypto"
	"github.com/lordofscripts/gofynex/fynex"
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
type KeyScheduleViewer struct {
	isDismissed bool
	modal       *dialog.CustomDialog
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// Create a modal editor but don't show it yet.
func NewKeyScheduleViewer(w fyne.Window, title string, schedule crypto.KeySchedule) *KeyScheduleViewer {
	me := &KeyScheduleViewer{
		isDismissed: false,
		modal:       nil,
	}
	me.build(w, title, schedule)

	return me
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

func (ksv *KeyScheduleViewer) Resize(size fyne.Size) {
	ksv.modal.Resize(size)
}

// Show the modal editor window if it has not been dismissed.
func (ksv *KeyScheduleViewer) Show() {
	if !ksv.isDismissed {
		ksv.modal.Show()
	}
}

// Hide the modal editor window if it has not been dismissed.
func (ksv *KeyScheduleViewer) Hide() {
	if !ksv.isDismissed {
		ksv.modal.Hide()
	}
}

// Destroy the modal editor window. The instance can no longer be shown.
// It is ignored if it has already been dismissed.
func (ksv *KeyScheduleViewer) Dismiss() {
	if !ksv.isDismissed {
		ksv.modal.Dismiss()
	}
}

// Check whether the modal editor has already been dismissed/destroyed.
// If true, don't call Show()
func (ksv *KeyScheduleViewer) IsDismissed() bool {
	return ksv.isDismissed
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

// build up the modal window containing the multi-line text box and
// the Accept & Cancel buttons at the bottom.
func (ksv *KeyScheduleViewer) build(w fyne.Window, title string, schedule crypto.KeySchedule) {
	dataWidget := widget.NewTable(
		// Table: Dimensions
		func() (rows int, cols int) {
			return len(schedule), 4
		},
		// Table: Create Cell
		func() fyne.CanvasObject {
			dlabl := fynex.NewDynamicLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{
				Monospace: true,
			}, nil)
			return dlabl
		},
		// Table: Populate
		func(id widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*fynex.DynamicLabel)
			item := schedule[id.Row]
			switch id.Col {
			case 0:
				label.SetText(string(item.KeyChar))
			case 1:
				label.SetText(fmt.Sprintf("%02d", item.KeyShift))
			case 2:
				label.SetText(string(item.Tabula))
			case 3:
				label.SetText(item.Comment)
			}
		})

	// 1. MUST set explicit column widths
	dataWidget.SetColumnWidth(0, 25)  // KeyChar
	dataWidget.SetColumnWidth(1, 30)  // KeyShift
	dataWidget.SetColumnWidth(2, 250) // Tabula
	dataWidget.SetColumnWidth(3, 100) // Comment

	// 2. Wrap in a Max container so the table fills the dialog
	editorContainer := container.NewStack(dataWidget)
	ksv.modal = dialog.NewCustom(title+" Key Schedule", "OK", editorContainer, w)

	// 3. MUST resize the dialog, otherwise it collapses to minimum size
	ksv.modal.Resize(fyne.NewSize(600, 400))
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/
