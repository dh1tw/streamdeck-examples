package cmd

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	sdeck "github.com/dh1tw/streamdeck"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/spf13/cobra"
)

var textbtnsCmd = &cobra.Command{
	Use:   "textbuttons",
	Short: "just buttons with custom text",
	Long:  `just buttons with custom text`,
	Run:   textbtns,
}

func init() {
	rootCmd.AddCommand(textbtnsCmd)
}

var monoFont *truetype.Font

func textbtns(cmd *cobra.Command, args []string) {

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

	// Load the font
	_monoFont, err := assetDirectory.Open("assets/fonts/mplus-1m-regular.ttf")
	if err != nil {
		log.Fatal(err)
	}
	defer _monoFont.Close()

	data, err := ioutil.ReadAll(_monoFont)
	if err != nil {
		log.Panic(err)
	}

	monoFont, err = freetype.ParseFont(data)
	if err != nil {
		log.Panic(err)
	}

	lineLabel := sdeck.TextLine{
		Font:      monoFont,
		FontColor: color.RGBA{255, 255, 0, 0}, // Yellow
		FontSize:  22,
		PosX:      10,
		PosY:      5,
		Text:      "STATE",
	}

	linePressed := sdeck.TextLine{
		Font:      monoFont,
		FontColor: color.RGBA{255, 255, 255, 0}, // White
		FontSize:  14,
		PosX:      12,
		PosY:      30,
		Text:      "PRESSED",
	}

	lineLongPressed := sdeck.TextLine{
		Font:      monoFont,
		FontColor: color.RGBA{255, 255, 0, 0}, // Yellow
		FontSize:  14,
		PosX:      12,
		PosY:      30,
		Text:      "LPRESSED",
	}

	lineReleased := sdeck.TextLine{
		Font:      monoFont,
		FontColor: color.RGBA{255, 0, 0, 0}, // Red
		FontSize:  14,
		PosX:      9,
		PosY:      30,
		Text:      "RELEASED",
	}

	pressedText := sdeck.TextButton{
		BgColor: color.RGBA{0, 0, 0, 0},
		Lines:   []sdeck.TextLine{lineLabel, linePressed},
	}

	longPressedText := sdeck.TextButton{
		BgColor: color.RGBA{0, 0, 0, 0},
		Lines:   []sdeck.TextLine{lineLabel, lineLongPressed},
	}

	releasedText := sdeck.TextButton{
		BgColor: color.RGBA{0, 0, 0, 0},
		Lines:   []sdeck.TextLine{lineLabel, lineReleased},
	}

	for i := 0; i < 15; i++ {
		sd.WriteText(i, releasedText)
	}

	btnEvtCb := func(btnIndex int, state sdeck.BtnState) {
		fmt.Printf("Button: %d, %s\n", btnIndex, state)
		switch state {
		case sdeck.BtnPressed:
			sd.WriteText(btnIndex, pressedText)
		case sdeck.BtnLongPressed:
			sd.WriteText(btnIndex, longPressedText)
		case sdeck.BtnReleased:
			sd.WriteText(btnIndex, releasedText)
		}
	}

	sd.SetBtnEventCb(btnEvtCb)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
}
