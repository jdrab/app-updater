// +build linux

package platform

import (
	"bytes"
	"os/exec"
)

// KillProcessByName kills all processes by name provided
func killProcessByName(appname string) (bool, error) {
	cmd := exec.Command("killall", "-9", appname)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, err
	}
	return true, nil
}

// startService using systemctl start serviceName
func startService(service string) (bool, error) {
	cmd := exec.Command("systemctl", "start", service)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, err
	}
	return true, nil
}

// stopService using systemctl stop serviceName
func stopService(service string) (bool, error) {
	cmd := exec.Command("systemctl", "stop", service)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, err
	}
	return true, nil

}
