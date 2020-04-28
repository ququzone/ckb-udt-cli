package main

import (
	"fmt"
	"os"

	"github.com/ququzone/ckb-udt-cli/cmd"
	"github.com/ququzone/ckb-udt-cli/config"
)

func main() {
	err := config.Init()
	if err != nil {
		fmt.Printf("read config file error: %v", err)
		os.Exit(1)
	}

	cmd.Execute()
}
