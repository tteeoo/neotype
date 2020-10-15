package game

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"time"
)

// Game represents the current game state.
type Game struct {
	Started    bool
	StartTime  time.Time
	TotalTime  time.Duration
	WordString string
	Width      int
	Height     int

	Lines []string
	Page  int

	Line       int
	Index      int
	TotalTyped int
	TotalWrong int
}

// Start will start the game's main execution loop and print to the screen.
func (g *Game) Start() error {

	var scrollMode bool
	if g.Width*g.Height < len(g.WordString) {
		scrollMode = true
	}

	// Initialize screen
	fmt.Print("\033[?1049h")
	fmt.Print("\033[H\033[2J")
	g.Line = 1
	g.Page = 1

	if scrollMode {
		chars := len(g.WordString)
		var begin, end int
		for {
			end = begin + g.Width
			if end >= chars {
				g.Lines = append(g.Lines, g.WordString[chars-(chars%g.Width):])
				break
			}
			g.Lines = append(g.Lines, g.WordString[begin:end])
			begin = end
		}
		fmt.Print(g.NewPrintBuffer())
	} else {
		fmt.Print(g.WordString)
	}

	fmt.Printf("\033[%dA", int(len(g.WordString)/g.Width)+1)
	fmt.Printf("\033[%dD", g.Width)

	// Main loop
	g.TotalTime = time.Now().Sub(g.StartTime)
	var b []byte = make([]byte, 1)
	for {
		// Read character
		_, err := os.Stdin.Read(b)
		if err != nil {
			err2 := exec.Command("stty", "-F", "/dev/tty", "echo").Run()
			if err2 != nil {
				return fmt.Errorf("Cannot run command \"stty -F /dev/tty echo\": %s", err)
			}
			return fmt.Errorf("Cannot read from standard input")
		}
		g.TotalTyped++

		// Don't start timer if you haven't typed yet
		if !g.Started {
			g.StartTime = time.Now()
			g.Started = true
		} else {
			// Get current time
			g.TotalTime = time.Now().Sub(g.StartTime)
			// Check for finish
			if g.Index == len(g.WordString)-1 {
				break
			}
		}

		// Correct character
		if string(b) == string(g.WordString[g.Index]) {
			fmt.Print("\033[1C")
			g.Index++
		} else {
			// Incorrect
			g.TotalWrong++
		}

		// Line break
		if g.Index == g.Width*g.Line {
			g.Line++
			if scrollMode {
				if g.Page == 1 && g.Line == g.Page*g.Height+1 {
					g.Page++
					fmt.Print("\033[H\033[2J")
					fmt.Print(g.NewPrintBuffer())
					fmt.Printf("\033[%dA", int(len(g.WordString)/g.Width)+1)
					fmt.Printf("\033[%dD", g.Width)
				} else if g.Page != 1 && g.Line == g.Page*g.Height-1 {
					g.Page++
					fmt.Print("\033[H\033[2J")
					fmt.Print(g.NewPrintBuffer())
					fmt.Printf("\033[%dA", int(len(g.WordString)/g.Width)+1)
					fmt.Printf("\033[%dD", g.Width)
				} else {
					fmt.Print("\033[1B")
					fmt.Printf("\033[%dD", g.Width)
				}
			} else {
				fmt.Print("\033[1B")
				fmt.Printf("\033[%dD", g.Width)
			}
		}
	}

	return nil
}

// NewPrintBuffer returns the next string to be printed when in scroll mode
func (g *Game) NewPrintBuffer() string {
	var printBuffer string
	var slice []string
	if g.Line-1+g.Height > len(g.Lines) {
		slice = g.Lines[g.Line-1:]
	} else {
		slice = g.Lines[g.Line-1 : g.Line-1+g.Height]
	}
	for i := 0; i < len(slice); i++ {
		printBuffer += slice[i]
	}
	return printBuffer
}

// WPM calculates words per minute.
func (g *Game) WPM() int {
	return int(math.Round(float64(g.TotalTyped-g.TotalWrong) / 5.0 / g.TotalTime.Minutes()))
}

// Raw calculates the raw words per minute.
func (g *Game) Raw() int {
	return int(math.Round(float64(g.TotalTyped) / 5.0 / g.TotalTime.Minutes()))
}

// Accuracy calculates the accuracy.
func (g *Game) Accuracy() int {
	if g.TotalTyped != 0 {
		return int(math.Round(float64(g.TotalTyped-g.TotalWrong) / float64(g.TotalTyped) * 100.0))
	}
	return 100
}
