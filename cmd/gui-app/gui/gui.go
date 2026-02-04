/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   APP_NAME
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package gui

import (
	_ "embed"
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"github.com/lordofscripts/caesardisk"
	"github.com/lordofscripts/caesardisk/crypto"
	"github.com/lordofscripts/caesardisk/internal/cipher"
	"github.com/lordofscripts/goapp/app/logx"
	"github.com/lordofscripts/gofynex/fynex"
	"github.com/lordofscripts/gofynex/fynex/dlg"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	APP_TITLE     = "Modern Caesar Cipher"
	APP_ID        = "com.lordofscripts.caesardisk_gui"
	WINDOW_WIDTH  = 380
	WINDOW_HEIGHT = 600

	// the default initial alphabet of the application (encode/decode)
	InitialAlphabetName = "English"
	// the initial Caesar mode
	InitialCipherMode = crypto.CaesarMode
)

// content for the About dialog
const aboutCONTENT string = `
## CaesarDisk GUI
A graphical user interface that integrates the
*CaesarDisk* disk generator application. As a
plus, you can also use various modes of the 
Caesar cipher:
* Caesar
* Didimus
* Fibonacci
* Primus

**CaesarDisk** and **CaesarDiskGUI** are
Copyright ©2025 LordOfScripts™
`

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ crypto.IViewNotifier = (*MainGUI)(nil)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

// the list of all alphabets
type AlphabetList map[string]*caesardisk.AlphabetModel

type MenuItemActionCB func()

