package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func RunAtVscode() (bool, error) {

	exePath, err := os.Executable()
	if err != nil {
		return false, fmt.Errorf("os.Executable error " + err.Error())
	}

	contains := strings.Contains(exePath, "go-build")
	return contains, nil
}

func GetWorkPath() (string, error) {

	at_vscode, err := RunAtVscode()
	if err != nil {
		return "", err
	}

	if !at_vscode {
		exePath, err := os.Executable()
		if err != nil {
			return "", err
		}
		return filepath.Dir(exePath), nil
	} else {

		pwdPath, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return pwdPath, nil
	}

}
