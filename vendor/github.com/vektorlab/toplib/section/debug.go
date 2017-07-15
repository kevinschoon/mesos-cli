package section

import (
	"fmt"
	ui "github.com/gizak/termui"
	"github.com/vektorlab/toplib"
	"github.com/vektorlab/toplib/sample"
	"runtime"
	"sort"
)

// Debug is a section that displays debug information
// which may be useful for developing toplib.
type Debug struct {
	Namespaces []sample.Namespace
}

func (d Debug) Name() string { return "debug" }

func (d Debug) Handlers(opts toplib.Options) map[string]func(ui.Event) {
	return map[string]func(ui.Event){}
}

func (d *Debug) Grid(opts toplib.Options) *ui.Grid {
	p := ui.NewPar(fmt.Sprintf("Samples Loaded: %d, Go Routines: %d",
		opts.Recorder.Counter, runtime.NumGoroutine()))
	p.Height = 3
	p.Width = 10
	l := ui.NewList()
	l.BorderLabel = "Handlers"
	l.Items = listHandlers()
	l.Width = 25
	l.Height = len(l.Items) + 1
	n := ui.NewList()
	n.BorderLabel = "Namespaces"
	n.Items = []string{}
	for _, ns := range d.Namespaces {
		latest := opts.Recorder.Latest(ns)
		n.Items = append(n.Items, fmt.Sprintf("(%s) Unique samples: %d", ns, len(latest)))
	}
	n.Height = len(n.Items) + 2
	return ui.NewGrid(
		ui.NewRow(
			ui.NewCol(6, 0, p),
		),
		ui.NewRow(
			ui.NewCol(6, 0, n),
			ui.NewCol(6, 0, l),
		),
	)
}

func listHandlers() []string {
	strs := []string{}
	for path, _ := range ui.DefaultEvtStream.Handlers {
		strs = append(strs, path)
	}
	sort.Strings(strs)
	return strs
}