type MainGUI struct {
	a        fyne.App
	w        fyne.Window
	dlgAbout *dlg.AboutBox

	gadgets     *appGadgets
	controllers *appControllers

	allAlphabets *AlphabetList
	wheelOpts    *caesardisk.CaesarWheelOptions
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

// application composite controls @todo convert to composite/custom widgets
type appGadgets struct {
	//     Type					Provides			Cascade
	Alpha  *AlphabetGadget   // IAlphabetService	Key,Offset,Wheel
	Cipher *CipherModeGadget // ICipherModeService  Wheel,Offset
	Wheel  *WheelGadget      //	IWheelUpdateService
	Key    *CaesarKeyGadget  // ICaesarKeyService	Wheel
	Offset *KeyOffsetGadget  // IKeyOffsetService	Wheel
	Data   *SecretDataGadget //
	Misc   *MiscOptionsGadget
}

// application controllers
type appControllers struct {
	Cipher *crypto.CipherController
}

/* ----------------------------------------------------------------
 *				I n i t i a l i z e r
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewGUI(alist *AlphabetList, wopts *caesardisk.CaesarWheelOptions) *MainGUI {
	// Data bindings that will be used by some widgets
	DataBindings.Init() // set defaults

	return &MainGUI{
		allAlphabets: alist,
		wheelOpts:    wopts,
		gadgets:      &appGadgets{},
		controllers:  &appControllers{},
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

func (g *MainGUI) Define() *MainGUI {
	myApp := app.NewWithID(APP_ID)
	myApp.SetIcon(applicationIcon)
	myApp.Lifecycle().SetOnStarted(g.OnStarted)
	myWindow := myApp.NewWindow(APP_TITLE)
	myWindow.SetIcon(applicationIcon)

	// Fix Metadata as it is usually the case
	meta := myApp.Metadata()
	if len(meta.ID) == 0 {
		meta.ID = APP_ID
	}
	if meta.Icon == nil {
		meta.Icon = applicationIcon
	}
	meta.Version = caesardisk.ShortVersion
	if len(meta.Name) == 0 {
		meta.Name = APP_TITLE
	}
	meta.Custom["url"] = "https://github.com/lordofscripts"
	meta.Custom["url.text"] = "GitHub"

	// · Help|About dialog
	me := fynex.NewPersonWithImage("Lord of Scripts™", "BScEE, Developer, Writer", developerIcon)
	g.dlgAbout = dlg.NewAboutBox(myWindow, applicationIcon, meta).
		WithText(aboutCONTENT, true, false).
		WithPersonModel(me)

	// · Main Menu Bar
	// 	 File menu : Open - Quit
	menuFileQuit := newMenuItemWithShortcut(fyne.KeyModifierAlt, fyne.KeyQ, "Quit", g.OnAppCloseEvent)
	menuFileOpen := newMenuItemWithShortcut(fyne.KeyModifierAlt, fyne.KeyO, "Open", g.OnFileOpeningEvent)
	menuFileOpen.Disabled = true

	menuTopFile := fyne.NewMenu("File", menuFileOpen, fyne.NewMenuItemSeparator(), menuFileQuit)

	//   Misc menu: KeySchedule
	menuMiscSchedule := newMenuItemWithShortcut(fyne.KeyModifierAlt, fyne.KeyK, "Key schedule", g.showKeySchedule)
	menuTopMisc := fyne.NewMenu("Misc", menuMiscSchedule)

	//   Help menu: About
	menuHelpAbout := newMenuItemWithShortcut(fyne.KeyModifierAlt, fyne.KeyH, "About", g.dlgAbout.ShowDialog)
	menuTopHelp := fyne.NewMenu("Help", menuHelpAbout)

	menuMain := fyne.NewMainMenu(menuTopFile, menuTopMisc, menuTopHelp)
	myWindow.SetMainMenu(menuMain)

	// · Instantiate Gadgets
	// ·· Options Tab
	var this IGadgetParent = g
	g.gadgets.Cipher = newCipherModeGadget(this)
	g.gadgets.Alpha = newAlphabetGadget(this, g.allAlphabets, InitialAlphabetName)
	g.gadgets.Misc = newMiscOptionsGadget(this, g.wheelOpts)

	// ·· Main Tab
	g.gadgets.Wheel = newWheelGadget(g.wheelOpts)
	g.gadgets.Key = newCaesarKeyGadget(this)
	g.gadgets.Offset = newKeyOffsetGadget(this)
	g.gadgets.Data = newSecretDataGadget(myWindow, g)

	// · Instantiate Controllers
	iniAlphabet := (*g.allAlphabets)[InitialAlphabetName]
	// ·· the Cipher Controller does encryption/decryption based on the model
	g.controllers.Cipher = crypto.NewCipherController(iniAlphabet, g)

	// · Create tabs
	tab2 := g.drawOptionsTab()
	tab1 := g.drawMainTab()
	tabs := container.NewAppTabs(
		tab1,
		tab2,
	)

	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(WINDOW_WIDTH, WINDOW_HEIGHT)) // Adjust window size
	myWindow.SetFixedSize(false)

	g.a = myApp
	g.w = myWindow

	return g
}

// set up a bound flag that, when set to true, it would initate a Wheel redraw
func (g *MainGUI) Bind() *MainGUI {
	logx.Visit()

	g.gadgets.Bind()

	DataBindings.Bind()   // bind listeners
	DataBindings.Reload() // refresh so that gadget.Define() gets that value

	// This triggers the listener but does not update the bound widget
	//BoundKeyShift.Set(5)

	return g
}

// show either an error or an informational dialog with a message.
// It implements crypto.IViewNotifier
func (g *MainGUI) Notify(isError bool, message string) {
	logx.Visit()

	if isError {
		errorDlg := dialog.NewError(errors.New(message), g.w)
		errorDlg.Show()
	} else {
		dialog.NewInformation("Warning/Info", message, g.w).Show()
	}
}

// Cascaded actions triggered by a gadget (serviceId) who had gotten activated
// due to user interaction. This contributes to cause-effect across gadgets
// handled by the main application rather than the gadget itself having knowledge
// of other sibling gadgets.
func (g *MainGUI) Cascade(serviceId GadgetId, value any) {
	switch serviceId {
	case GadgetWheel:
		// @todo for now the Wheel does not generate changes. In the future, the
		// we could bind the Mouse-Wheel to more left/right to change the Main Key.

	case GadgetMainKey:
		if v, ok := value.(int); ok {
			logx.OnCascade("mainKey", v)

			g.gadgets.Wheel.Update()
		}

	case GadgetOffset:
		if v, ok := value.(int); ok {
			logx.OnCascade("altKey", v)

			g.gadgets.Wheel.Update()
		}

	case GadgetAlphabet:
		if v, ok := value.(string); ok {
			logx.OnCascade("alphabet", v)
			BoundKeyShift.Set(0)
			BoundKeyOffset.Set(0)

			g.gadgets.Wheel.Update()
		}

	case GadgetCipherMode:
		if v, ok := value.(crypto.CaesarCipherMode); ok {
			logx.OnCascade("cipherMode", v)
			if v == crypto.DidimusMode || v == crypto.PrimusMode {
				BoundKeyOffset.Set(0)
				g.gadgets.Offset.Show()
				g.gadgets.Wheel.CanShowOffset(true)
			} else {
				g.gadgets.Wheel.CanShowOffset(false)
				g.gadgets.Offset.Hide()
			}
		}

	case GadgetOtherOpts:

	}
}

func (g *MainGUI) Run() {
	g.w.ShowAndRun()
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

func (g *MainGUI) drawMainTab() *container.TabItem {
	const TAB_NAME = "Cipher"

	// -- Wheel Gadget
	//		· Provides: IWheelUpdateService
	//		· Uses: SessionController
	g.gadgets.Wheel.Define()

	// -- CaesarKey Gadget
	//		· Provides: ICaesarKeyService
	//		· Uses	  : SessionController
	g.gadgets.Key.Define()

	// -- KeyOffset Gadget
	//	· Provides: IKeyOffsetService
	//	· Uses: 	SessionController
	g.gadgets.Offset.Define()

	// -- SecretData Gadget
	//	· Provides:
	//	· Uses: 	SessionController
	g.gadgets.Data.With(g.controllers.Cipher).Define().Bind()

	// -- Content Layout
	tab1Content := container.New(layout.NewVBoxLayout(),
		g.gadgets.Wheel.Container(),
		g.gadgets.Key.Container(),
		g.gadgets.Offset.Container(),
		g.gadgets.Data.Container(),
	)

	// -- Tab
	tab1 := container.NewTabItem(TAB_NAME, tab1Content)

	g.gadgets.Wheel.PostRender()
	g.gadgets.Key.PostRender()
	g.gadgets.Offset.PostRender()
	g.gadgets.Data.PostRender()

	return tab1
}

func (g *MainGUI) drawOptionsTab() *container.TabItem {
	const TAB_NAME = "Options"

	// -- Cipher Mode Gadget
	//		· Provides: ICipherModeService
	g.gadgets.Cipher.Define()

	// -- Alphabet Gadget
	//		· Provides: IAlphabetService
	g.gadgets.Alpha.Define()

	// -- Rendering: Orthogonality
	g.gadgets.Misc.Define()

	// -- Content Layout
	tab2Content := container.New(layout.NewVBoxLayout(),
		g.gadgets.Cipher.Canvas(),
		g.gadgets.Alpha.Canvas(), // but this gadget uses a card, so we need this.
		//checkOrtho,
		//checkPDU,
		g.gadgets.Misc.Canvas(),
	)

	// -- Tab
	tab2 := container.NewTabItem(TAB_NAME, tab2Content)

	g.gadgets.Cipher.PostRender()
	g.gadgets.Alpha.PostRender()
	g.gadgets.Misc.PostRender()

	return tab2
}

func (g *MainGUI) showKeySchedule() {
	// alphabet instance
	if alphaName, err := BoundAlphaName.Get(); err == nil {
		alpha := caesardisk.AlphabetFactory(alphaName)
		// cipher controller with currently selected alphabet
		ctrl := g.controllers.Cipher.CloneWith(alpha)
		// current encryption parameters
		sm := DataBindings.GetSessionModel()
		param := cipher.CaesarParameters{
			//Alphabet: alpha,
			KeyValue: sm.MainKey.Shift,
			Offset:   sm.Offset,
		}
		// get programmed key schedule for current key/offset setting
		var err error
		var schedule crypto.KeySchedule
		switch g.gadgets.Cipher.GetCipherMode() { // @note Update when new cipher modes
		case crypto.CaesarMode:
			schedule, err = ctrl.GetCaesarSchedule(param.KeyValue)
		case crypto.DidimusMode:
			schedule, err = ctrl.GetDidimusSchedule(param.KeyValue, param.Offset)
		case crypto.FibonacciMode:
			schedule, err = ctrl.GetFibonacciSchedule(param.KeyValue)
		case crypto.PrimusMode:
			schedule, err = ctrl.GetPrimusSchedule(param.KeyValue, param.Offset)
		}
		// any errors?
		if err != nil {
			g.Notify(true, err.Error())
		} else {
			// show schedule to user
			title, _ := BoundCipherModeName.Get()
			modal := NewKeyScheduleViewer(g.w, title, schedule)
			modal.Resize(fyne.NewSize(WINDOW_WIDTH+50, WINDOW_HEIGHT/2))
			modal.Show()
		}
	} else {
		logx.Print(err)
	}

}

// call Bind on all gadgets
func (ag *appGadgets) Bind() {
	ag.Alpha.Bind()
	ag.Cipher.Bind()
	ag.Key.Bind()
	ag.Offset.Bind()
	ag.Wheel.Bind()
	ag.Data.Bind()
	ag.Misc.Bind()
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

// creates a menu item with a shortcut and action callback
func newMenuItemWithShortcut(modifier fyne.KeyModifier, keyname fyne.KeyName, label string, callback MenuItemActionCB) *fyne.MenuItem {
	shortcut := &desktop.CustomShortcut{
		KeyName:  keyname,
		Modifier: modifier,
	}
	menuItem := fyne.NewMenuItem(label, callback)
	menuItem.Shortcut = shortcut
	return menuItem
}

/* ----------------------------------------------------------------
 *				E v e n t   C a l l b a c k s
 *-----------------------------------------------------------------*/

func (g *MainGUI) OnStarted() {
	g.Bind()
}

// called whenever the GUI application is closed
func (g *MainGUI) OnAppCloseEvent() {
	g.a.Quit()
}

// @todo whenever File|Open is clicked. It displays a file selector
// and loads the file contents accordingly.
func (g *MainGUI) OnFileOpeningEvent() {
	logx.Print("Not implemented")
}
