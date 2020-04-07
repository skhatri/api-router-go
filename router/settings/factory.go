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

var routeSettings *_RouteSettings
var empty = &EmptySettings{}

func ApplySettings(settings *string) error {
	if routeSettings != nil {
		return errors.New("settings has already been loaded")
	}
	settingsMux.Lock()
	defer settingsMux.Unlock()
	settingsFile := ""
	if settings == nil {
		if rs := os.Getenv("ROUTE_SETTINGS"); rs != "" {
			settingsFile = rs
		}
	} else {
		settingsFile = *settings
	}
	errorMode := true
	if settingsFile == "" {
		errorMode = false
		settingsFile = "router.json"
	}
	var localSettings _RouteSettings
	if _, statErr := os.Stat(settingsFile); statErr == nil {
		f, err := os.OpenFile(settingsFile, os.O_RDONLY, os.ModeType)
		if err != nil {
			return errors.New(fmt.Sprintf("settings file %s could not be read. %s", settingsFile, err.Error()))
		}
		err = json.NewDecoder(f).Decode(&localSettings)
		if err != nil {
			return err
		}
		log.Println("settings file loaded from", settingsFile)
	} else {
		if errorMode {
			return errors.New(fmt.Sprintf("settings file %s is not readable", settingsFile))
		}
	}
	if routeSettings == nil {
		routeSettings = &localSettings
	}
	return nil
}

func GetSettings() RouteSettings {
	if routeSettings == nil {
		return empty
	}
	return routeSettings
}
