// +build windows

package platform

import (
	"bytes"
	"os/exec"
)

// KillProcessByName kills all processes by name provided
func killProcessByName(appname string) (bool, error) {
	cmd := exec.Command("taskkill", "/F", "/IM", appname)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, err
	}
	return true, nil
}

// naive startService using net cmd :D
func startService(service string) (bool, error) {
	cmd := exec.Command("net", "start", service)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, err
	}
	return true, nil
}

// naive stopService
func stopService(service string) (bool, error) {
	cmd := exec.Command("net", "Stop", service)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, err
	}
	return true, nil
}
