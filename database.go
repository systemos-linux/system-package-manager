package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"systemos.org/pkg/semver"
)

type Database struct {
	repository *Repository
	installed map[string]*Archive
}

func new_database() *Database {
	return &Database {
		repository: new_repository(),
		installed: load_installed(),
	}
}

func (d *Database) Update() {
	d.repository.Update()
}

func (d *Database) InstallPackages(packages []string) {
	for _, p_name := range packages {
		p_name = strings.ToLower(p_name)
		var p_version string
		if strings.ContainsRune(p_name, '@') {
			// package name contains version info
			// parse it out
			parts := strings.Split(p_name, "@")
			p_name = parts[0]
			p_version = parts[1]
		} else {
			p_version = "latest"
		}
		p := d.repository.find(p_name, p_version)
		if p != nil {
			
			i := d.installed[p_name+"@"+p_version]
			if i != nil {
				v_list := i.Version
				semver.Compare()
			}
		}
	}
}

func load_installed() (installed map[string]&Archive) {
	err := filepath.WalkDir("/var/pkg/installed", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if filepath.Ext(path) == "json" {
				archive := load_archive(path)
				installed[archive.Name+"@"+archive.Version] = archive
			}
		}
		return nil
	})
	CheckIfError(err)
	return installed
}