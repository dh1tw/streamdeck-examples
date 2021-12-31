package cmd

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"os/signal"
	"strconv"

	sdeck "github.com/dh1tw/streamdeck"
	_ "github.com/dh1tw/streamdeck-buttons" // needed to load stream-deck-button static assets
	label "github.com/dh1tw/streamdeck-buttons/label"
	"github.com/dh1tw/streamdeck-buttons/ledbutton"
	"github.com/spf13/cobra"
)

var pagesCmd = &cobra.Command{
	Use:   "pages",
	Short: "Custom nested pages with navigation",
	Long:  `Custom nested pages with navigation`,
	Run:   pages,
}

func init() {
	rootCmd.AddCommand(pagesCmd)
}

func pages(cmd *cobra.Command, args []string) {
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

	p := NewStackPage(sd, nil)
	p.Draw()

	cb := func(keyIndex int, state sdeck.BtnState) {
		newPage := p.Set(keyIndex, state)
		if newPage != nil {
			p = newPage
			sd.ClearAllBtns()
			p.Draw()
		}
	}

	sd.SetBtnEventCb(cb)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

}

type stackPage struct {
	sd         *sdeck.StreamDeck
	parent     sdeck.Page
	btns       map[int]*ledbutton.LedButton
	rotators   map[int]*label.Label
	stackState map[int]bool
}

func NewStackPage(sd *sdeck.StreamDeck, parent sdeck.Page) sdeck.Page {

	sp := &stackPage{
		sd:         sd,
		parent:     parent,
		btns:       make(map[int]*ledbutton.LedButton),
		rotators:   make(map[int]*label.Label),
		stackState: make(map[int]bool),
	}

	smBtns := map[int]string{
		14: "NA",
		13: "EU",
		12: "OB11",
	}

	for pos, txt := range smBtns {
		led, err := ledbutton.NewLedButton(sd, pos,
			ledbutton.LedColor(ledbutton.LEDGreen),
			ledbutton.Text(txt))
		if err != nil {
			log.Panic(err)
		}
		sp.btns[pos] = led
		sp.stackState[pos] = false
	}

	rot := map[int]string{
		9: "300°",
		8: "042°",
		7: "148°",
	}

	for pos, txt := range rot {
		lbl, err := label.NewLabel(sd, pos, label.Text(txt))
		if err != nil {
			log.Panic(err)
		}
		sp.rotators[pos] = lbl
	}

	return sp
}

func (sp *stackPage) SetActive(active bool) {
}

func (sp *stackPage) Set(btnIndex int, state sdeck.BtnState) sdeck.Page {

	if state == sdeck.BtnReleased {
		return nil
	}

	btn, ok := sp.btns[btnIndex]
	if ok {
		if state == sdeck.BtnPressed {
			sp.stackState[btnIndex] = !sp.stackState[btnIndex]
			btn.SetState(sp.stackState[btnIndex])
			if err := btn.Draw(); err != nil {
				log.Println(err)
			}
			return nil
		}
	}
	_, ok = sp.rotators[btnIndex]
	if ok {
		return NewRotatorPage(sp.sd, sp)
	}
	return nil
}

func (sp *stackPage) Draw() {
	for _, btn := range sp.btns {
		if err := btn.Draw(); err != nil {
			log.Println(err)
		}
	}
	for _, rot := range sp.rotators {
		if err := rot.Draw(); err != nil {
			log.Println(err)
		}
	}
}

func (sp *stackPage) Parent() sdeck.Page {
	return sp.parent
}

type rotatorPage struct {
	sd            *sdeck.StreamDeck
	parent        sdeck.Page
	numPad        map[int]*label.Label
	newPos        *label.Label
	back          *label.Label
	set           *label.Label
	preset        *label.Label
	newPosText    string
	keyPadMapping map[int]int
}

