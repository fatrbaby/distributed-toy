package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
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
			r.notify(patch{
				Removed: []patchEntry{
					patchEntry{
						Name: service.Name,
						URL:  ServersURL,
					},
				},
			})

			r.mutex.Lock()
			r.services = append(r.services[:i], r.services[i+1:]...)
			r.mutex.Unlock()
			return nil
		}
	}

	return fmt.Errorf("service at URL: %s not found", url)
}

func (r *registry) notify(fullPatch patch) {
	for _, svc := range r.services {
		go func(s Service) {
			for _, required := range s.RequiredServices {
				p := newPatch()

				shouldUpdate := false

				for _, entry := range fullPatch.Added {
					if entry.Name == required {
						p.Added = append(p.Added, entry)
						shouldUpdate = true
					}
				}
				for _, entry := range fullPatch.Removed {
					if entry.Name == required {
						fmt.Printf("remove service: %v\n", required)
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

func (r *registry) heartbeat(freq time.Duration) {
	for {
		var wg sync.WaitGroup
		for _, service := range r.services {
			wg.Add(1)
			go func(svc Service) {
				defer wg.Done()

				success := true

				for retries := 0; retries < 3; retries++ {
					response, err := http.Get(svc.HeartbeatURL)

					if err != nil {
						log.Println(err)
					} else if response.StatusCode == http.StatusOK {
						log.Printf("Heartbeat check passed for: %v", svc.Name)

						if !success {
							r.add(svc)
						}

						break
					}

					log.Printf("Heartbeat check failed for: %v", svc.Name)

					if success {
						success = false
						r.remove(svc.URL)
					}

					time.Sleep(1 * time.Second)
				}
			}(service)

			wg.Wait()
			time.Sleep(freq)
		}
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
