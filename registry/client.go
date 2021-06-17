package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

const (
	ServiceLogging  = ServiceName("service.logging")
	ServiceCalendar = ServiceName("service.calendar")
)

var prov = providers{
	services: make(map[ServiceName][]string),
	mutex:    new(sync.RWMutex),
}

type providers struct {
	services map[ServiceName][]string
	mutex    *sync.RWMutex
}

func (that *providers) Update(p patch) {
	that.mutex.Lock()
	defer that.mutex.Unlock()

	for _, entry := range p.Added {
		if _, ok := that.services[entry.Name]; !ok {
			that.services[entry.Name] = make([]string, 0)
		}

		that.services[entry.Name] = append(that.services[entry.Name], entry.URL)
	}

	for _, entry := range p.Removed {
		if urls, ok := that.services[entry.Name]; ok {
			for i := range urls {
				if urls[i] == entry.URL {
					that.services[entry.Name] = append(urls[:i], urls[:i+1]...)
				}
			}
		}
	}
}

func GetProvider(name ServiceName) (string, error) {
	return prov.get(name)
}

func (that *providers) get(name ServiceName) (string, error) {
	services, ok := that.services[name]

	if !ok {
		return "", fmt.Errorf("no provider available for service %v", name)
	}

	index := int(rand.Float32() * float32(len(services)))

	return services[index], nil
}

func RegisterService(r Registrar) error {
	updateURL, err := url.Parse(r.ServiceUpdateURL)

	if err != nil {
		return err
	}

	http.Handle(updateURL.Path, &serviceUpdateHandler{})

	buff := new(bytes.Buffer)
	encoder := json.NewEncoder(buff)

	err = encoder.Encode(r)

	if err != nil {
		return err
	}

	response, err := http.Post(ServersURL, "application/json", buff)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register service, responded with code: %v", response.StatusCode)
	}

	return nil
}

func ShutdownService(url string) error {
	request, err := http.NewRequest(http.MethodDelete, ServersURL, bytes.NewBuffer([]byte(url)))

	if err != nil {
		return err
	}

	request.Header.Add("Content-type", "text/plain")

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register service, responded with code: %v", response.StatusCode)
	}

	return nil
}

type serviceUpdateHandler struct{}

func (s serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var p patch

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	prov.Update(p)
}
