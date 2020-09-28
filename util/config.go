package util

// Config represents the serialized config file.
type Config struct {
	// Words is the number of words to be typed.
	Words int `json:"default-num-words"`
}
