package project

import (
	"os"
	"strconv"
)

// devMode indicates whether the project is running in development mode.
var devMode, _ = strconv.ParseBool(os.Getenv("DEV_MODE"))

// DevMode indicates whether the project is running in development mode.
func DevMode() bool {
	return devMode
}
