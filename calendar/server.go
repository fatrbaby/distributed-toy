package calendar

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	api = "http://api.tianapi.com/txapi/lunar/index"
	key = "01de3c7b6fadebfba692e57891d87b45"
)

func RegisterHandler() {
	http.HandleFunc("/lunar", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		query := request.URL.Query()

		response, err := http.Get(fmt.Sprintf("%s?key=%s&date=%s", api, key, query.Get("date")))

		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		bytes, _ := ioutil.ReadAll(response.Body)

		writer.Header().Add("content-type", "application/json")
		writer.Write(bytes)
	})
}
