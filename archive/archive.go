package archive

import (
	"encoding/json"
	"io/ioutil"

	"github.com/systemos-linux/go-debian/deb"
	"systemos.org/pkg/common"
)

type Archive struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Keywords     []string `json:"keywords"`
	Url          string   `json:"url"`
	Homepage     string   `json:"homepage"`
	Version      string   `json:"version"`
	Dependencies []string `json:"dependencies"`
	Maintainer   string   `json:"maintainer"`
}

func Load(path string) *Archive {
	data, err := ioutil.ReadFile(path)
	common.CheckIfError(err)

	var archive Archive
	err = json.Unmarshal(data, &archive)
	common.CheckIfError(err)
	return &archive
}

func (a *Archive) Install() error {
	return nil
}

func (a *Archive) Uninstall() error {
	return nil
}

func LoadDeb(path string) *deb.Deb {
	deb, _, err := deb.LoadFile(path)
	common.CheckIfError(err)
	return deb
}
