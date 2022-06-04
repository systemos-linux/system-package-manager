package main

import (
	"log"
)

func CheckIfError(err error) {
	if err != nil {
		log.Fatalln("Error: ", err)
	}
}
