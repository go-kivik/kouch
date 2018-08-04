package root

import (
	"testing"

	"github.com/flimzy/testy"
	"github.com/spf13/viper"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name string
		conf *viper.Viper
		err  string
	}{
		{
			name: "defaults",
			conf: viper.New(),
			err:  "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateConfig(test.conf)
			testy.Error(t, test.err, err)
		})
	}
}
