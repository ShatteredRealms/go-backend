package main

import (
	"fmt"
	"github.com/ShatteredRealms/go-backend/cmd/updater/updater"
)

func main() {
	err := SetupConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	updater.Execute()
}
