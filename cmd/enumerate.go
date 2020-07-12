package cmd

import (
	"html/template"
	"log"
	"os"

	"github.com/dh1tw/hid"
	sd "github.com/dh1tw/streamdeck"
	"github.com/spf13/cobra"
)

var enumerateCmd = &cobra.Command{
	Use:   "enumerate",
	Short: "enumerate all connected Stream Decks",
	Long:  `This program enumerates all Corsair/Elgato Stream Decks connected to this computer.`,
	Run:   enumerate,
}

func init() {
	rootCmd.AddCommand(enumerateCmd)
}

var tmpl = template.Must(template.New("").Parse(
	`Found {{. | len}} Elgato Stream Deck(s): {{range .}}
	SerialNumber:        {{.Serial}}
	{{end}}`,
))

func enumerate(cmd *cobra.Command, args []string) {
	devices := hid.Enumerate(sd.VendorID, sd.ProductID)
	if err := tmpl.Execute(os.Stdout, devices); err != nil {
		log.Fatal(err)
	}
}
