package toplib

import (
	ui "github.com/gizak/termui"
)

var (
	padding  = 2
	minWidth = 30
)

type Menu struct {
	ui.Block
	Items       []string
	TextFgColor ui.Attribute
	TextBgColor ui.Attribute
	Selectable  bool
	CursorPos   int
	Handlers    map[string]func(ui.Event)
}

func NewMenu(items ...string) *Menu {
	m := &Menu{
		Block:       *ui.NewBlock(),
		Items:       items,
		TextFgColor: ui.ThemeAttr("par.text.fg"),
		TextBgColor: ui.ThemeAttr("par.text.bg"),
		Selectable:  true,
		CursorPos:   0,
		Handlers:    make(map[string]func(ui.Event)),
	}
	m.Width, m.Height = calcSize(items)
	return m
}

func (m *Menu) Buffer() ui.Buffer {
	var cell ui.Cell
	buf := m.Block.Buffer()

	for n, item := range m.Items {
		x := padding
		for _, ch := range item {
			// invert bg/fg colors on currently selected row
			if m.Selectable && n == m.CursorPos {
				cell = ui.Cell{Ch: ch, Fg: m.TextBgColor, Bg: m.TextFgColor}
			} else {
				cell = ui.Cell{Ch: ch, Fg: m.TextFgColor, Bg: m.TextBgColor}
			}
			buf.Set(x, n+padding, cell)
			x++
		}
	}

	return buf
}

func (m *Menu) Up() string {
	if m.CursorPos > 0 {
		m.CursorPos--
	}
	return m.Items[m.CursorPos]
}

func (m *Menu) Down() string {
	if m.CursorPos < (len(m.Items) - 1) {
		m.CursorPos++
	}
	return m.Items[m.CursorPos]
}

// return width and height based on menu items
func calcSize(items []string) (w, h int) {
	h = len(items) + (padding * 2)

	w = minWidth
	for _, s := range items {
		if len(s) > w {
			w = len(s)
		}
	}
	w += (padding * 2)

	return w, h
}
