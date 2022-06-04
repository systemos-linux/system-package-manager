package main

import (
	"fmt"
	"os"
	"strings"
)
var (
	db *Database
)

func main() {
	args := os.Args[1:]
	fmt.Println("System Package Manager")
	command := strings.Trim(strings.ToLower(args[0]))

	db = new_database()
	db.init()

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

func update_package_database() {
	
	db.update()
}

func list_installed_packages() {

}

func install_packages([]string packages) {

}

func uninstall_packages([]string packages) {

}

func info_on_package(string package) {

}

func cleanup_packages() {

}

func build_package(string package) {

}

func unpack_package(string package) {

}
