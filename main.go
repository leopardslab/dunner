package main

import (
	"github.com/leopardslab/Dunner/cmd"
	"github.com/leopardslab/Dunner/internal/settings"
)

func main() {
	settings.Initialize()
	cmd.Execute()
}
