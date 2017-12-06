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
			Run:   dbinfo,
		}
		dbinfoCmd.Flags().StringVarP(&dbname, "dbname", "d", "", "Database name")
		conf.BindPFlag("dbname", dbinfoCmd.PersistentFlags().Lookup("dbname"))

		return dbinfoCmd
	})

}

func dbinfo(cmd *cobra.Command, args []string) {
	fmt.Printf("dbinfo called for %s / %s\n", viper.GetString("server"), viper.GetString("dbname"))
}
