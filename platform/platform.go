package platform

import (
	"log"
)

// Verbose if set to true package will print more information
var Verbose bool

// ServiceConfig - app and service config structure
type ServiceConfig struct {
	ServiceName string `json:"name"`
	AppName     string `json:"appName"`
	Version     string `json:"version"`
}

// KillProcessByName calls killProcessByName for required platform
func KillProcessByName(appname string) {
	ok, err := killProcessByName(appname)

	if err != nil {
		log.Printf("error killing process %v\n %v", appname, err)
	}

	if ok && Verbose {
		log.Printf("%v killed", appname)
	}

}

// StartService starts ServiceConfig.ServiceName
func StartService(service string) {
	ok, err := startService(service)
	if err != nil {
		log.Printf("error starting service %v\n %v", service, err)
	}

	if ok && Verbose {
		log.Printf("service %v started successfully", service)
	}

}

// StopService stops ServiceConfig.ServiceName
func StopService(service string) {
	ok, err := stopService(service)
	if err != nil {
		log.Printf("error stoping service %v\n %v", service, err)
	}
	if ok && Verbose {
		log.Printf("service %v stopped successfully", service)
	}

}
