package toplib

import (
	"fmt"
	"time"

	ui "github.com/gizak/termui"
)

type Header struct {
	Time  *ui.Par
	Count *ui.Par
}

func NewHeader() *Header {
	return &Header{
		Time:  headerPar(timeStr()),
		Count: headerPar("-"),
	}
}

func (h *Header) Row() *ui.Row {
	h.Time.Text = timeStr()
	return ui.NewRow(
		ui.NewCol(2, 0, h.Time),
		ui.NewCol(2, 0, h.Count),
	)
}

func (c *Header) SetCount(val int) {
	c.Count.Text = fmt.Sprintf("%d samples", val)
}

func timeStr() string {
	return time.Now().Local().Format("15:04:05 MST")
}

func headerPar(s string) *ui.Par {
	p := ui.NewPar(fmt.Sprintf(" %s", s))
	p.Border = false
	p.Height = 1
	p.Width = 20
	//p.TextFgColor = ui.ColorDefault
	//p.TextBgColor = ui.ColorWhite
	//p.Bg = ui.ColorWhite
	return p
}
