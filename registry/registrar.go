package registry

type Registrar struct {
	Name ServiceName
	URL  string
}

type ServiceName string
