package cmd

import (
	"embed"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

//go:embed assets
var assetDirectory embed.FS

// RootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "",
	Short: "This application contains several demos for the golang stream deck driver",
	Long: `This application contains several demos for the golang stream deck driver
which is being developed at https://github.com/dh1tw/streamdeck.
`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("device", "d", "", "Serial Number of stream deck to be used")
}
