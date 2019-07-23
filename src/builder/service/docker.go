package service

import (
	"builder/constant"
	"io/ioutil"
	"net/http"
)

// DockerService is relative docker services
type DockerService struct{}

func init() {

}

// GetCatalog returns docker registry catalog
func (d *DockerService) GetCatalog() []byte {
	// needs admin logon
	// needs token

	resp, err := http.Get(basicinfo.GetRegistryURL(constant.PathRegistryCatalog))
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
