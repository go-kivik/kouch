package kouch

import (
	"os"
	"path"
)

// Home returns the kouch home dir, or an empty string if the user has no
// home directory.
func Home() string {
	home := os.Getenv("HOME")
	if home == "" {
		return ""
	}
	return path.Join(home, HomeDir)
}
