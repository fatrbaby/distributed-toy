package registry

type Registrar struct {
	Name             ServiceName
	URL              string
	RequiredServices []ServiceName
	ServiceUpdateURL string
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
