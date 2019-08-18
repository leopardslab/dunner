package settings

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/leopardslab/dunner/internal"
	"github.com/spf13/viper"
)

func TestInit(t *testing.T) {
	Init()
	fmt.Print(viper.AllSettings())
	defaultSettings := map[string]interface{}{
		"dunnertaskfile":   internal.DefaultDunnerTaskFileName,
		"dotenvfile":       ".env",
		"globallogfile":    "/var/log/dunner/logs/",
		"workingdirectory": "./",
		"async":            false,
		"verbose":          false,
		"dry-run":          false,
		"force-pull":       false,
		"dockerapiversion": "1.39",
		"no-color":         false,
	}

	if !reflect.DeepEqual(viper.AllSettings(), defaultSettings) {
		t.Fatal("Default not equal to as expected")
	}
}
