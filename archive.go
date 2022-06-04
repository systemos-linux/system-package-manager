package main

import (
	"ioutil"
	"json"
)

type Archive struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	Url         string   `json:"url"`
	Version     string   `json:"version"`
	Maintainer  string   `json:"maintainer"`
}

func load_archive(path string) *Archive {
	data, err := ioutil.ReadFile(path)
	CheckIfError(err)

	var archive Archive
	err = json.Unmarshal(data, &archive)
	CheckIfError(err)
	return &archive
}
