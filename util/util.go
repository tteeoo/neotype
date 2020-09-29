package util

import (
	"errors"
	"fmt"
	"os"
	"os/user"
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
	share := os.Getenv("NEOTYPE_DATA")
	if share == "" {
		share = os.Getenv("XDG_DATA_HOME")
		if share == "" {
			user, err := user.Current()
			if err != nil {
				return "", errors.New("cannot get current user: " + err.Error())
			}
			share = user.HomeDir + "/.local/share/neotype"
		} else {
			share += "/neotype"
		}
	}
	_, err := os.Stat(share)
	if os.IsNotExist(err) {
		err = os.Mkdir(share, 0755)
		if err != nil {
			return "", errors.New("cannot create neotype data directory \"" + share + "\": " + err.Error())
		}
	} else if err != nil {
		return "", errors.New("cannot access neotype data directory \"" + share + "\": " + err.Error())
	}
	return share, nil
}

// ResolveConfig is the equivalent of ResolveShare but for the config file.
func ResolveConfig() (string, error) {
	config := os.Getenv("NEOTYPE_CONFIG")
	if config == "" {
		config = os.Getenv("XDG_CONFIG_HOME")
		if config == "" {
			user, err := user.Current()
			if err != nil {
				return "", errors.New("cannot get current user: " + err.Error())
			}
			config = user.HomeDir + "/.config/neotype/config.json"
		} else {
			config += "/neotype/config.json"
		}
	}
	_, err := os.Stat(config)
	if err != nil {
		return "", err
	}
	return config, nil
}
