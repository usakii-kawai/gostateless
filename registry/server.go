package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

const ServerPort = ":8080"
const ServiceURL = "http://127.0.0.1" + ServerPort + "/services"

type registry struct {
	registrations []Registration
	mutex         *sync.Mutex
}

func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()
	return nil
}

func (r *registry) remove(url string) error {
	for i := range reg.registrations {
		if reg.registrations[i].ServiceURL == url {
			r.mutex.Lock()
			// fmt.Println("--------------")
			// bs, _ := json.Marshal(r.registrations)
			// fmt.Println(string(bs))
			// fmt.Println("--------------")
			reg.registrations = append(reg.registrations[:i], reg.registrations[i+1:]...)
			// fmt.Println("--------------")
			// bs, _ = json.Marshal(r.registrations)
			// fmt.Println(string(bs))
			// fmt.Println("--------------")

			fmt.Println("delete service:", url)
			rest, _ := json.Marshal(r.registrations)
			fmt.Println("rest service now:", string(rest))

			r.mutex.Unlock()
			return nil
		}
	}
	return fmt.Errorf("service at url %s not found", url)
}

var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.Mutex),
}

type RegistryService struct{}

func (s RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received")

	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		var r Registration
		err := dec.Decode(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Adding service: %v with URL: %v\n", r.ServiceName, r.ServiceURL)
		err = reg.add(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	case http.MethodDelete:
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		url := string(payload)
		log.Printf("Removing service at url: %s", url)
		err = reg.remove(url)
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
