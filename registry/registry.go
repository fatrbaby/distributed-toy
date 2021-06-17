package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type registry struct {
	registrars []Registrar
	mutex      *sync.RWMutex
}

func (rs *registry) add(registrar Registrar) error {
	rs.mutex.Lock()
	rs.registrars = append(rs.registrars, registrar)
	rs.mutex.Unlock()

	return rs.sendRequiredServices(registrar)
}

func (rs *registry) remove(url string) error {
	for i, r := range rs.registrars {
		if r.URL == url {
			rs.mutex.Lock()
			rs.registrars = append(rs.registrars[:i], rs.registrars[:i+1]...)
			rs.mutex.Unlock()
			return nil
		}
	}

	return fmt.Errorf("service at URL: %s not found", url)
}

func (rs *registry) sendRequiredServices(registrar Registrar) error {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	var p patch

	for _, r := range rs.registrars {
		for _, registered := range registrar.RequiredServices {
			if registered == r.Name {
				p.Added = append(p.Added, patchEntry{
					Name: r.Name,
					URL:  r.ServiceUpdateURL,
				})
			}
		}
	}

	return rs.sendPatch(p, registrar.ServiceUpdateURL)
}

func (rs *registry) sendPatch(p patch, url string) error {
	data, err := json.Marshal(p)

	if err != nil {
		return err
	}

	_, err = http.Post(url, "application/json", bytes.NewBuffer(data))

	if err != nil {
		return err
	}

	return nil
}
