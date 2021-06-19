package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type registry struct {
	services []Service
	mutex    *sync.RWMutex
}

func (r *registry) add(service Service) error {
	r.mutex.Lock()
	r.services = append(r.services, service)
	r.mutex.Unlock()

	err := r.sendRequiredServices(service)
	r.notify(patch{
		Added: []patchEntry{
			patchEntry{
				Name: service.Name,
				URL:  service.URL,
			},
		},
	})

	return err
}

func (r *registry) remove(url string) error {
	for i, service := range r.services {
		if service.URL == url {
			r.mutex.Lock()
			r.services = append(r.services[:i], r.services[:i+1]...)
			r.mutex.Unlock()
			r.notify(patch{
				Removed: []patchEntry{
					{
						Name: service.Name,
						URL:  ServersURL,
					},
				},
			})
			return nil
		}
	}

	return fmt.Errorf("service at URL: %s not found", url)
}

func (r *registry) notify(fullPatch patch) {
	for _, svc := range r.services {
		go func(s Service) {
			for _, required := range s.RequiredServices {
				p := patch{
					Added: []patchEntry{},
					Removed: []patchEntry{},
				}
				shouldUpdate := false

				for _, entry := range fullPatch.Added {
					if entry.Name == required {
						p.Added = append(p.Added, entry)
						shouldUpdate = true
					}
				}
				for _, entry := range fullPatch.Removed {
					if entry.Name == required {
						p.Removed = append(p.Removed, entry)
						shouldUpdate = true
					}
				}

				if shouldUpdate {
					err := r.sendPatch(p, s.UpdateURL)

					if err != nil {
						fmt.Println(err)
						return
					}
				}
			}
		}(svc)
	}
}

func (r *registry) sendRequiredServices(service Service) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var p patch

	for _, r := range r.services {
		for _, registered := range service.RequiredServices {
			if registered == r.Name {
				p.Added = append(p.Added, patchEntry{
					Name: r.Name,
					URL:  r.UpdateURL,
				})
			}
		}
	}

	return r.sendPatch(p, service.UpdateURL)
}

func (r *registry) sendPatch(p patch, url string) error {
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
