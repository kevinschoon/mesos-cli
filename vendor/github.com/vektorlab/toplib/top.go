package toplib

import (
	ui "github.com/gizak/termui"
	"github.com/gizak/termui/extra"
	"github.com/vektorlab/toplib/sample"
	"time"
)

// Section returns a renderable ui.Grid
type Section interface {
	Name() string
	Grid(Options) *ui.Grid
	Handlers(Options) map[string]func(ui.Event)
}

type Options struct {
	Recorder *Recorder
	Render   func()
}

// Top renders Sections which are periodically updated
type Top struct {
	Exit     chan bool
	Errors   chan error
	Samples  chan []*sample.Sample
	Recorder *Recorder // Holds samples
	Sections []Section
	Tabpane  *extra.Tabpane
	Grid     *ui.Grid
	section  int
	Options  Options
}

func NewTop(sections []Section) *Top {
	top := &Top{
		Exit:     make(chan bool),
		Errors:   make(chan error),
		Samples:  make(chan []*sample.Sample),
		Recorder: NewRecorder(),
		Sections: sections,
		Tabpane:  extra.NewTabpane(),
		Grid:     ui.NewGrid(),
	}
	top.Options = Options{
		Recorder: top.Recorder,
		Render: func() {
			render(top)
		},
	}
	return top
}

func handlers(top *Top) {
	ui.DefaultEvtStream.ResetHandlers()
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		top.Exit <- true
	})
	ui.Handle("/sys/kbd/j", func(ui.Event) {
		top.Tabpane.SetActiveLeft()
		render(top)
	})
	ui.Handle("/sys/kbd/k", func(ui.Event) {
		top.Tabpane.SetActiveRight()
		render(top)
	})
	for path, fn := range top.Sections[top.section].Handlers(top.Options) {
		ui.Handle(path, fn)
	}
}

func render(top *Top) {
	handlers(top)
	tabs := []extra.Tab{}
	for _, section := range top.Sections {
		grid := section.Grid(top.Options)
		grid.Width = ui.TermWidth()
		grid.Align()
		tab := extra.NewTab(section.Name())
		tab.AddBlocks(grid)
		tabs = append(tabs, *tab)
	}
	top.Tabpane.SetTabs(tabs...)
	top.Tabpane.Width = ui.TermWidth()
	top.Tabpane.Align()
	ui.Clear()
	ui.Render(top.Tabpane)
}

func collect(top *Top, fn sample.SampleFunc) {
	for {
		samples, err := fn()
		if err != nil {
			top.Errors <- err
			return
		}
		top.Samples <- samples
		time.Sleep(500 * time.Millisecond)
	}
}

func Run(top *Top, funcs ...sample.SampleFunc) (err error) {
	if err = ui.Init(); err != nil {
		return err
	}

	defer ui.Close()

	for _, fn := range funcs {
		go collect(top, fn)
	}

	go func() {
		for {
			select {
			case err = <-top.Errors:
				ui.StopLoop()
				break
			case <-top.Exit:
				ui.StopLoop()
				break
			case samples := <-top.Samples:
				if len(samples) > 0 {
					top.Recorder.Load(samples[0].Namespace(), samples)
					handlers(top)
					render(top)
				}
			}
		}
	}()
	render(top)
	ui.Loop()
	return err
}
