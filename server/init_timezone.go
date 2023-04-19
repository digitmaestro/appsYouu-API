package server

import (
	"os"
	"runtime"
	"time"
)

// InitTimezoneOnWindows This provides as fix where TZ environment variable is not recognized when running on windows.
func InitTimezoneOnWindows() {
	timezone, timezoneSet := os.LookupEnv("TZ")
	if runtime.GOOS == "windows" && timezoneSet {
		time.Local, _ = time.LoadLocation(timezone)
	}
}
