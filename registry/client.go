package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func RegisterService(r Registration) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(r)
	if err != nil {
		return err
	}

	res, err := http.Post(ServiceURL, "application/json", buf)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register service. Registry service "+"responded with code %v", res.StatusCode)
	}
	return nil
}

func ShutdownService(url string) error {
	req, err := http.NewRequest(http.MethodDelete, ServiceURL, bytes.NewBuffer([]byte(url)))
	if err != nil {
		return nil
	}
	req.Header.Add("Content-Type", "text/plain")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to deregister servie."+
			"service responded with code %v", res.StatusCode)
	}
	return nil
}
