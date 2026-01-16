/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package gui

import (
	"bytes"
	"embed"
	"image"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"github.com/lordofscripts/caesardisk"
	"github.com/lordofscripts/caesardisk/crypto"
	"github.com/lordofscripts/goapp/app/logx"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	wheel_WIDTH          = 300
	wheel_HEIGHT         = 300
	wheelCaption_MAINKEY = "Main key"
	wheelCaption_ALTKEY  = "Alternate key"
)

//go:embed caesar_disk_sample.png
var caesarDiskSample embed.FS

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ crypto.IWheelUpdateService = (*WheelGadget)(nil)
var _ IGadget = (*WheelGadget)(nil)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type WheelGadget struct {
	image *canvas.Image
	//imageW    *ClickImage
	sample    image.Image
	caption   *DynamicLabel //*widget.Label
	labelMode *DynamicLabel
	container *fyne.Container

	mutex sync.Mutex

	last        lastParameters
	wheelOpts   *caesardisk.CaesarWheelOptions
	showMainKey bool
	isBipolar   bool
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

type lastParameters struct {
	Alpha       string
	KeyShift    int
	OffsetShift int
}

/* ----------------------------------------------------------------
 *				I n i t i a l i z e r
 *-----------------------------------------------------------------*/
func init() {}

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func newWheelGadget(wopts *caesardisk.CaesarWheelOptions) *WheelGadget {
	return &WheelGadget{
		showMainKey: true,
		isBipolar:   false, // by default non-Didimus i.e. no Offset shown
		last:        lastParameters{},
		wheelOpts:   wopts,
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// Instantiates all its widgets. But does not set initial values
// that would trigger OnChange.
func (g *WheelGadget) Define() IGadget {
	displaySize := fyne.NewSize(wheel_WIDTH, wheel_HEIGHT)

	if sample, err := loadEmbeddedImage("caesar_disk_sample.png"); err != nil {
		logx.Fatal(err.Error()) // @audit this doesn't need to be fatal

		return g
	} else {
		//smallTextLabel := widget.NewLabelWithStyle(wheelCaption_MAINKEY, fyne.TextAlignCenter, fyne.TextStyle{
		smallTextLabel := NewDynamicLabelWithStyle(wheelCaption_MAINKEY, fyne.TextAlignCenter, fyne.TextStyle{
			Bold:      true,
			Monospace: true,
		}, nil)
		smallTextLabel.SizeName = theme.SizeNameCaptionText
		smallTextLabel.Wrapping = fyne.TextWrapOff
		g.caption = smallTextLabel

		// · Cipher Mode label (informational)
		g.labelMode = NewDynamicLabelWithStyle(InitialCipherMode.String(), fyne.TextAlignCenter, fyne.TextStyle{
			Bold:      true,
			Monospace: true,
		}, nil)
		g.labelMode.SizeName = theme.SizeNameCaptionText
		g.labelMode.Wrapping = fyne.TextWrapOff
		g.labelMode.Bind(BoundCipherModeName)

		img := canvas.NewImageFromImage(sample)
		img.SetMinSize(displaySize)
		img.Resize(displaySize)
		img.FillMode = canvas.ImageFillContain

		g.image = img
		//g.imageW = NewClickImageFromImage(img)
		g.sample = sample

		labelContainer := container.NewBorder(nil, nil, g.labelMode, g.caption)

		g.container = container.NewVBox(img, labelContainer)
	}
	return g
}

// Lets the gadget bind data to the widgets
func (g *WheelGadget) Bind() IGadget {
	return g
}

// After the widgets and gadgets are defined and rendered in the
// application window, but prior to run, we call PostRender() to
// set widget values that may/will trigger onChange cascade events.
func (g *WheelGadget) PostRender() IGadget {
	g.caption.OnTapped = g.onTapped
	//g.imageW.OnTapped = g.onTapped
	return g
}

// Get the container that would be composited into another or a window
func (g *WheelGadget) Container() *fyne.Container {
	return g.container
}

// Update should be called if there are unbound data model values
// that would require widget state to change.
func (g *WheelGadget) Update() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	logx.OnUpdate()

	g.updateVisual()
}

// Hide gadget
func (g *WheelGadget) Hide() {
	g.container.Hide()
}

// Show gadget
func (g *WheelGadget) Show() {
	g.container.Show()
}

// Enable gadget
func (g *WheelGadget) Enable() {
}

// Disable gadget
func (g *WheelGadget) Disable() {
}

// Clears all fields of a gadget
func (g *WheelGadget) Clear() {
	displaySize := fyne.NewSize(wheel_WIDTH, wheel_HEIGHT)
	img := canvas.NewImageFromImage(g.sample)
	img.SetMinSize(displaySize)
	img.Resize(displaySize)
	img.FillMode = canvas.ImageFillContain

	g.container.Objects[0] = img
}

func (g *WheelGadget) CanShowOffset(can bool) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.isBipolar = can
	if !can {
		g.showMainKey = true // only main key can be shown
		g.updateVisual()
	}
}

// Only meaninful for bi-alphabetic ciphers. It indicates
// whether the wheel is currently showing the MAIN key (true)
// rather than the alternate key (false). Useful to avoid
// unnecesary wheel updates.
func (g *WheelGadget) IsShowingMain() bool {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	return g.showMainKey
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

// func (g *WheelGadget) onTapped(evt *fyne.PointEvent) {
func (g *WheelGadget) onTapped() {
	logx.OnClick()

	if g.isBipolar {
		logx.EventEnter()
		// @todo switch between Main & Alternate disk
		g.showMainKey = !g.showMainKey
		g.updateVisual()
		logx.EventLeave()
	}
}

func (g *WheelGadget) updateVisual() {
	var sm crypto.SessionModel
	sm = DataBindings.GetSessionModel()

	// · retrieve alphabet name & length
	reqAlphaLen := sm.Alpha.Length()
	reqAlphaChars := sm.Alpha.String()
	// · retrieve main Key Shift value or Alternate Shift (Didimus)
	var reqKeyShift int
	if g.showMainKey || !g.isBipolar {
		reqKeyShift = sm.MainKey.Shift
		g.caption.SetText(wheelCaption_MAINKEY)
	} else {
		reqKeyShift = sm.AltKey.Shift
		g.caption.SetText(wheelCaption_ALTKEY)
	}

	// · compose requested parameters
	reqValues := lastParameters{
		Alpha:       reqAlphaChars,
		KeyShift:    reqKeyShift,
		OffsetShift: sm.Offset,
	}
	// · only do work if there has been a change
	if g.last.Equal(reqValues) {
		return
	}

	logx.Printf("wheelUpdate α:%s k:%02d O:%02d", sm.Alpha.Name, reqKeyShift, sm.Offset)

	const GENERATE_DUAL_ALPHABET_DISK bool = false
	var imgBase, imgOverlay, imgComposite image.Image
	var err error = nil

	// base/outer
	if imgBase, err = caesardisk.GenerateCaesarWheelImage(
		reqAlphaChars, false, *g.wheelOpts); err == nil {
		// overlay/inner
		if imgOverlay, err = caesardisk.GenerateCaesarWheelImage(
			reqAlphaChars, true, *g.wheelOpts); err == nil {
			if imgComposite, err = caesardisk.SuperimposeDisksByShiftImage(
				reqKeyShift,
				reqAlphaLen,
				imgBase,
				imgOverlay,
				GENERATE_DUAL_ALPHABET_DISK,
				*g.wheelOpts,
			); err == nil {
				g.image.Image = imgComposite
				g.image.Refresh()
				//g.imageW.SetImage(canvas.NewImageFromImage(imgComposite))
				//g.imageW.Refresh()
			}
		}
	}

	if err != nil {
		logx.AttentionAlways("updateWheel", err)
	} else {
		g.last = reqValues
	}
}

func (l lastParameters) Equal(other lastParameters) bool {
	return strings.EqualFold(l.Alpha, other.Alpha) &&
		l.KeyShift == other.KeyShift &&
		l.OffsetShift == other.OffsetShift
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

func loadEmbeddedImage(path string) (image.Image, error) {
	data, err := caesarDiskSample.ReadFile(path)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return img, nil
}
