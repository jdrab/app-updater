package serviceconfig

import (
	"runtime"
)

type ServiceConfig struct {
	ServiceName string `json:"name"`
	AppName     string `json:"appName"`
	Version     string `json:"version"`
}

// Load now you know what
func Load() ServiceConfig {
	configuration := make(map[string]ServiceConfig)
	defaultVersion := "app-updater-0.1.0"

	var winConfig = ServiceConfig{
		ServiceName: "my-service",
		AppName:     "Sample Client App.exe",
		Version:     defaultVersion,
	}

	var linuxConfig = ServiceConfig{
		ServiceName: "my-service",
		AppName:     "sample-client-app",
		Version:     defaultVersion,
	}

	configuration["windows"] = winConfig
	configuration["linux"] = linuxConfig

	return configuration[runtime.GOOS]
}
