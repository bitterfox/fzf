package tui

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Types of user action
type EventType int

const (
	Rune EventType = iota

	CtrlA
	CtrlB
	CtrlC
	CtrlD
	CtrlE
	CtrlF
	CtrlG
	CtrlH
	Tab
	CtrlJ
	CtrlK
	CtrlL
	CtrlM
	CtrlN
	CtrlO
	CtrlP
	CtrlQ
	CtrlR
	CtrlS
	CtrlT
	CtrlU
	CtrlV
	CtrlW
	CtrlX
	CtrlY
	CtrlZ
	ESC
	CtrlSpace

	// https://apple.stackexchange.com/questions/24261/how-do-i-send-c-that-is-control-slash-to-the-terminal
	CtrlBackSlash
	CtrlRightBracket
	CtrlCaret
	CtrlSlash

	Invalid
	Resize
	Mouse
	DoubleClick
	LeftClick
	RightClick

	BTab
	BSpace

	Del
	PgUp
	PgDn

	Up
	Down
	Left
	Right
	Home
	End
	Insert

	SUp
	SDown
	SLeft
	SRight

	F1
	F2
	F3
	F4
	F5
	F6
	F7
	F8
	F9
	F10
	F11
	F12

	Change
	BackwardEOF
	Loaded

	AltBS

	AltUp
	AltDown
	AltLeft
	AltRight

	AltSUp
	AltSDown
	AltSLeft
	AltSRight

	Alt
	CtrlAlt
)

func (t EventType) AsEvent() Event {
	return Event{t, 0, nil}
}

func (t EventType) Int() int {
	return int(t)
}

func (t EventType) Byte() byte {
	return byte(t)
}

func (e Event) Comparable() Event {
	// Ignore MouseEvent pointer
	return Event{e.Type, e.Char, nil}
}

func Key(r rune) Event {
	return Event{Rune, r, nil}
}

func AltKey(r rune) Event {
	return Event{Alt, r, nil}
}

func CtrlAltKey(r rune) Event {
	return Event{CtrlAlt, r, nil}
}

const (
	doubleClickDuration = 500 * time.Millisecond
)

type Color int32

func (c Color) IsDefault() bool {
	return c == colDefault
}

func (c Color) is24() bool {
	return c > 0 && (c&(1<<24)) > 0
}

type ColorAttr struct {
	Color Color
	Attr  Attr
	BackgroundColor Color
}

func NewColorAttr() ColorAttr {
	return ColorAttr{Color: colUndefined, Attr: AttrUndefined, BackgroundColor: colUndefined}
}

const (
	colUndefined Color = -2
	colDefault   Color = -1
)

const (
	colBlack Color = iota
	colRed
	colGreen
	colYellow
	colBlue
	colMagenta
	colCyan
	colWhite
)

type FillReturn int

const (
	FillContinue FillReturn = iota
	FillNextLine
	FillSuspend
)

type ColorPair struct {
	fg   Color
	bg   Color
	attr Attr
}

func HexToColor(rrggbb string) Color {
	r, _ := strconv.ParseInt(rrggbb[1:3], 16, 0)
	g, _ := strconv.ParseInt(rrggbb[3:5], 16, 0)
	b, _ := strconv.ParseInt(rrggbb[5:7], 16, 0)
	return Color((1 << 24) + (r << 16) + (g << 8) + b)
}

func NewColorPair(fg Color, bg Color, attr Attr) ColorPair {
	return ColorPair{fg, bg, attr}
}

func (p ColorPair) Fg() Color {
	return p.fg
}

func (p ColorPair) Bg() Color {
	return p.bg
}

func (p ColorPair) Attr() Attr {
	return p.attr
}

func (p ColorPair) HasBg() bool {
	return p.attr&Reverse == 0 && p.bg != colDefault ||
		p.attr&Reverse > 0 && p.fg != colDefault
}

func (p ColorPair) merge(other ColorPair, except Color) ColorPair {
	dup := p
	dup.attr = dup.attr.Merge(other.attr)
	if other.fg != except {
		dup.fg = other.fg
	}
	if other.bg != except {
		dup.bg = other.bg
	}
	return dup
}

func (p ColorPair) WithAttr(attr Attr) ColorPair {
	dup := p
	dup.attr = dup.attr.Merge(attr)
	return dup
}

