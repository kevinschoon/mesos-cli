package section

import (
	ui "github.com/gizak/termui"
	"github.com/vektorlab/toplib"
	"github.com/vektorlab/toplib/cursor"
	"github.com/vektorlab/toplib/sample"
	"github.com/vektorlab/toplib/toggle"
)

type Samples struct {
	Namespace sample.Namespace
	SortField string
	Fields    []string
	Cursor    *cursor.Cursor
	SortMenu  *toplib.Menu
	Toggles   toggle.Toggles
}

func NewSamples(ns sample.Namespace, fields ...string) *Samples {
	return &Samples{
		Namespace: ns,
		SortField: "ID",
		Fields:    fields,
		Cursor:    cursor.NewCursor(),
		Toggles:   toggle.NewToggles(&toggle.Toggle{Name: "sort"}),
	}
}

func (s Samples) ordered(rec *toplib.Recorder) []*sample.Sample {
	latest := rec.Latest(s.Namespace)
	sample.Sort(s.SortField, latest)
	return latest
}

func (d Samples) Name() string { return string(d.Namespace) }

func (d Samples) Handlers(opts toplib.Options) map[string]func(ui.Event) {
	return map[string]func(ui.Event){}
}

func (s Samples) Grid(opts toplib.Options) *ui.Grid {
	samples := s.ordered(opts.Recorder)
	rows := [][]string{s.Fields}
	for _, sample := range samples {
		rows = append(rows, sample.Strings(s.Fields))
	}
	table := ui.NewTable()
	table.Rows = rows
	table.Separator = false
	table.Border = false
	table.SetSize()
	table.Analysis()
	//table.BgColors[s.Cursor.IDX(opts.Recorder.Items())] = ui.ColorRed
	l := ui.NewList()
	l.Items = s.Fields
	l.Height = 30
	l.Width = 25
	if s.Toggles.State("sort") {
		return ui.NewGrid(
			ui.NewRow(
				ui.NewCol(3, 0, l),
				ui.NewCol(9, 0, table),
			),
		)
	}
	return ui.NewGrid(
		ui.NewRow(
			ui.NewCol(12, 0, table),
		),
	)
}
