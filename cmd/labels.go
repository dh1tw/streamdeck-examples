package cmd

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cobra"

	sdeck "github.com/dh1tw/streamdeck"
	label "github.com/dh1tw/streamdeck-buttons/label"
)

var labelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "show a bunch of labeled icons on the streamdeck",
	Long: `This example will instantiate 15 labels on the streamdeck. Each Label
is setup as a counter which will increment every 100ms. If a button is
pressed it will be colored blue until it is released.`,
	Run: labels,
}

func init() {
	rootCmd.AddCommand(labelsCmd)
}

func labels(cmd *cobra.Command, args []string) {

	var mu sync.Mutex

	sd, err := sdeck.NewStreamDeck()
	if err != nil {
		log.Panic(err)
	}
	defer sd.ClearAllBtns()

	labels := make(map[int]*label.Label)

	for i := 0; i < 15; i++ {
		label, err := label.NewLabel(sd, i, label.Text(strconv.Itoa(i)))
		if err != nil {
			fmt.Println(err)
		}
		label.Draw()
		labels[i] = label
	}

	handleBtnEvents := func(btnIndex int, state sdeck.BtnState) {
		fmt.Printf("Button: %d, %s\n", btnIndex, state)
		mu.Lock()
		defer mu.Unlock()
		if state == sdeck.BtnPressed {
			col := color.RGBA{0, 0, 153, 0}
			labels[btnIndex].SetBgColor(image.NewUniform(col))
			labels[btnIndex].Draw()
		} else { // must be BtnReleased
			col := color.RGBA{0, 0, 0, 255}
			labels[btnIndex].SetBgColor(image.NewUniform(col))
			labels[btnIndex].Draw()
		}
	}

	sd.SetBtnEventCb(handleBtnEvents)

	ticker := time.NewTicker(time.Millisecond * 100)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	counter := 0

	for {
		select {
		case <-c:
			return
		case <-ticker.C:
			for i := 0; i < 15; i++ {
				mu.Lock()
				labels[i].SetText(fmt.Sprintf("%03d", counter))
				labels[i].Draw()
				mu.Unlock()
			}
			counter++
		}
	}
}
