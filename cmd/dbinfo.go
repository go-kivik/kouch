package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-kivik/couchdb/chttp"
	"github.com/go-kivik/kouch/log"
)

func init() {
	var dbname string
	registerCommand(func(log log.Logger, conf *viper.Viper) *cobra.Command {
		dbinfoCmd := &cobra.Command{
			Use:   "dbinfo",
			Short: "Returns meta data about the database",
			Long:  `Returns the result of GET /{db}`,
			Run:   dbinfo(log, conf),
		}
		dbinfoCmd.Flags().StringVarP(&dbname, "dbname", "d", "", "Database name")
		conf.BindPFlag("dbname", dbinfoCmd.Flags().Lookup("dbname"))

		return dbinfoCmd
	})

}

func dbinfo(log log.Logger, conf *viper.Viper) func(*cobra.Command, []string) {
	return func(_ *cobra.Command, _ []string) {
		response, err := getInfo(context.Background(), conf.GetString("server"), conf.GetString("dbname"))
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		log.Println(string(response))
		fmt.Printf("dbinfo called for %s / %s\n", conf.GetString("server"), conf.GetString("dbname"))
	}
}

func getInfo(ctx context.Context, server, dbname string) (json.RawMessage, error) {
	client, err := chttp.New(context.Background(), server)
	if err != nil {
		return nil, err
	}
	res, err := client.DoReq(context.Background(), http.MethodGet, dbname, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}
