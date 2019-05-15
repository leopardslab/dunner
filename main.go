package main

import (
	"github.com/leopardslab/dunner/cmd"
	G "github.com/leopardslab/dunner/pkg/global"
)

var version string

func main() {
	G.VERSION = version
	cmd.Execute()
}
