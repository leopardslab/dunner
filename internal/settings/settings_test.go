package settings

import (
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

func TestInit(t *testing.T) {
	Init()
	defaultSettings := map[string]interface{}{
		"dunnertaskfile":   ".dunner.yaml",
		"dotenvfile":       ".env",
		"globallogfile":    "/var/log/dunner/logs/",
		"workingdirectory": "./",
		"async":            false,
		"verbose":          false,
		"dry-run":          false,
		"dockerapiversion": "1.39",
		"no-color":         false,
	}

	if !reflect.DeepEqual(viper.AllSettings(), defaultSettings) {
		t.Fatal("Default not equal to as expected")
	}
}
