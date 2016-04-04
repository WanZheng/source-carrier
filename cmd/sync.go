package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: sync,
}

func init() {
	clientCmd.AddCommand(syncCmd)
}

func sync(cmd *cobra.Command, args []string) {
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/sync", httpPort))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
	fmt.Println()
	if resp.StatusCode != 200 {
		fmt.Println("Failed: ", resp.Status)
		os.Exit(1)
	}
}
