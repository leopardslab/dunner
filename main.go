package main

import (
	"github.com/leopardslab/dunner/cmd"
	"github.com/leopardslab/dunner/internal/settings"
	G "github.com/leopardslab/dunner/pkg/global"
)

var version string

func main() {
	settings.Init()
	G.VERSION = version
	cmd.Execute()
}
