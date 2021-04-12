package cmd

import (
	"bufio"
	"image/gif"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/markbates/pkger"
	"github.com/spf13/cobra"

	sdeck "github.com/dh1tw/streamdeck"
	_ "github.com/dh1tw/streamdeck-buttons" // needed to load stream-deck-button static assets
)

var gifCmd = &cobra.Command{
	Use:   "gif",
	Short: "show an animated gif on the streamdeck",
	Long:  `This example will show an animated gif on the streamdeck.`,
	Run:   animatedGif,
}

func init() {
	rootCmd.AddCommand(gifCmd)
}

func animatedGif(cmd *cobra.Command, args []string) {

	var mu sync.Mutex
	mu.Lock()
	mu.Unlock()

	sd, err := sdeck.NewStreamDeck()
	if err != nil {
		log.Panic(err)
	}
	defer sd.ClearAllBtns()

	_cat, err := pkger.Open("/assets/images/cat.gif")
	if err != nil {
		log.Panic(err)
	}
	defer _cat.Close()

	// resize, center and crop on the fly
	cat, err := gif.DecodeAll(bufio.NewReader(_cat))
	if err != nil {
		log.Panic(err)
	}

	_cat72, err := pkger.Open("/assets/images/cat-72px.gif")
	if err != nil {
		log.Panic(err)
	}
	defer _cat.Close()

	cat72, err := gif.DecodeAll(bufio.NewReader(_cat72))
	if err != nil {
		log.Panic(err)
	}

	// resize, center and crop on the fly
	_banana, err := pkger.Open("/assets/images/banana.gif")
	if err != nil {
		log.Panic(err)
	}
	defer _banana.Close()

	banana, err := gif.DecodeAll(bufio.NewReader(_banana))
	if err != nil {
		log.Panic(err)
	}

	// take in the raw; resize, center and crop on the fly
	_banana72, err := pkger.Open("/assets/images/banana-72px.gif")
	if err != nil {
		log.Panic(err)
	}
	defer _banana.Close()

	banana72, err := gif.DecodeAll(bufio.NewReader(_banana72))
	if err != nil {
		log.Panic(err)
	}

	if err := sd.FillGif(0, *cat); err != nil {
		log.Panic(err)
	}

	if err := sd.FillGif(4, *banana); err != nil {
		log.Panic(err)
	}

	if err := sd.FillGif(10, *cat72); err != nil {
		log.Panic(err)
	}

	if err := sd.FillGif(14, *banana72); err != nil {
		log.Panic(err)
	}

	// will cause system crash!!!!
	// for i := 0; i < 15; i++ {
	// 	if err := sd.FillGif(i, *cat); err != nil {
	// 		log.Panic(err)
	// 	}
	// }

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	for {
		select {
		case <-c:
			return
		}
	}
}
