package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
)

var settingsMux = sync.Mutex{}

var routeSettings RouteSettings
var empty = &EmptySettings{}

func init() {

	settingsMux.Lock()
	defer settingsMux.Unlock()

	settingsFile := "router.json"

	if rs := os.Getenv("ROUTE_SETTINGS"); rs != "" {
		settingsFile = rs
	} else {
		if _, derr := os.Stat(settingsFile); derr != nil {
			routeSettings = empty
			return
		}
	}

	var localSettings _RouteSettings
	if _, statErr := os.Stat(settingsFile); statErr == nil {
		f, err := os.OpenFile(settingsFile, os.O_RDONLY, os.ModeType)
		if err != nil {
			errors.New(fmt.Sprintf("settings file %s could not be read. %s", settingsFile, err.Error()))
		}
		err = json.NewDecoder(f).Decode(&localSettings)
		if err != nil {
			return
		}
		log.Println("settings file loaded from", settingsFile)
		routeSettings = &localSettings
	} else {
		routeSettings = empty
	}
}

func GetSettings() RouteSettings {
	return routeSettings
}
