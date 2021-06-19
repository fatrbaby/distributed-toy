package registry

type Service struct {
	Name             ServiceName
	URL              string
	RequiredServices []ServiceName
	UpdateURL        string
}

type ServiceName string

type patchEntry struct {
	Name ServiceName
	URL  string
}

type patch struct {
	Added   []patchEntry `json:"added"`
	Removed []patchEntry `json:"removed"`
}

func newPatch() patch {
	return patch{
		Added:   []patchEntry{},
		Removed: []patchEntry{},
	}
}
