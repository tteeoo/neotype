package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"os/exec"
	"time"
	"strings"
	"math"
	"math/rand"
	"os/signal"
	"golang.org/x/crypto/ssh/terminal"
	"github.com/tteeoo/neotype/util"
)

func main() {

	// Handle command-line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-V", "--version":
			fmt.Println("neotype version 0.1.0\n" +
				"created by Theo Henson\n" +
				"source available at https://github.com/tteeoo/neotype\n" +
				"licensed under the Unlicense")
			return
		}
	}

	// Get data directory
	share, err := util.ResolveShare()
	util.DieIf(err, "neotype: error: %s\n", err)

	// Get terminal width
	w, _, err := terminal.GetSize(0)
	util.DieIf(err, "neotype: error: cannot get terminal width: %s\n", err)

	// Start
	started := false
	var startTime time.Time
	var totalTime time.Duration

	// Generate words
	dictionaryB, err := ioutil.ReadFile(share + "/words.txt")
	util.DieIf(err, "neotype: error: cannot read file '%s/words.txt': %s\n", share, err)
	dictionary := strings.Split(string(dictionaryB), "\n")
	rand.Seed(time.Now().Unix())
	var chosen []string
	for i := 0; i < 10; i++ {
		chosen = append(chosen, dictionary[rand.Intn(len(dictionary))])
	}

	// Level
	fmt.Print("\033[?1049h")
	fmt.Print("\033[H\033[2J")
	wordstring := strings.Join(chosen, " ")
	fmt.Println(wordstring)
	fmt.Printf("\033[%dA", int(len(wordstring) / w) + 1)
	fmt.Printf("\033[%dD", w)

	// Tracking vars
	line := 1
	var index int
	var totalTyped int
	var totalWrong int

	// Handle Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		for range c {
			totalTime = time.Now().Sub(startTime)
			fmt.Print("\033[H\033[2J")
			wpm := float64(totalTyped - totalWrong) / 5 / totalTime.Minutes()
			raw:= float64(totalTyped) / 5 / totalTime.Minutes()
			var acc float64
			if totalTyped != 0 {
				acc = float64(totalTyped - totalWrong) / float64(totalTyped)
			}
			fmt.Printf("wpm: %d\n", int(math.Round(wpm)))
			fmt.Printf("acc: %d\n", int(math.Round(acc * 100)))
			fmt.Printf("raw: %d\n", int(math.Round(raw)))

			err = exec.Command("stty", "-F", "/dev/tty", "echo").Run()
			util.DieIf(err, "neotype: error: cannot run command 'stty -F /dev/tty echo': %s\n", err)

			os.Exit(0)
		}
	}()

	// Hide input and remove CR buffer
	err = exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	util.DieIf(err, "neotype: error: cannot run command 'stty -F /dev/tty cbreak min 1': %s\n", err)
	err = exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	util.DieIf(err, "neotype: error: cannot run command 'stty -F /dev/tty -echo': %s\n", err)

	// Loop over num characters
	var b []byte = make([]byte, 1)
	for {
		_, err := os.Stdin.Read(b)
		if err != nil {
			err2 := exec.Command("stty", "-F", "/dev/tty", "echo").Run()
			util.DieIf(err2, "neotype: error: cannot run command 'stty -F /dev/tty echo': %s\n", err2)
			util.DieIf(err, "neotype: error: cannot read stdin")
		}
		totalTyped++

		// Don't start timer if you haven't typed yet
		if !started {
			startTime = time.Now()
			started = true
		}

		if string(b) == string(wordstring[index]) {
			fmt.Print("\033[1C")
			index++
		} else {
			totalWrong++
		}

		if index == len(wordstring) {
			totalTime = time.Now().Sub(startTime)
			break
		} else if index == w * line {
			line++
			fmt.Print("\033[1B")
			fmt.Printf("\033[%dD", w)
		}
	}

	fmt.Println("")
	wpm := float64(totalTyped - totalWrong) / 5 / totalTime.Minutes()
	raw:= float64(totalTyped) / 5 / totalTime.Minutes()
	var acc float64
	if totalTyped != 0 {
		acc = float64(totalTyped - totalWrong) / float64(totalTyped)
	}
	fmt.Print("\033[?1049l")
	fmt.Printf("wpm: %d\n", int(math.Round(wpm)))
	fmt.Printf("acc: %d\n", int(math.Round(acc * 100)))
	fmt.Printf("raw: %d\n", int(math.Round(raw)))

	err = exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	util.DieIf(err, "neotype: error: cannot run command 'stty -F /dev/tty echo': %s\n", err)
}