func NewRotatorPage(sd *sdeck.StreamDeck, parent sdeck.Page) sdeck.Page {

	sp := &rotatorPage{
		sd:     sd,
		parent: parent,
		numPad: make(map[int]*label.Label),
		keyPadMapping: map[int]int{
			10: 0,
			3:  1,
			2:  2,
			1:  3,
			8:  4,
			7:  5,
			6:  6,
			13: 7,
			12: 8,
			11: 9,
		},
	}

	newPos, err := label.NewLabel(sd, 0,
		label.BgColor(color.RGBA{0, 255, 0, 255}),
		label.TextColor(color.RGBA{0, 0, 0, 255}))
	if err != nil {
		log.Panic(err)
	}
	sp.newPos = newPos

	for pos, num := range sp.keyPadMapping {
		l, err := label.NewLabel(sd, pos,
			label.Text(strconv.Itoa(num)))
		// label.BgColor(color.RGBA{255, 255, 0, 0}),
		// label.TextColor(color.RGBA{0, 0, 0, 255}),

		if err != nil {
			log.Panic(err)
		}
		sp.numPad[pos] = l
	}

	set, err := label.NewLabel(sd, 5, label.Text("SET"))
	if err != nil {
		log.Panic(err)
	}
	sp.set = set

	ret, err := label.NewLabel(sd, 4, label.Text("BACK"))
	if err != nil {
		log.Panic(err)
	}
	sp.back = ret

	preset, err := label.NewLabel(sd, 9, label.Text("PSET"))
	if err != nil {
		log.Panic(err)
	}
	sp.preset = preset

	return sp
}

func (sp *rotatorPage) SetActive(active bool) {
	return
}

func (sp *rotatorPage) Set(btnIndex int, state sdeck.BtnState) sdeck.Page {
	if state == sdeck.BtnReleased {
		return nil
	}

	switch btnIndex {
	case 4, 5:
		return sp.parent
	case 9:
		return NewPresetPage(sp.sd, sp)
	}

	_, ok := sp.numPad[btnIndex]
	if ok {
		if len(sp.newPosText) > 2 {
			return nil
		}
		num := sp.keyPadMapping[btnIndex]
		sp.newPosText = sp.newPosText + strconv.Itoa(num)
		sp.newPos.SetText(sp.newPosText)
		sp.Draw()
	}

	return nil
}

func (sp *rotatorPage) Draw() {
	for _, btn := range sp.numPad {
		btn.Draw()
	}
	sp.newPos.Draw()
	sp.preset.Draw()
	sp.back.Draw()
	sp.set.Draw()
}

func (sp *rotatorPage) Parent() sdeck.Page {
	return sp.parent
}

type presetPage struct {
	sd         *sdeck.StreamDeck
	parent     sdeck.Page
	btns       map[int]*label.Label
	btnMapping map[int]string
	back       *label.Label
}

func NewPresetPage(sd *sdeck.StreamDeck, parent sdeck.Page) sdeck.Page {

	pp := &presetPage{
		sd:     sd,
		parent: parent,
		btns:   make(map[int]*label.Label),
		btnMapping: map[int]string{
			3:  "NW",
			2:  "N",
			1:  "NE",
			8:  "W",
			6:  "E",
			13: "SW",
			12: "S",
			11: "SE",
			// 0:  "NA",
			// 5:  "EU",
			// 10: "VK",
		},
	}

	for pos, txt := range pp.btnMapping {
		l, err := label.NewLabel(sd, pos, label.Text(txt))
		if err != nil {
			log.Panic(err)
		}
		pp.btns[pos] = l
	}

	back, err := label.NewLabel(sd, 4, label.Text("BACK"))
	if err != nil {
		log.Panic(err)
	}
	pp.back = back

	return pp
}

func (sp *presetPage) SetActive(active bool) {
	return
}

func (pp *presetPage) Set(btnIndex int, state sdeck.BtnState) sdeck.Page {
	if state == sdeck.BtnReleased {
		return nil
	}

	switch btnIndex {
	case 4:
		return pp.parent
	}

	_, ok := pp.btns[btnIndex]
	if ok {
		return pp.parent.Parent()
	}

	return nil
}

func (pp *presetPage) Draw() {
	for _, btn := range pp.btns {
		btn.Draw()
	}
	pp.back.Draw()
}

func (pp *presetPage) Parent() sdeck.Page {
	return pp.parent
}
