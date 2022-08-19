package service

import (
	"fmt"
	"os"
)

func response(msg string, isError, isWarning, newline bool) {
	if isError {
		fmt.Printf("=> yatt[error]: %s", msg)
	} else if isWarning {
		fmt.Printf("=> yatt[warning]: %s", msg)
	}

	if isError {
		os.Exit(1)
	}
}
