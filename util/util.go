package util

import (
	"fmt"
	"os"
	"os/user"
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
	}
	if err != nil {
		return false, err
	}
	return false, nil
}

// fileExists returns true if the file exists and is not a directory.
func fileExists(filename string) (bool, error) {
	info, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
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
// for the given file, then returns the filepath if available.
func ResolveFilePath(file string) (string, error) {
	matched, err := fileExists(file)
	if err != nil {
		return "", err
	}
	if matched {
		return file, nil
	}

	share, err := getSharePath()
	if err != nil {
		return "", err
	}

	created, err := createDirectory(share)
	if err != nil {
		return "", err
	}
	if created {
		return "", fmt.Errorf("No file was found at: %s. A share was created at: %s", file, share)
	}

	shareFile := filepath.Join(share, file)
	matched, err = fileExists(shareFile)
	if err != nil {
		return "", err
	}
	if matched {
		return shareFile, nil
	}

	return "", fmt.Errorf("No file match at %s or %s", file, shareFile)
}
