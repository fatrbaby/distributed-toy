package logger

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var logger *log.Logger

type fileLogger string

// implements io.Writer
func (fl fileLogger) Write(data []byte) (n int, err error) {
	file, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)

	if err != nil {
		return 0, err
	}

	defer file.Close()

	return file.Write(data)
}

func Run(destination string)  {
	logger = log.New(fileLogger(destination), "logger: ", log.LstdFlags)
}

func RegisterHandlers()  {
	http.HandleFunc("/log", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodPost:
			msg, err := ioutil.ReadAll(request.Body)

			if err != nil || len(msg) == 0 {
				writer.WriteHeader(http.StatusBadRequest)
				return
			}

			write(msg)
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}

func write(data []byte)  {
	logger.Printf("%v\n", string(data))
}
