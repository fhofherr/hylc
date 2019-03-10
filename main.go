package main

import (
	"fmt"

	"github.com/fhofherr/hylc/cmd"
)

func main() {
	fmt.Printf("Build time: %s\n", cmd.BuildTime)
	fmt.Printf("Git hash  : %s\n", cmd.GitHash)
	fmt.Printf("Version   : %s\n", cmd.Version)
}
