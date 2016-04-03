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
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
	"syncfile/data"
	"syncfile/server"

	"github.com/spf13/cobra"
)

var (
	serverPort int
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: serve,
}

func init() {
	RootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	serverCmd.Flags().StringVarP(&root, "root", "r", "", "Root path")
	serverCmd.Flags().IntVarP(&serverPort, "port", "p", 6001, "Local port")
}

func serve(cmd *cobra.Command, args []string) {
	log.Print("root: ", root)
	log.Print("port: ", serverPort)

	if len(root) <= 0 {
		fmt.Println(cmd.Flags().FlagUsages())
		os.Exit(1)
		return
	}
	if err := runServer(root, serverPort); err != nil {
		log.Fatal("Failed to run sync server: ", err)
		os.Exit(1)
	}
}

func runServer(root string, port int) error {
	info, err := os.Stat(root)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return errors.New(fmt.Sprintf("'%s' is not a directory", root))
	}

	// db
	log.Print("open db")
	dbPath := filepath.Join(root, "../server.sqlite3")
	db, err := data.OpenSqlDB(dbPath)
	if err != nil {
		return err
	}

	// rpc server
	s := server.NewSyncServer(root, db)
	rpc.RegisterName("Sync", s)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	log.Print("listen at ", port)
	return http.Serve(l, nil)
}
