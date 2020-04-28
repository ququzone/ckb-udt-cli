package cmd

import (
	"fmt"
	"os"
)

func Fatalf(format string, v ...interface{}) {
	fmt.Printf(format, v)
	os.Exit(1)
}
