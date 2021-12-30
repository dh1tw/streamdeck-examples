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
	{{end}}
`,
))

func enumerate(cmd *cobra.Command, args []string) {
	streamdecks := deDuplicateHidDevices(hid.Enumerate(sd.VendorID, sd.ProductID))

	if err := tmpl.Execute(os.Stdout, streamdecks); err != nil {
		log.Fatal(err)
	}
}

// deDuplicateHidDevices removes duplicate HID devices based on the Serial Number.
// Sometimes a HID device can register itself for several Usages. This happens
// for example with Streamdecks on MacOS.
func deDuplicateHidDevices(deviceList []hid.DeviceInfo) []hid.DeviceInfo {

	devices := make(map[string]hid.DeviceInfo)

	for _, nextDevice := range deviceList {
		if _, dupe := devices[nextDevice.Serial]; !dupe {
			devices[nextDevice.Serial] = nextDevice
		}
	}

	deDupedDevicelist := make([]hid.DeviceInfo, 0, len(devices))

	for _, value := range devices {
		deDupedDevicelist = append(deDupedDevicelist, value)
	}

	return deDupedDevicelist
}
