package cmd

import (
	"bufio"
	"fmt"
	"image"
	"log"
	"os"
	"os/signal"

	sdeck "github.com/dh1tw/streamdeck"
	"github.com/markbates/pkger"
	"github.com/spf13/cobra"
)

var iconsCmd = &cobra.Command{
	Use:   "icons",
	Short: "place a bunch of icons on the streamdeck",
	Long: `This example loads icons and places them on buttons in the first row
of the Stream Deck. The lightbulb icon on button 0 can be toggled.
	`,
	Run: icons,
}

func init() {
	rootCmd.AddCommand(iconsCmd)
}

func icons(cmd *cobra.Command, args []string) {
	sd, err := sdeck.NewStreamDeck()
	if err != nil {
		log.Panic(err)
	}
	defer sd.ClearAllBtns()

	_user, err := pkger.Open("/assets/images/user.png")
	if err != nil {
		log.Panic(err)
	}
	defer _user.Close()

	user, _, err := image.Decode(bufio.NewReader(_user))
	if err != nil {
		log.Panic(err)
	}

	_tux, err := pkger.Open("/assets/images/tux.png")
	if err != nil {
		log.Panic(err)
	}
	defer _tux.Close()

	tux, _, err := image.Decode(bufio.NewReader(_tux))
	if err != nil {
		log.Panic(err)
	}

	_warning, err := pkger.Open("/assets/images/warning.png")
	if err != nil {
		log.Panic(err)
	}
	defer _warning.Close()

	warning, _, err := image.Decode(bufio.NewReader(_warning))
	if err != nil {
		log.Panic(err)
	}

	_doctor, err := pkger.Open("/assets/images/doctor.png")
	if err != nil {
		log.Panic(err)
	}
	defer _doctor.Close()

	doctor, _, err := image.Decode(bufio.NewReader(_doctor))
	if err != nil {
		log.Panic(err)
	}

	_lightbulbOn, err := pkger.Open("/assets/images/lightbulb_on.png")
	if err != nil {
		log.Panic(err)
	}
	defer _lightbulbOn.Close()

	lightbulbOn, _, err := image.Decode(bufio.NewReader(_lightbulbOn))
	if err != nil {
		log.Panic(err)
	}

	_lightbulbOff, err := pkger.Open("/assets/images/lightbulb_off.png")
	if err != nil {
		log.Panic(err)
	}
	defer _lightbulbOff.Close()

	lightbulbOff, _, err := image.Decode(bufio.NewReader(_lightbulbOff))
	if err != nil {
		log.Panic(err)
	}

	if err := sd.FillImage(4, warning); err != nil {
		log.Panic(err)
	}
	if err := sd.FillImage(3, doctor); err != nil {
		log.Panic(err)
	}
	if err := sd.FillImage(2, tux); err != nil {
		log.Panic(err)
	}
	if err := sd.FillImage(1, user); err != nil {
		log.Panic(err)
	}
	if err := sd.FillImage(0, lightbulbOff); err != nil {
		log.Panic(err)
	}

	lightbulb := false

	onPressedCb := func(btnIndex int, state sdeck.BtnState) {
		fmt.Printf("Button: %d, %s\n", btnIndex, state)
		if btnIndex == 0 && state == sdeck.BtnPressed {
			if lightbulb {
				if err := sd.FillImage(0, lightbulbOff); err != nil {
					log.Panic(err)
				}
				lightbulb = false
			} else {
				if err := sd.FillImage(0, lightbulbOn); err != nil {
					log.Panic(err)
				}
				lightbulb = true
			}
		}
	}

	sd.SetBtnEventCb(onPressedCb)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
}
