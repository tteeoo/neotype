package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/tteeoo/neotype/game"
	"github.com/tteeoo/neotype/util"
	"golang.org/x/crypto/ssh/terminal"
)

var version = flag.Bool("version", false, "Print version information and exit.")
var words = flag.Int("words", 20, "The number of words to test with.")
var wordFile = flag.String("wordfile", "words.txt", "The path to the file of newline-separated words to use.\nThis can be an absolute or relative path.\nIf it is invalid it will be treated as a relative path from the data directory.")
var textFile = flag.String("textfile", "", "The path to a file containing exactly what to type.\nThe path works the same as wordfile. (overrides wordfile)")

func main() {

	flag.Parse()

	// Handle command-line arguments
	if *version {
		fmt.Println("NeoType version 0.1.2\n" +
			"Created by Theo Henson\n" +
			"Source available at https://github.com/tteeoo/neotype\n" +
			"Licensed under the Unlicense")
		os.Exit(1)
	}

	// Get terminal dimensions
	w, h, err := terminal.GetSize(0)
	util.DieIf(err, "NeoType: Error: Cannot get terminal size: %s\n", err)

	// Get word source
	var wordString string

	if *textFile == "" {
		// Get word file
		wordFilePath, err := util.ResolveFilePath(*wordFile)
		util.DieIf(err, "NeoType: Error: Cannot find word file: %s\n", err)

		// Generate words
		dictionaryB, err := ioutil.ReadFile(wordFilePath)
		util.DieIf(err, "NeoType: Error: Cannot read file \"%s\": %s\n", wordFilePath, err)
		dictionary := strings.Split(string(dictionaryB), "\n")
		var fixedDictionary []string
		for _, v := range dictionary {
			if v != "" {
				fixedDictionary = append(fixedDictionary, v)
			}
		}
		rand.Seed(time.Now().Unix())
		var chosen []string
		for i := 0; i < *words; i++ {
			chosen = append(chosen, fixedDictionary[rand.Intn(len(fixedDictionary))])
		}
		wordString = strings.Join(chosen, " ")
	} else {
		textFilePath, err := util.ResolveFilePath(*textFile)
		util.DieIf(err, "NeoType: Error: Cannot find text file: %s\n", err)

		b, err := ioutil.ReadFile(textFilePath)
		util.DieIf(err, "NeoType: Error: Cannot read file: %s\n", err)

		// TODO: newline support
		wordString = strings.ReplaceAll(string(b), "\n", " ")
	}

	// Declare Game
	g := game.Game{
		WordString: strings.ReplaceAll(wordString, "\t", "    "),
		Width:      w,
		Height:     h,
	}

	// Handle Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Print("\033[H\033[2J")
			fmt.Print("\033[?1049l")

			fmt.Printf("WPM: %d\n", g.WPM())
			fmt.Printf("Acc: %d\n", g.Accuracy())
			fmt.Printf("Raw: %d\n", g.Raw())

			err = exec.Command("stty", "-F", "/dev/tty", "echo").Run()
			util.DieIf(err, "NeoType: Error: Cannot run command \"stty -F /dev/tty echo\": %s\n", err)

			os.Exit(0)
		}
	}()

	// Hide input and remove CR buffer
	err = exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	util.DieIf(err, "NeoType: Error: Cannot run command \"stty -F /dev/tty cbreak min 1\": %s\n", err)
	err = exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	util.DieIf(err, "NeoType: Error: Cannot run command \"stty -F /dev/tty -echo\": %s\n", err)

	// Start game
	err = g.Start()
	util.DieIf(err, "NeoType: Error: Cannot start the game: %s\n", err)

	fmt.Println("")
	fmt.Print("\033[?1049l")
	fmt.Printf("WPM: %d\n", g.WPM())
	fmt.Printf("Acc: %d\n", g.Accuracy())
	fmt.Printf("Raw: %d\n", g.Raw())

	err = exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	util.DieIf(err, "NeoType: Error: Cannot run command \"stty -F /dev/tty echo\": %s\n", err)

	os.Exit(0)
}
