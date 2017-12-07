package cmd

import (
	"testing"

	"github.com/flimzy/testy"
	"github.com/spf13/viper"
)

func TestIsValidFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Raw",
			input:    "raw",
			expected: true,
		},
		{
			name:     "JSON",
			input:    "json",
			expected: true,
		},
		{
			name:     "Foo",
			input:    "foo",
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := IsValidFormat(test.input)
			if test.expected != result {
				t.Errorf("Unexpected result: %v", result)
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name string
		conf *viper.Viper
		err  string
	}{
		{
			name: "defaults",
			conf: func() *viper.Viper {
				v := viper.New()
				v.SetDefault("format", "raw")
				return v
			}(),
			err: "",
		},
		{
			name: "format=json",
			conf: func() *viper.Viper {
				v := viper.New()
				v.SetDefault("format", "raw")
				v.Set("format", "json")
				return v
			}(),
			err: "",
		},
		{
			name: "format=foo",
			conf: func() *viper.Viper {
				v := viper.New()
				v.SetDefault("format", "raw")
				v.Set("format", "foo")
				return v
			}(),
			err: "Invalid output format 'foo'",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateConfig(test.conf)
			testy.Error(t, test.err, err)
		})
	}
}
