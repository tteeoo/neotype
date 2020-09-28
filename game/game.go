package game

import (
	"fmt"
	"time"
	"os/exec"
	"os"
	"errors"
	"math"
)

// Game represents the current game state.
type Game struct {
	Started bool
	StartTime time.Time
	TotalTime time.Duration
	WordString string
	Width int

	Line int
	Index int
	TotalTyped int
	TotalWrong int
}

// Start the game.
func (g *Game) Start() error {

	// Initialize screen
	fmt.Print("\033[?1049h")
	fmt.Print("\033[H\033[2J")
	fmt.Println(g.WordString)
	fmt.Printf("\033[%dA", int(len(g.WordString) / g.Width) + 1)
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
				return errors.New("cannot run command 'stty -F /dev/tty echo': " + err.Error())
			}
			return errors.New("cannot read stdin")
		}
		g.TotalTyped++

		// Don't start timer if you haven't typed yet
		if !g.Started {
			g.StartTime = time.Now()
			g.Started = true
		} else  {
			// Get current time
			g.TotalTime = time.Now().Sub(g.StartTime)
			// Check for finish
			if g.Index == len(g.WordString) {
				break
			}
		}

		// Correct character
		if string(b) == string(g.WordString[g.Index]) {
			fmt.Print("\033[1C")
			g.Index++
		// Incorrect
		} else {
			g.TotalWrong++
		}

		// Line break
		if g.Index == g.Width * g.Line {
			g.Line++
			fmt.Print("\033[1B")
			fmt.Printf("\033[%dD", g.Width)
		}
	}

	return nil
}

// Calculate words per minute.
func (g *Game) WPM() int {
	return int(math.Round(float64(g.TotalTyped - g.TotalWrong) / 5.0 / g.TotalTime.Minutes()))
}

// Calculate the raw words per minute.
func (g *Game) Raw() int {
	return int(math.Round(float64(g.TotalTyped) / 5.0 / g.TotalTime.Minutes()))
}

// Calculate the accuracy.
func (g *Game) Accuracy() int {
	if g.TotalTyped != 0 {
		return int(math.Round(float64(g.TotalTyped - g.TotalWrong) / float64(g.TotalTyped) * 100.0))
	}
	return 100
}

