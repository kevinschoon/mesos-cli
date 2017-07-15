package toggle

type Toggle struct {
	Name  string
	state bool
}

func (t *Toggle) Toggle() bool {
	t.state = t.state == false
	return t.state
}

func (t *Toggle) Off() {
	t.state = false
}

type Toggles []*Toggle

func NewToggles(toggles ...*Toggle) Toggles {
	return Toggles(toggles)
}

func (toggles Toggles) State(name string) (state bool) {
	for _, toggle := range toggles {
		if toggle.Name == name {
			state = toggle.state
		}
	}
	return state
}

func (toggles Toggles) Toggle(name string, exclusive bool) (state bool) {
	for _, toggle := range toggles {
		if toggle.Name == name {
			state = toggle.Toggle()
		}
		if exclusive && toggle.Name != name {
			toggle.Off()
		}
	}
	return state
}
