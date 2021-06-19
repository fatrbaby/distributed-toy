package registry

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	ServerPort = ":4600"
	ServersURL = "http://localhost" + ServerPort + "/services"
)

var center = registry{
	services: make([]Service, 0),
	mutex:    new(sync.RWMutex),
}

var once sync.Once

func SetupRegistry()  {
	once.Do(func() {
		go center.heartbeat(3 * time.Second)
	})
}

type Server struct{}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received")

	switch r.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var registrar Service
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
