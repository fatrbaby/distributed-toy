package registry

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

const (
	ServerPort = ":4600"
	ServersURL = "http://localhost" + ServerPort + "/services"
)

var center = registry{
	registrars: make([]Registrar, 0),
	mutex:      new(sync.RWMutex),
}

type Service struct{}

func (s Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received")

	switch r.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var registrar Registrar
		err := decoder.Decode(&registrar)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("adding service: %v with URL: %s \n", registrar.Name, registrar.URL)
		err = center.add(registrar)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodDelete:
		payload, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		url := string(payload)
		log.Printf("Removing service at URL: %s", url)
		err = center.remove(url)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
