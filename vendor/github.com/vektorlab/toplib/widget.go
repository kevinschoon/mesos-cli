package toplib

import (
	ui "github.com/gizak/termui"
	"github.com/vektorlab/toplib/toggle"
	"strings"
)

type ToggleMenu struct {
	ui.Block
	ToggleFgColor ui.Attribute
	ToggleBgColor ui.Attribute
	Toggles       toggle.Toggles
}

func NewToggleMenu(toggles toggle.Toggles) *ToggleMenu {
	return &ToggleMenu{
		Toggles:       toggles,
		ToggleFgColor: ui.ThemeAttr("list.item.fg"),
		ToggleBgColor: ui.ThemeAttr("list.item.bg"),
	}
}

func (tm *ToggleMenu) names() []string {
	names := []string{}
	for _, toggle := range tm.Toggles {
		names = append(names, toggle.Name)
	}
	return names
}

func (tm *ToggleMenu) Buffer() ui.Buffer {
	buf := tm.Block.Buffer()
	cs := ui.DefaultTxBuilder.Build(strings.Join(tm.names(), "\n"), tm.ToggleFgColor, tm.ToggleBgColor)
	i, j := 0, 0
	for i < len(cs) {
		if cs[j].Ch == '\n' {
		}
	}

	return buf
}
