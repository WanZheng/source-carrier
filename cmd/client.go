// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"syncfile/client"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// fs-watcherCmd represents the fs-watcher command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: watcher,
}

func init() {
	RootCmd.AddCommand(clientCmd)

	clientCmd.Flags().StringP("server", "s", "", "Server address")

	viper.BindPFlags(clientCmd.Flags())
}

func watcher(cmd *cobra.Command, args []string) {
	servAddr := viper.GetString("server")

	log.Print("root: ", root)
	log.Print("server: ", servAddr)
	log.Print("http_port: ", httpPort)

	if len(root) <= 0 || len(servAddr) <= 0 {
		fmt.Println(cmd.Flags().FlagUsages())
		os.Exit(1)
		return
	}

	if err := runClient(root, servAddr, httpPort); err != nil {
		log.Fatal("Failed to start sync client: ", err)
		os.Exit(1)
	}
}

func runClient(root, servAddr string, httpPort int) error {
	info, err := os.Stat(root)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return errors.New(fmt.Sprintf("'%s' is not a directory", root))
	}

	c := client.NewSyncClient(root, servAddr, httpPort)
	if err := c.Open(); err != nil {
		return err
	}

	return c.Run()
}
