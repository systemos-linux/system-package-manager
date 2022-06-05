package common

import (
	"fmt"
	"os"
)

func CheckIfError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}
