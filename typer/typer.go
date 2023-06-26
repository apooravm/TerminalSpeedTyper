package typer

import (
	"fmt"
	"time"
	// "os"
	// "bufio"
	"github.com/gdamore/tcell/v2"
	// "github.com/eiannone/keyboard"
	// "github.com/gookit/color"
)

type Typer struct {
	x int
	y int
	passageText string
	inputQueue []rune
	displayQueue []rune
	screen tcell.Screen
	style tcell.Style
	styleCorrect tcell.Style
	styleIncorrect tcell.Style
	WIN_SIZE int
	head int
	temp int
	startTime time.Time
	timeFlag bool
	timerStop chan struct{}
}

func (this *Typer) init() {
	this.head = 0
	if len(this.passageText) < 40 {
		this.WIN_SIZE = len(this.passageText)
	} else {
		this.WIN_SIZE = 40
	}
	this.timeFlag = false
	this.refresh()

	this.show()
}

func (this *Typer) refresh() {
	if this.head + this.WIN_SIZE >= len(this.passageText) {
		this.displayQueue = []rune(this.passageText[len(this.passageText)-this.WIN_SIZE:])
	} else {
		this.displayQueue = []rune(this.passageText[this.head:this.head + this.WIN_SIZE])
	}
}

func (this *Typer) showTimeRoutine() {
	stop := make(chan struct{})
	this.timerStop = stop
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				clearLineHoriz(this.screen, this.y+3)
				duration := time.Since(this.startTime).String()
				var comb []rune
				x := 10

				for _, char := range duration {
					charString := string(char) // Convert the character to a string
					this.screen.SetContent(x, this.y+3, []rune(charString)[0], comb, this.style)
					x++
				}

				this.screen.Show()
				time.Sleep(time.Second) // Wait for 1 second before the next update
			}
		}
	}()

}

func (this *Typer) showTime() {
	clearLineHoriz(this.screen, this.y+3)
	duration := time.Since(this.startTime).String()
	var comb []rune
	x := this.x

	for _, char := range duration {
		charString := string(char) // Convert the character to a string
		this.screen.SetContent(x, this.y+3, []rune(charString)[0], comb, this.style)
		x++
	}

	this.screen.Show()
}

func (this *Typer) show() {
	clearLineHoriz(this.screen, this.y)
	var comb []rune
	// for idx, char := range this.displayQueue {
	// 	this.screen.SetContent(this.x + idx - this.WIN_SIZE/2, this.y-1, char, comb, this.style)
	// }

	for idx, char := range this.displayQueue {
		this.screen.SetContent(this.x + idx - this.WIN_SIZE/2, this.y, char, comb, this.style)
	}

	// Only consider the last this.WIN_SIZE/2 elements of this.inputQueue
	currInput := this.inputQueue
	if this.head + this.WIN_SIZE >= len(this.passageText) {
		currInput = this.inputQueue[len(this.inputQueue)-(this.WIN_SIZE/2) - this.temp:]
		this.temp += 1
	} else {
		if len(this.inputQueue) > this.WIN_SIZE/2 {
			currInput = this.inputQueue[len(this.inputQueue)-(this.WIN_SIZE/2):]
		}
	}

	for idx, char := range currInput {
		value := rune(this.displayQueue[idx])
		if char == value {
			this.screen.SetContent(this.x + idx - this.WIN_SIZE/2, this.y, value, comb, this.styleCorrect)
		} else {
			this.screen.SetContent(this.x + idx - this.WIN_SIZE/2, this.y, value, comb, this.styleIncorrect)
		}
	}
}

func (this *Typer) appendText(char rune) {
	if !this.timeFlag {
		this.startTime = time.Now()
		this.timeFlag = true
	}
	this.showTime()
	if len(this.inputQueue) >= len(this.passageText) {
		// close(this.timerStop)
		return
	}
	this.inputQueue = append(this.inputQueue, char)
	if len(this.inputQueue) > this.WIN_SIZE/2 {
		this.head += 1
	}
	this.refresh()

	this.show()
}

func (this *Typer) popChar() {
	if len(this.inputQueue) >= len(this.passageText) {
		return
	}
	if len(this.inputQueue) == 0 {
		return
	}

	if this.head != 0 {
		this.head -= 1
	}

	this.inputQueue = this.inputQueue[:len(this.inputQueue)-1]
	this.refresh()
	this.show()
}

func (this *Typer) getScore() {
	
}

func Start(passageText string) {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Println("Error Creating Screen!", err)
		return
	}

	if err := screen.Init(); err != nil {
		fmt.Println("Error initializing", err)
		return
	}

	defer screen.Fini()

	screen.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite))

	screen.Clear()

	width, height := screen.Size()
	emptyFunc(width, height)

	BACK_STYLED := tcell.StyleDefault.
					Foreground(tcell.ColorBlack).
					Background(tcell.ColorWhite)

	GREY_STYLED := tcell.StyleDefault.
					Foreground(tcell.ColorLightGray).
					Background(tcell.ColorWhite)

	RED_STYLED := tcell.StyleDefault.
					Foreground(tcell.ColorBlack).
					Background(tcell.ColorRed)

	GREEN_STYLED := tcell.StyleDefault.
					Foreground(tcell.ColorBlack).
					Background(tcell.ColorGreen)

	addStrHoriz(0, height-1, "quit: esc", BACK_STYLED, screen)

	t1 := Typer{
		passageText: passageText,
		x: width/2,
		y: height/2-4,
		screen: screen,
		style: GREY_STYLED,
		styleCorrect: GREEN_STYLED,
		styleIncorrect: RED_STYLED,
	}
	t1.init()

	// Main Event Loop
	for {
		event := screen.PollEvent()
		switch event := event.(type) {
			case *tcell.EventKey:
				if event.Key() == tcell.KeyEscape {
					return
				} else if event.Key() == tcell.KeyLeft {
					continue

				} else if event.Key() == tcell.KeyRight {
					continue

				} else if event.Key() == tcell.KeyUp {
					continue

				} else if event.Key() == tcell.KeyDown {
					continue

				} else if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
					t1.popChar()

				} else {
					t1.appendText(event.Rune())
				}
		}
		// t1.showTime()
		screen.Show()
	}
}

func emptyFunc(x, y int) {
	return
}

func addStrHoriz(x, y int, str string, style tcell.Style, screen tcell.Screen) {
	clearLineHoriz(screen, y)
	var comb []rune
	for ind, char := range str {
		screen.SetContent(x+ind, y, char, comb, style)
	}
}

func addStrVert(x, y int, str string, style tcell.Style, screen tcell.Screen) {
	clearLineVert(screen, x)
	var comb []rune
	for ind, char := range str {
		screen.SetContent(x, y+ind, char, comb, style)
	}
	screen.ShowCursor(0, 0)
	screen.Show()
}

func clearLineHoriz(screen tcell.Screen, y int) {
	width, _ := screen.Size()

	for x := 0; x < width; x++ {
		screen.SetContent(x, y, ' ', nil, tcell.StyleDefault)
	}
}

func clearLineVert(screen tcell.Screen, x int) {
	_, height := screen.Size()

	for y := 0; y < height; y++ {
		screen.SetContent(x, y, ' ', nil, tcell.StyleDefault)
	}
}