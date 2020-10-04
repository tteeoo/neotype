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

// ResolveShare checks multiple environment variables, finds the location of the data directory,
// ensures it exists and it is accessable, then returns the path.
func ResolveShare() (string, error) {
	neotypeData := os.Getenv("NEOTYPE_DATA")
	xdgDataHome := os.Getenv("XDG_DATA_HOME")

	var share string
	if neotypeData != "" {
		share = neotypeData
	} else if xdgDataHome != "" {
		share = filepath.Join(xdgDataHome, "/neotype")
	} else {
		user, err := user.Current()
		if err != nil {
			return "", err
		}
		share = filepath.Join(user.HomeDir, "/.local/share/neotype")
	}

	_, err := os.Stat(share)
	if os.IsNotExist(err) {
		if err = os.Mkdir(share, 0755); err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}

	return share, nil
}
