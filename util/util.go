package util

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
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

// appendUnique will append any unique string to the slice.
// If the string is already in the slice, the slice is returned unchanged.
func appendUnique(slice []string, str string) []string {
	for _, s := range slice {
		if s == str {
			return slice
		}
	}
	return append(slice, str)
}

// GetSharePaths returns an array of all unique defined
// locations that may contain wordlist data.
func GetSharePaths() (shares []string) {
	workingDir, err := os.Getwd()
	if err == nil {
		shares = append(shares, workingDir)
	}

	neotypeData := os.Getenv("NEOTYPE_DATA")
	if neotypeData != "" {
		shares = appendUnique(shares, neotypeData)
	}

	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome != "" {
		shares = appendUnique(shares, filepath.Join(xdgDataHome, "/neotype"))
	}

	user, err := user.Current()
	if err == nil {
		shares = appendUnique(shares, filepath.Join(user.HomeDir, "/.local/share/neotype"))
	}

	return shares
}

// ResolveFilePath checks multiple paths for the given wordlist,
// then returns the path if available. It returns an error if more
// than one potential wordlist is found.
func ResolveFilePath(filename string, shares ...string) (string, error) {
	matchedFiles := []string{}

	for _, d := range shares {
		created, err := createDirectory(d)
		if err != nil {
			return "", err
		}
		if !created {
			filepath := path.Join(d, filename)
			matched, err := fileExists(filepath)
			if err != nil {
				return "", err
			}
			if matched {
				matchedFiles = append(matchedFiles, filepath)
			}
		}
	}

	if len(matchedFiles) == 1 {
		return matchedFiles[0], nil
	}

	if len(matchedFiles) > 1 {
		return "", errors.New("Matched multiple files: " + strings.Join(matchedFiles, ", "))
	}

	return "", errors.New("Wordfile was not found in any of these locations: " + strings.Join(shares, ", "))
}
