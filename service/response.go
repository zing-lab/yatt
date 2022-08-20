package service

import (
	"fmt"
)

func response(msg string, isError, isWarning, newline bool) error {
	if isError {
		fmt.Printf("=> yatt[error]: %s", msg)
	} else if isWarning {
		fmt.Printf("=> yatt[warning]: %s", msg)
	}
	return nil
}
