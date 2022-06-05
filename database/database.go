package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"systemos.org/pkg/archive"
	"systemos.org/pkg/common"
	"systemos.org/pkg/semver"
)

type Database struct {
	Repository *Repository
	Installed  map[string]*archive.Archive
}

func New() *Database {
	return &Database{
		Repository: new_repository(),
		Installed:  load_installed(),
	}
}

func (d *Database) Update() {
	d.Repository.Update()
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

		// Fetch the archive struct for this package and version
		p, err := d.Repository.Get(p_name, p_version)
		common.CheckIfError(err)

		if p != nil {
			// we either fetched their version
			// or we fetched "latest" and now
			// have a version, either way, reuse.
			p_version = p.Version
			// see if it's already installed
			i := d.Installed[p_name]

			// if it is
			if i != nil {
				v := i.Version
				// compare the installed version with the fetched version.
				// compare will return +1 if fetched version is greater than
				// installed version. Otherwise it will return -1 if less
				// and 0 if same.
				if semver.Compare(p_version, v) > 0 {
					// the fetched version is greater, let's install it.
					err := p.Install()
					// any errors during install?
					if err == nil {
						// no errors, let's track it.
						d.Installed[p_name] = p
						// TODO: do we uninstall previous versions?
						// leave them? track them for removal? or unused?
						// (shrug)
					} else {
						log.Fatalf("Error installing package %s@%s\n%v", p_name, p_version, err)
					}
				}
			}
		}
	}
}

func load_installed() (installed map[string]*archive.Archive) {
	err := filepath.Walk("/var/pkg/installed", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if filepath.Ext(path) == "json" {
				a := archive.Load(path)
				installed[a.Name] = a
			}
		}
		return nil
	})
	common.CheckIfError(err)
	if err != nil {
		if errors.Is(err, syscall.ENOENT) {
			fmt.Println("Unable to determine installed packages. /var/pkg/installed is not a directory.")
		} else if errors.Is(err, syscall.EPERM) {
			fmt.Println("Unable to determine installed packages. Please run as an administrator.")
			fmt.Printf("sudo %v\n", os.Args)
		} else {
			fmt.Println("Unable to determine installed packages.")
			fmt.Printf("%T", err)
		}
	}
	return installed
}
