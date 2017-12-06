package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	var dbname string
	registerCommand(func(conf *viper.Viper) *cobra.Command {
		dbinfoCmd := &cobra.Command{
			Use:   "dbinfo",
			Short: "Returns meta data about the database",
			Long:  `Returns the result of GET /{db}`,
			Run:   dbinfo(conf),
		}
		dbinfoCmd.Flags().StringVarP(&dbname, "dbname", "d", "", "Database name")
		conf.BindPFlag("dbname", dbinfoCmd.Flags().Lookup("dbname"))

		return dbinfoCmd
	})

}

func dbinfo(conf *viper.Viper) func(*cobra.Command, []string) {
	return func(_ *cobra.Command, _ []string) {
		fmt.Printf("dbinfo called for %s / %s\n", conf.GetString("server"), conf.GetString("dbname"))
	}
}
