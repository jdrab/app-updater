package config

import "runtime"

// ServiceConfig
// type part struct {
// 	Service     string `json:"name"`
// 	Application string `json:"appName"`
// }

type config struct {
	Service string `json:"service"`
	App     string `json:"app"`
}

// Load now you know what
func Load() config {
	configuration := make(map[string]config)

	// default platform configuration, can be changed using makefile
	configuration["linux"] = config{
		Service: "my-service",
		App:     "sample-client-app",
	}

	configuration["windows"] = config{
		Service: "my-service",
		App:     "My App.exe",
	}

	return configuration[runtime.GOOS]
}