func (p ColorPair) MergeAttr(other ColorPair) ColorPair {
	return p.WithAttr(other.attr)
}

func (p ColorPair) Merge(other ColorPair) ColorPair {
	return p.merge(other, colUndefined)
}

func (p ColorPair) MergeNonDefault(other ColorPair) ColorPair {
	return p.merge(other, colDefault)
}

type ColorTheme struct {
	Colored      bool
	Input        ColorAttr
	Disabled     ColorAttr
	Fg           ColorAttr
	Bg           ColorAttr
	PreviewFg    ColorAttr
	PreviewBg    ColorAttr
	DarkBg       ColorAttr
	Gutter       ColorAttr
	Prompt       ColorAttr
	Match        ColorAttr
	Current      ColorAttr
	CurrentMatch ColorAttr
	Spinner      ColorAttr
	Info         ColorAttr
	Cursor       ColorAttr
	Selected     ColorAttr
	Header       ColorAttr
	Border       ColorAttr
}

type Event struct {
	Type       EventType
	Char       rune
	MouseEvent *MouseEvent
}

type MouseEvent struct {
	Y      int
	X      int
	S      int
	Left   bool
	Down   bool
	Double bool
	Mod    bool
}

type BorderShape int

const (
	BorderNone BorderShape = iota
	BorderRounded
	BorderSharp
	BorderHorizontal
	BorderVertical
	BorderTop
	BorderBottom
	BorderLeft
	BorderRight
)

type BorderStyle struct {
	shape       BorderShape
	horizontal  rune
	vertical    rune
	topLeft     rune
	topRight    rune
	bottomLeft  rune
	bottomRight rune
}

type BorderCharacter int

func MakeBorderStyle(shape BorderShape, unicode bool) BorderStyle {
	if unicode {
		if shape == BorderRounded {
			return BorderStyle{
				shape:       shape,
				horizontal:  '─',
				vertical:    '│',
				topLeft:     '╭',
				topRight:    '╮',
				bottomLeft:  '╰',
				bottomRight: '╯',
			}
		}
		return BorderStyle{
			shape:       shape,
			horizontal:  '─',
			vertical:    '│',
			topLeft:     '┌',
			topRight:    '┐',
			bottomLeft:  '└',
			bottomRight: '┘',
		}
	}
	return BorderStyle{
		shape:       shape,
		horizontal:  '-',
		vertical:    '|',
		topLeft:     '+',
		topRight:    '+',
		bottomLeft:  '+',
		bottomRight: '+',
	}
}

func MakeTransparentBorder() BorderStyle {
	return BorderStyle{
		shape:       BorderRounded,
		horizontal:  ' ',
		vertical:    ' ',
		topLeft:     ' ',
		topRight:    ' ',
		bottomLeft:  ' ',
		bottomRight: ' '}
}

type Renderer interface {
	Init()
	Pause(clear bool)
	Resume(clear bool, sigcont bool)
	Clear()
	RefreshWindows(windows []Window)
	Refresh()
	Close()

	GetChar() Event

	MaxX() int
	MaxY() int

	NewWindow(top int, left int, width int, height int, preview bool, borderStyle BorderStyle) Window
}

type Window interface {
	Top() int
	Left() int
	Width() int
	Height() int

	Refresh()
	FinishFill()
	Close()

	X() int
	Y() int
	Enclose(y int, x int) bool

	Move(y int, x int)
	MoveAndClear(y int, x int)
	Print(text string)
	CPrint(color ColorPair, text string)
	Fill(text string) FillReturn
	CFill(fg Color, bg Color, attr Attr, text string) FillReturn
	Erase()
}

type FullscreenRenderer struct {
	theme        *ColorTheme
	mouse        bool
	forceBlack   bool
	prevDownTime time.Time
	clickY       []int
}

func NewFullscreenRenderer(theme *ColorTheme, forceBlack bool, mouse bool) Renderer {
	r := &FullscreenRenderer{
		theme:        theme,
		mouse:        mouse,
		forceBlack:   forceBlack,
		prevDownTime: time.Unix(0, 0),
		clickY:       []int{}}
	return r
}

