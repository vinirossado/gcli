package main

import (
	"github.com/vinirossado/gcli/cmd/gcli"
)

func main() {
	err := gcli.Execute()
	if err != nil {
		return
	}
}
