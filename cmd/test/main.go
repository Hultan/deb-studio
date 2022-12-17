package main

import (
	"fmt"

	"github.com/hultan/deb-studio/internal/engine"
)

func main() {
	ws, err := engine.Open("/home/per/installs/softtube")
	if err != nil {
		fmt.Println(err)
	}

	for _, w := range ws.Versions {
		fmt.Println("Version:", w.Name)
		for _, a := range w.Architectures {
			fmt.Println("Architecture", a.Name)
		}
	}
}
