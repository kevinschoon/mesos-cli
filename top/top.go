package top

import (
	ui "github.com/gizak/termui"
	"github.com/vektorlab/toplib"
	//"time"
)

func initTop() *toplib.Top {
	t := toplib.NewTop()
	var (
		cursor  = toplib.NewCursor()
		toggles = toplib.NewToggles(
			&toplib.Toggle{Name: "sort"},
		)

		sortMenu = toplib.NewMenu("ID", "FRAMEWORK", "STATE", "CPU", "MEM", "GPU", "DISK")
		table    = toplib.NewTable("ID", "FRAMEWORK", "STATE", "CPU", "MEM", "GPU", "DISK")
	)

	defaultView := toplib.NewView(func() []*ui.Row {
		return []*ui.Row{
			toplib.NewHeader().Row(),
			// Main section
			ui.NewRow(
				ui.NewCol(12, 0, table.Buffers(t.Recorder, cursor)...),
			),
			// Bottom toggles
			ui.NewRow(
				ui.NewCol(12, 0, toggles.Buffers()...),
			),
		}
	})

	defaultView.Handlers["/sys/kbd/<up>"] = func(ui.Event) {
		if cursor.Up(t.Recorder.Samples()) {
			t.Render()
		}
	}

	defaultView.Handlers["/sys/kbd/<down>"] = func(ui.Event) {
		if cursor.Down(t.Recorder.Samples()) {
			t.Render()
		}
	}

	defaultView.Handlers["/sys/kbd/s"] = func(ui.Event) {
		if toggles.Toggle("sort", true) {
			t.Views.Set("sort")
		} else {
			t.Views.Set("default")
		}
		t.Render()
	}

	sortView := toplib.NewView(func() []*ui.Row {
		return []*ui.Row{
			ui.NewRow(
				ui.NewCol(3, 0, sortMenu),
				ui.NewCol(9, 0, table.Buffers(t.Recorder, cursor)...),
			),
			ui.NewRow(
				ui.NewCol(12, 0, toggles.Buffers()...),
			),
		}
	})

	sortView.Handlers["/sys/kbd/<up>"] = func(ui.Event) {
		t.Recorder.SortField = sortMenu.Up()
		t.Render()
	}
	sortView.Handlers["/sys/kbd/<down>"] = func(ui.Event) {
		t.Recorder.SortField = sortMenu.Down()
		t.Render()
	}

	t.Views.Add("default", defaultView)
	t.Views.Add("sort", sortView)
	t.Views.Set("default")

	return t
}

// TODO: Update toplib so samples can be sent individually
/*
func collect(client *Client) ([]*toplib.Sample, error) {
	samples := []*toplib.Sample{}
	tasks, err := client.Tasks(TaskFilterAll)
	if err != nil {
		return nil, err
	}
	for _, task := range tasks {
		sample := toplib.NewSample(task.GetTaskId().GetValue())
		sample.SetString("FRAMEWORK", task.GetFrameworkId().GetValue())
		sample.SetString("STATE", task.GetState().String())
		sample.SetFloat64("CPU", FilterScalar(task.GetResources(), "cpus"))
		sample.SetFloat64("MEM", FilterScalar(task.GetResources(), "mem"))
		sample.SetFloat64("GPU", FilterScalar(task.GetResources(), "gpu"))
		sample.SetFloat64("DISK", FilterScalar(task.GetResources(), "disk"))
		samples = append(samples, sample)
	}
	return samples, nil
}

func RunTop(client *Client) (err error) {
	top := initTop()
	tick := time.NewTicker(1500 * time.Millisecond)
	go func() {
	loop:
		for {
			select {
			case <-top.Exit:
				close(top.Samples)
				break loop
			case <-tick.C:
				samples, err := collect(client)
				if err != nil {
					break loop
				}
				top.Samples <- samples
			}
		}
		tick.Stop()
	}()
	return toplib.Run(top)
}
*/
