package util

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

// DieIf printfs and exits if err != nil.
func DieIf(err error, format string, a ...interface{}) {
	if err != nil {
		fmt.Fprintf(os.Stderr, format, a...)
		os.Exit(1)
	}
}

// createDirectory will create a directory and set permissions,
// if the directory does not already exist. It returns true
// if a directory was created, and false if it was not.
func createDirectory(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err = os.Mkdir(path, 0755); err != nil {
			return false, err
		}
		return true, nil
	} else if err != nil {
		return false, err
	}
	return false, nil
}

func fileExists(filename string) (bool, error) {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return !info.IsDir(), nil
}

func getSharePath() (string, error) {
	neotypeData := os.Getenv("NEOTYPE_DATA")
	if neotypeData != "" {
		return neotypeData, nil
	}

	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome != "" {
		return filepath.Join(xdgDataHome, "/neotype"), nil
	}

	user, err := user.Current()
	if err == nil {
		return filepath.Join(user.HomeDir, "/.local/share/neotype"), nil
	}

	return "", err
}

// ResolveFilePath checks the working directory and the share
// for the given file, then returns the path if available.
func ResolveFilePath(filename string) (string, error) {
	matched, err := fileExists(filename)
	if err != nil {
		return "", err
	}
	if matched {
		return filename, nil
	}

	sharePath, err := getSharePath()
	if err != nil {
		return "", err
	}

	created, err := createDirectory(sharePath)
	if err != nil {
		return "", err
	}
	if !created {
		sharePathFile := path.Join(sharePath, filename)
		matched, err := fileExists(sharePathFile)
		if err != nil {
			return "", err
		}
		if matched {
			return sharePathFile, nil
		}
	}

	return "", errors.New("No searched directories contained the requested file")
}
