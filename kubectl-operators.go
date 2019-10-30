package main

import (
	"os"
	"github.com/komish/kubectl-operators/cmd"
)

func main() {
	os.Exit(cmd.Run())
}