var (
	Default16 *ColorTheme
	Dark256   *ColorTheme
	Light256  *ColorTheme

	ColPrompt               ColorPair
	ColNormal               ColorPair
	ColInput                ColorPair
	ColDisabled             ColorPair
	ColMatch                ColorPair
	ColCursor               ColorPair
	ColCursorEmpty          ColorPair
	ColSelected             ColorPair
	ColCurrent              ColorPair
	ColCurrentMatch         ColorPair
	ColCurrentCursor        ColorPair
	ColCurrentCursorEmpty   ColorPair
	ColCurrentSelected      ColorPair
	ColCurrentSelectedEmpty ColorPair
	ColSpinner              ColorPair
	ColInfo                 ColorPair
	ColHeader               ColorPair
	ColBorder               ColorPair
	ColPreview              ColorPair
	ColPreviewBorder        ColorPair
)

func EmptyTheme() *ColorTheme {
	return &ColorTheme{
		Colored:      true,
		Input:        ColorAttr{colUndefined, AttrUndefined, colUndefined},
		Disabled:     ColorAttr{colUndefined, AttrUndefined, colUndefined},
		Fg:           ColorAttr{colUndefined, AttrUndefined, colUndefined},
		Bg:           ColorAttr{colUndefined, AttrUndefined, colUndefined},
		PreviewFg:    ColorAttr{colUndefined, AttrUndefined, colUndefined},
		PreviewBg:    ColorAttr{colUndefined, AttrUndefined, colUndefined},
		DarkBg:       ColorAttr{colUndefined, AttrUndefined, colUndefined},
		Gutter:       ColorAttr{colUndefined, AttrUndefined, colUndefined},
		Prompt:       ColorAttr{colUndefined, AttrUndefined, colUndefined},
		Match:        ColorAttr{colUndefined, AttrUndefined, colUndefined},
		Current:      ColorAttr{colUndefined, AttrUndefined, colUndefined},
		CurrentMatch: ColorAttr{colUndefined, AttrUndefined, colUndefined},
		Spinner:      ColorAttr{colUndefined, AttrUndefined, colUndefined},
		Info:         ColorAttr{colUndefined, AttrUndefined, colUndefined},
		Cursor:       ColorAttr{colUndefined, AttrUndefined, colUndefined},
		Selected:     ColorAttr{colUndefined, AttrUndefined, colUndefined},
		Header:       ColorAttr{colUndefined, AttrUndefined, colUndefined},
		Border:       ColorAttr{colUndefined, AttrUndefined, colUndefined}}
}

func NoColorTheme() *ColorTheme {
	return &ColorTheme{
		Colored:      false,
		Input:        ColorAttr{colDefault, AttrRegular, colUndefined},
		Disabled:     ColorAttr{colDefault, AttrRegular, colUndefined},
		Fg:           ColorAttr{colDefault, AttrRegular, colUndefined},
		Bg:           ColorAttr{colDefault, AttrRegular, colUndefined},
		PreviewFg:    ColorAttr{colDefault, AttrRegular, colUndefined},
		PreviewBg:    ColorAttr{colDefault, AttrRegular, colUndefined},
		DarkBg:       ColorAttr{colDefault, AttrRegular, colUndefined},
		Gutter:       ColorAttr{colDefault, AttrRegular, colUndefined},
		Prompt:       ColorAttr{colDefault, AttrRegular, colUndefined},
		Match:        ColorAttr{colDefault, Underline, colUndefined},
		Current:      ColorAttr{colDefault, Reverse, colUndefined},
		CurrentMatch: ColorAttr{colDefault, Reverse | Underline, colUndefined},
		Spinner:      ColorAttr{colDefault, AttrRegular, colUndefined},
		Info:         ColorAttr{colDefault, AttrRegular, colUndefined},
		Cursor:       ColorAttr{colDefault, AttrRegular, colUndefined},
		Selected:     ColorAttr{colDefault, AttrRegular, colUndefined},
		Header:       ColorAttr{colDefault, AttrRegular, colUndefined},
		Border:       ColorAttr{colDefault, AttrRegular, colUndefined}}
}

func errorExit(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(2)
}

