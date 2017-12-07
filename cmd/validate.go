package cmd

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	// FormatRaw outputs raw format, as it was received from the server.
	FormatRaw = "raw"
	// FormatJSON formats the output into human-readable JSON, with indentation.
	FormatJSON = "json"
)

// IsValidFormat returns true if the specified output format is valid.
func IsValidFormat(format string) bool {
	return format == FormatRaw || format == FormatJSON
}

// ValidateConfig validates the configuration, returning an error if something
// is wrong.
func ValidateConfig(conf *viper.Viper) error {
	if !IsValidFormat(conf.GetString("format")) {
		return fmt.Errorf("Invalid output format '%s'", conf.GetString("format"))
	}
	return nil
}
