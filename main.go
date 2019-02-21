package main

import (
	"github.com/leopardslab/Dunner/cmd"
	G "github.com/leopardslab/Dunner/pkg/global"
)

var version string = ""

func main() {
	G.VERSION = version
	cmd.Execute()
}