func init() {
	Default16 = &ColorTheme{
		Colored:      true,
		Input:        ColorAttr{colDefault, AttrUndefined, colUndefined},
		Disabled:     ColorAttr{colUndefined, AttrUndefined, colUndefined},
		Fg:           ColorAttr{colDefault, AttrUndefined, colUndefined},
		Bg:           ColorAttr{colDefault, AttrUndefined, colUndefined},
		PreviewFg:    ColorAttr{colUndefined, AttrUndefined, colUndefined},
		PreviewBg:    ColorAttr{colUndefined, AttrUndefined, colUndefined},
		DarkBg:       ColorAttr{colBlack, AttrUndefined, colUndefined},
		Gutter:       ColorAttr{colUndefined, AttrUndefined, colUndefined},
		Prompt:       ColorAttr{colBlue, AttrUndefined, colUndefined},
		Match:        ColorAttr{colGreen, AttrUndefined, colUndefined},
		Current:      ColorAttr{colYellow, AttrUndefined, colUndefined},
		CurrentMatch: ColorAttr{colGreen, AttrUndefined, colUndefined},
		Spinner:      ColorAttr{colGreen, AttrUndefined, colUndefined},
		Info:         ColorAttr{colWhite, AttrUndefined, colUndefined},
		Cursor:       ColorAttr{colRed, AttrUndefined, colUndefined},
		Selected:     ColorAttr{colMagenta, AttrUndefined, colUndefined},
		Header:       ColorAttr{colCyan, AttrUndefined, colUndefined},
		Border:       ColorAttr{colBlack, AttrUndefined, colUndefined}}
	Dark256 = &ColorTheme{
		Colored:      true,
		Input:        ColorAttr{colDefault, AttrUndefined, colUndefined},
		Disabled:     ColorAttr{colUndefined, AttrUndefined, colUndefined},
		Fg:           ColorAttr{colDefault, AttrUndefined, colUndefined},
		Bg:           ColorAttr{colDefault, AttrUndefined, colUndefined},
		PreviewFg:    ColorAttr{colUndefined, AttrUndefined, colUndefined},
		PreviewBg:    ColorAttr{colUndefined, AttrUndefined, colUndefined},
		DarkBg:       ColorAttr{236, AttrUndefined, colUndefined},
		Gutter:       ColorAttr{colUndefined, AttrUndefined, colUndefined},
		Prompt:       ColorAttr{110, AttrUndefined, colUndefined},
		Match:        ColorAttr{108, AttrUndefined, colUndefined},
		Current:      ColorAttr{254, AttrUndefined, colUndefined},
		CurrentMatch: ColorAttr{151, AttrUndefined, colUndefined},
		Spinner:      ColorAttr{148, AttrUndefined, colUndefined},
		Info:         ColorAttr{144, AttrUndefined, colUndefined},
		Cursor:       ColorAttr{161, AttrUndefined, colUndefined},
		Selected:     ColorAttr{168, AttrUndefined, colUndefined},
		Header:       ColorAttr{109, AttrUndefined, colUndefined},
		Border:       ColorAttr{59, AttrUndefined, colUndefined}}
	Light256 = &ColorTheme{
		Colored:      true,
		Input:        ColorAttr{colDefault, AttrUndefined, colUndefined},
		Disabled:     ColorAttr{colUndefined, AttrUndefined, colUndefined},
		Fg:           ColorAttr{colDefault, AttrUndefined, colUndefined},
		Bg:           ColorAttr{colDefault, AttrUndefined, colUndefined},
		PreviewFg:    ColorAttr{colUndefined, AttrUndefined, colUndefined},
		PreviewBg:    ColorAttr{colUndefined, AttrUndefined, colUndefined},
		DarkBg:       ColorAttr{251, AttrUndefined, colUndefined},
		Gutter:       ColorAttr{colUndefined, AttrUndefined, colUndefined},
		Prompt:       ColorAttr{25, AttrUndefined, colUndefined},
		Match:        ColorAttr{66, AttrUndefined, colUndefined},
		Current:      ColorAttr{237, AttrUndefined, colUndefined},
		CurrentMatch: ColorAttr{23, AttrUndefined, colUndefined},
		Spinner:      ColorAttr{65, AttrUndefined, colUndefined},
		Info:         ColorAttr{101, AttrUndefined, colUndefined},
		Cursor:       ColorAttr{161, AttrUndefined, colUndefined},
		Selected:     ColorAttr{168, AttrUndefined, colUndefined},
		Header:       ColorAttr{31, AttrUndefined, colUndefined},
		Border:       ColorAttr{145, AttrUndefined, colUndefined}}
}

