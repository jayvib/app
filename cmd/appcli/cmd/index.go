/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/jayvib/app/config"
	"github.com/jayvib/app/pkg/elasticsearch"
	"github.com/olivere/elastic/v7"
	"github.com/prometheus/common/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var (
	type_     string
	env       string
	file      string
	indexName string
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Use to create an index in Elasticsearch",
	Long: `Use to create an index in Elasticsearch from the json file.

Example:
	./appcli elasticsearch create index --type=user --env=development --file=user.json --name=user
`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := elasticsearch.NewSimpleClient()
		validate()
		handleError(err)
		printFlags()
		mapping := readFile(file)
		res, err := client.CreateIndex(indexName).BodyString(string(mapping)).Do(context.Background())
		handleError(err)
		printResult(res)
	},
}

func init() {
	createCmd.AddCommand(indexCmd)
	indexCmd.Flags().StringVarP(&type_, "type", "t", "",
		"Use to indicate the app model type. Possible values are [user|author|article]")
	indexCmd.Flags().StringVarP(&env, "env", "e", "local",
		"Environment where index will be created. Possible values are [local|staging|prod]")
	indexCmd.Flags().StringVarP(&file, "file", "f", "",
		"Path to the file.")
	indexCmd.Flags().StringVarP(&indexName, "name", "n", "",
		"Index name.")

	err := indexCmd.MarkFlagRequired("type")
	handleError(err)
	err = indexCmd.MarkFlagRequired("env")
	handleError(err)
	err = indexCmd.MarkFlagRequired("file")
	handleError(err)
	err = indexCmd.MarkFlagRequired("name")
	handleError(err)
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func printFlags() {
	logrus.Println("Creating index type:", type_)
	logrus.Println("Environment:", env)
	logrus.Println("File:", file)
}

func printResult(r *elastic.IndicesCreateResult) {
	logrus.Println("Index:", r.Index)
	logrus.Println("Acknowledged:", r.Acknowledged)
	logrus.Println("Shard Acknowledged:", r.ShardsAcknowledged)
}

func readFile(fileName string) []byte {
	contents, err := ioutil.ReadFile(fileName)
	handleError(err)
	return contents
}

func message(env string) string {
	return fmt.Sprintf("Environment not match. Check the system environment value for '%s' expecting an '%s' value",
		config.AppEnvironmentKey, env)
}

func validate() {
	env_ := os.Getenv(config.AppEnvironmentKey)
	// Environment validation
	switch env {
	case "development":
		if env_ != config.DevelopmentEnv {
			log.Fatalf(message(config.DevelopmentEnv))
		}
	case "staging":
		if env_ != config.StagingEnv {
			log.Fatalf(message(config.StagingEnv))
		}
	case "prod":
		if env_ != config.ProductionEnv {
			log.Fatalf(message(config.ProductionEnv))
		}

	default:
		log.Fatalf("Environment: %s not yet supported", env)
	}
}
