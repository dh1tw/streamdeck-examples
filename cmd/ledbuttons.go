package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	sdeck "github.com/dh1tw/streamdeck"
	_ "github.com/dh1tw/streamdeck-buttons" // needed to load stream-deck-button static assets
	"github.com/dh1tw/streamdeck-buttons/ledbutton"
	"github.com/spf13/cobra"
)

var ledbtnsCmd = &cobra.Command{
	Use:   "ledbuttons",
	Short: "show a bunch of buttons with status LED",
	Long: `This example shows how to use the 'streamdeck/LedButton‘. It will
enumerate all the buttons on the panel with their ID and with a red, yellow or 
green LED which can be activated / deactivated with a button press.`,
	Run: ledbtns,
}

func init() {
	rootCmd.AddCommand(ledbtnsCmd)
}

func ledbtns(cmd *cobra.Command, args []string) {
	var sd *sdeck.StreamDeck
	var err error

	sdSerial := rootCmd.Flag("device").Value.String()

	if len(sdSerial) > 0 {
		sd, err = sdeck.NewStreamDeck(sdSerial)
	} else {
		sd, err = sdeck.NewStreamDeck()
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sd.ClearAllBtns()

	fmt.Println("using stream deck device with serial number", sd.Serial())

	defer sd.ClearAllBtns()

	btns := make(map[int]*ledbutton.LedButton)

	// Red Buttons
	for i := 0; i < 5; i++ {
		text := fmt.Sprintf("%03d", i)
		btn, err := ledbutton.NewLedButton(sd, i, ledbutton.Text(text), ledbutton.LedColor(ledbutton.LEDRed))
		if err != nil {
			fmt.Println(err)
		}
		btn.Draw()
		btns[i] = btn
	}

	// Yellow Buttons
	for i := 5; i < 10; i++ {
		text := fmt.Sprintf("%03d", i)
		btn, err := ledbutton.NewLedButton(sd, i, ledbutton.Text(text), ledbutton.LedColor(ledbutton.LEDYellow))
		if err != nil {
			fmt.Println(err)
		}
		btn.Draw()
		btns[i] = btn
	}

	// Green Buttons
	for i := 10; i < 15; i++ {
		text := fmt.Sprintf("%03d", i)
		btn, err := ledbutton.NewLedButton(sd, i, ledbutton.Text(text), ledbutton.LedColor(ledbutton.LEDGreen))
		if err != nil {
			fmt.Println(err)
		}
		btn.Draw()
		btns[i] = btn
	}

	btnChangedCb := func(btnIndex int, state sdeck.BtnState) {
		fmt.Printf("Button: %d, %s\n", btnIndex, state)
		if state == sdeck.BtnPressed {
			btn := btns[btnIndex]
			btn.SetState(!btn.State())
			if err := btn.Draw(); err != nil {
				log.Fatal(err)
			}
		}
	}
	sd.SetBtnEventCb(btnChangedCb)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
}