func initTheme(theme *ColorTheme, baseTheme *ColorTheme, forceBlack bool) {
	if forceBlack {
		theme.Bg = ColorAttr{colBlack, AttrUndefined, colUndefined}
	}

	o := func(a ColorAttr, b ColorAttr) ColorAttr {
		c := a
		if b.Color != colUndefined {
			c.Color = b.Color
		}
		if b.Attr != AttrUndefined {
			c.Attr = b.Attr
		}
		if b.BackgroundColor != colUndefined {
			c.BackgroundColor = b.BackgroundColor
		}
		return c
	}
	theme.Input = o(baseTheme.Input, theme.Input)
	theme.Disabled = o(theme.Input, o(baseTheme.Disabled, theme.Disabled))
	theme.Fg = o(baseTheme.Fg, theme.Fg)
	theme.Bg = o(baseTheme.Bg, theme.Bg)
	theme.PreviewFg = o(theme.Fg, o(baseTheme.PreviewFg, theme.PreviewFg))
	theme.PreviewBg = o(theme.Bg, o(baseTheme.PreviewBg, theme.PreviewBg))
	theme.DarkBg = o(baseTheme.DarkBg, theme.DarkBg)
	theme.Gutter = o(theme.DarkBg, o(baseTheme.Gutter, theme.Gutter))
	theme.Prompt = o(baseTheme.Prompt, theme.Prompt)
	theme.Match = o(baseTheme.Match, theme.Match)
	theme.Current = o(baseTheme.Current, theme.Current)
	theme.CurrentMatch = o(baseTheme.CurrentMatch, theme.CurrentMatch)
	theme.Spinner = o(baseTheme.Spinner, theme.Spinner)
	theme.Info = o(baseTheme.Info, theme.Info)
	theme.Cursor = o(baseTheme.Cursor, theme.Cursor)
	theme.Selected = o(baseTheme.Selected, theme.Selected)
	theme.Header = o(baseTheme.Header, theme.Header)
	theme.Border = o(baseTheme.Border, theme.Border)

	initPalette(theme)
}

func initPalette(theme *ColorTheme) {
	pair := func(fg, bg ColorAttr) ColorPair {
		if fg.Color == colDefault && (fg.Attr&Reverse) > 0 {
			bg.Color = colDefault
		}
		if fg.BackgroundColor != colDefault && fg.BackgroundColor != colUndefined {
			bg.Color = fg.BackgroundColor
		}
		return ColorPair{fg.Color, bg.Color, fg.Attr}
	}
	blank := theme.Fg
	blank.Attr = AttrRegular

	ColPrompt = pair(theme.Prompt, theme.Bg)
	ColNormal = pair(theme.Fg, theme.Bg)
	ColInput = pair(theme.Input, theme.Bg)
	ColDisabled = pair(theme.Disabled, theme.Bg)
	ColMatch = pair(theme.Match, theme.Bg)
	ColCursor = pair(theme.Cursor, theme.Gutter)
	ColCursorEmpty = pair(blank, theme.Gutter)
	ColSelected = pair(theme.Selected, theme.Gutter)
	ColCurrent = pair(theme.Current, theme.DarkBg)
	ColCurrentMatch = pair(theme.CurrentMatch, theme.DarkBg)
	ColCurrentCursor = pair(theme.Cursor, theme.DarkBg)
	ColCurrentCursorEmpty = pair(blank, theme.DarkBg)
	ColCurrentSelected = pair(theme.Selected, theme.DarkBg)
	ColCurrentSelectedEmpty = pair(blank, theme.DarkBg)
	ColSpinner = pair(theme.Spinner, theme.Bg)
	ColInfo = pair(theme.Info, theme.Bg)
	ColHeader = pair(theme.Header, theme.Bg)
	ColBorder = pair(theme.Border, theme.Bg)
	ColPreview = pair(theme.PreviewFg, theme.PreviewBg)
	ColPreviewBorder = pair(theme.Border, theme.PreviewBg)
}
