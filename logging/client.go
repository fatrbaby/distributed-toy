package logging

import (
	"bytes"
	"fmt"
	"github.com/fatrbaby/distributed-toy/registry"
	"log"
	"net/http"
)

func UseClientLogger(serviceURL string, clientService registry.ServiceName) {
	log.SetPrefix(fmt.Sprintf("[%v]", clientService))
	log.SetFlags(0)
	log.SetOutput(&clientLogger{
		url: serviceURL,
	})
}

type clientLogger struct {
	url string
}

func (c clientLogger) Write(p []byte) (n int, err error) {
	b := bytes.NewBuffer(p)

	response, err := http.Post(c.url, "text/plain", b)

	if err != nil {
		return 0, err
	}

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to send logging message. service responed with status %v", response.StatusCode)
	}

	return len(p), nil
}
