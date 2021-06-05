package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func RegisterService(r Registrar) error {
	buff := new(bytes.Buffer)
	encoder := json.NewEncoder(buff)

	err := encoder.Encode(r)

	if err != nil {
		return err
	}

	response, err := http.Post(ServersURL, "application/json", buff)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK  {
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

	if response.StatusCode != http.StatusOK  {
		return fmt.Errorf("failed to register service, responded with code: %v", response.StatusCode)
	}

	return nil
}
