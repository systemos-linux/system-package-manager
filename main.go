package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/schollz/progressbar/v3"
	"systemos.org/pkg/archive"
	"systemos.org/pkg/database"
)

var (
	db *database.Database
)

func main() {

	args := os.Args[1:]
	if len(args) < 1 {
		print_help()
		return
	}
	fmt.Println("System Package Manager")
	command := strings.Trim(strings.ToLower(args[0]), " ")

	db = database.New()

	switch command {
	case "update", "u":
		update_package_database()
	case "install", "i":
		install_packages(args[1:])
	case "uninstall", "rm":
		uninstall_packages(args[1:])
	case "info":
		info_on_package(args[1:][0])
	case "list", "l":
		list_installed_packages()
	case "clean", "c":
		cleanup_packages()
	case "build", "b":
		build_package(args[1:][0])
	case "unpack":
		unpack_package(args[1:][0])
	}
}

func print_help() {
	fmt.Println("SystemOS GNU/Linux Package Manager")
	fmt.Println("-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
	fmt.Println("Usage:")
	fmt.Println("\tpkg [option] <[arguments]>")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Printf("\t%-20s-\t%s\n", "build <dir>", "Builds a package from the current directory.")
	fmt.Printf("\t%-20s-\t%s\n", "clean", "Cleans the directory of a package build.")
	fmt.Printf("\t%-20s-\t%s\n", "info <pkg>", "Shows info of a package from the repository. Use './' prefix for local package.")
	fmt.Printf("\t%-20s-\t%s\n", "install <...pkg>", "Installs a package from the repository. Use './' prefix for local packages.")
	fmt.Printf("\t%-20s-\t%s\n", "list", "Lists installed packages on this machine.")
	fmt.Printf("\t%-20s-\t%s\n", "uninstall <...pkg>", "Uninstall removes packages on this machine.")
	fmt.Printf("\t%-20s-\t%s\n", "unpack <pkg> <dir>", "Unpack a package to the specified directory.")
	fmt.Printf("\t%-20s-\t%s\n", "update", "Updates the local package repository with new package definitions.")
}
func update_package_database() {

	db.Update()
}

func list_installed_packages() {

}

func install_packages(package_names []string) {
	bar := progressbar.Default(
		-1,
		"Installing Packages...",
	)
	for _, p := range package_names {
		bar.Describe(fmt.Sprintf("Package: %s", p))
		bar.Add(1)
		pkg := archive.LoadDeb(p)
		fmt.Printf("Package: %v\n", pkg.Control)
	}
	bar.Finish()
}

func uninstall_packages(package_names []string) {

}

func info_on_package(package_name string) {

}

func cleanup_packages() {

}

func build_package(package_name string) {

}

func unpack_package(package_name string) {

}
