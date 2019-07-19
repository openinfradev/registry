package service

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// HTTPSampleService is http sample service
type HTTPSampleService struct{}

func init() {

}

// GetDaum returns daum homepage
func (h *HTTPSampleService) GetDaum() []byte {

	resp, err := http.Get("http://daum.net")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return r
}

// PostTest is sample function. so not works
func (h *HTTPSampleService) PostTest() string {

	reqBody := bytes.NewBufferString("post request body")
	resp, err := http.Post("http://test", "text/plain", reqBody)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(r)
}
