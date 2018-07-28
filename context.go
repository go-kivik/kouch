package kouch

import (
	"github.com/spf13/viper"

	"github.com/go-kivik/kouch/io"
	"github.com/go-kivik/kouch/log"
)

// Context is the command execution context.
type Context struct {
	Logger   log.Logger
	Conf     *viper.Viper
	Outputer io.OutputProcessor
}
