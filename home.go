package kouch

import (
	"os"
	"path"
)

const (
	// kouchHome is the default directory where config is stored under the
	// users's home directory.
	kouchHome = ".kouch"
)

// Home returns the kouch home dir, or an empty string if the user has no
// home directory.
func Home() string {
	home := os.Getenv("HOME")
	if home == "" {
		return ""
	}
	return path.Join(home, kouchHome)
}
