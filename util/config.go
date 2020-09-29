package util

import (
	"encoding/json"
	"io/ioutil"
)

// Config represents the serialized config file.
type Config struct {
	// Words is the number of words to be typed.
	Words int `json:"words"`
}

// Read reads and sets the config.
func (c *Config) Read(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &c)
	if err != nil {
		return err
	}

	return nil
}
