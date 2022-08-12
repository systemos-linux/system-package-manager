package archive

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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
	data, err := os.ReadFile(path)
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

func Build(path string) (*Archive, error) {
	a := Load(filepath.Join(path, "pkg.json"))
	if a == nil {
		return nil, fmt.Errorf("couldn't find pkg.json, aborting")
	}
	files := make([]string, 0)
	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	_ = os.MkdirAll("./output", 0755)
	fp, e := os.OpenFile(filepath.Join("output", fmt.Sprintf("%s-%s.sys", strings.ToLower(a.Name), a.Version)), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	if e != nil {
		return nil, e
	}
	e = create_archive(files, fp)
	if e != nil {
		return nil, e
	}
	return a, e
}

func create_archive(files []string, buf io.Writer) error {
	gw := gzip.NewWriter(buf)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, file := range files {
		err := archive_add(tw, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func archive_add(tw *tar.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}
	header.Name = filename

	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}

	return nil
}
