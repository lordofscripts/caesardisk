# Designing the GUI

 Mark  | Description |
| :---: | :---:       |
| 🛈    | The purpose of the widget |
| ⟫     | Condition that should be fulfilled for the widget to be enabled |
| ↯     | When the widget is clicked, i.e. buttons, menuItem    |
| 🢱     | Widgets down the chain that get affected by the `OnChange` event      |
| 🢰     | Widgets down the chain that get affected by the `OnValidate` event   |
| Ⓤ     | Unbound Data model that gets updated | 
| Ⓑ     | Bound Data model that gets updated |

## Data Model - 

* CurrentKeySequencer (IKeySequencer)
* CurrentTranscoder (ITranscoder)
* CaesarParams.Alphabet
* CaesarParams.KeyValue
* CaesarParams.KeyOffset

## Main Menu Bar

· File|Open
	🛈 Opens a text file and loads it into the input TextEntry widget
	
· File|Quit
	🛈 Terminate the Graphical User Interface main window loop and return to CLI


## Main Canvas

### ✅ AlphabetGadget (Options Tab)
Displays the encoding alphabet (English, Spanish, German, Punctuation, etc.)
Provides: IAlphabetService
Uses    : IParamProvider

· {AlphabetGadget}.Selector
	🛈 User selects encoding alphabet
	Ⓤ {Data.Alphabet} gets updated
	Ⓤ {Data.KeyValue} set to 0
	↯ [changed] {DataGadget}.Clear
	↯ [changed] {KeyGadget} set key to zero (Clear)
	↯ [changed] {WheelGadget}
	
· {AlphabetGadget}.Alphabet
	🛈 A label that displays all the characters in the selected alphabet

### ✅ CipherModeGadget (Options Tab)
Displays the alphabet selection
Provides: ICipherModeService
Uses    : none

· {CipherModeGadget}.Selector
	🛈 User selects encoding algorithm (Caesar, Didimus, Fibonacci)
	↯ [value = Didimus] {OffsetGadget}.Show()
	↯ [value != Didimus] {OffsetGadget}.Hide()

### ✅ WheelGadget (Main Tab)
Displays the Caesar encoder wheel. It gets updated when the key is set.
Provides: IWheelUpdateService
Uses    : IAlphabetService

· {WheelGadget}.Image
	🛈 Displays the Caesar wheel with the current Key or Offset
	↯ Toggle between main `KeyShift` and `KeyOffset` disks.
	
### ✅ CaesarKeyGadget (Main Tab)
Handles all UI actions related to the main Caesar key.
Provides: ICaesarKeyService
Uses    : IAlphabetService, IParamProvider

⟫ {CipherModeGadget.Selected} is any of Caesar, Didimus or Fibonacci

· {CaesarKeyGadget}.LabelKey
	🛈 Displays the character corresponding to `KeyShift` for selected alphabet

· {CaesarKeyGadget}.Slider
	🛈 The user slides left & right to set the `KeyShift` value (Caesar main key)
	🢱 {CaesarKeyGadget}.LabelKey
	🢱 {CaesarKeyGadget}.LabelKeyShift
	Ⓤ {Data.KeyValue} 

· {KeyGadget}.LabelKeyShift
	🛈  Displays the integer correspoding to the `KeyShift` selected in Slider
	
### ✅ OffsetGadget (Main Tab)
Handles every UI action related to the Offset shift used in Didimus & Fibonacci.
Provides: IKeyOffsetService
Uses    : IAlphabetService, IParamsService

⟫ {CipherModeGadget.Selected} is Didimus
 
· {OffsetGadget}.LabelKey
	🛈 Displays the character corresponding to `KeyShift` for selected alphabet

· {OffsetGadget}.Slider
	🛈 The user slides left & right to set the `KeyShift` value (Caesar main key)
	🢱 {OffsetGadget}.LabelKey
	🢱 {OffsetGadget}.LabelKeyShift
	Ⓤ {Data.Offset} 

· {OffsetGadget}.LabelKeyShift
	🛈  Displays the integer correspoding to the `KeyShift` selected in Slider
	
### DataGadget (Main Tab)
Displays the input and output texts as well as the action buttons.
Provides: 
Uses    : 

· {DataGadget}.InputText (widget.MultiLineEntry)

· {DataGadget}.OutputText (widget.MultiLineEntry)

· {DataGadget}.Encode (widget.Button)
	⟫ {DataGadget}.InputText length > 0
	↯ {DataGadget}.OutputText updated with encrypted text result

· {DataGadget}.Decode (widget.Button)
	⟫ {DataGadget}.InputText length > 0
	↯ {DataGadget}.OutputText updated with decrypted text result

· {DataGadget}.CLR (widget.Button)
	⟫ {DataGadget}.InputText length > 0
	↯ {DataGadget}.InputText clear
	↯ {DataGadget}.OutputText clear

· {DataGadget}.Exchange (widget.Button)
	⟫ {DataGadget}.InputText length > 0
	↯ swap contents of {DataGadget}.InputText and {DataGadget}.OutputText
	