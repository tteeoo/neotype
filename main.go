package main

import (
	"flag"
	"fmt"
	"github.com/tteeoo/neotype/game"
	"github.com/tteeoo/neotype/util"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

var version = flag.Bool("version", false, "Print version information and exit")
var words = flag.Int("words", 0, "The number of words to test with")

func main() {

	flag.Parse()

	// Default config
	config := util.Config{
		Words: 20,
	}

	// Handle command-line arguments
	if *version {
		fmt.Println("NeoType version 0.1.0\n" +
			"Created by Theo Henson\n" +
			"Source available at https://github.com/tteeoo/neotype\n" +
			"Licensed under the Unlicense")
		os.Exit(1)
	}
	if *words != 0 {
		config.Words = *words
	}

	// Get data directory
	shareDir, err := util.ResolveShare()
	util.DieIf(err, "NeoType: error: cannot find data directory: %s\n", err)

	// Get terminal dimensions
	w, h, err := terminal.GetSize(0)
	util.DieIf(err, "NeoType: error: cannot get terminal size: %s\n", err)

	// Generate words
	dictionaryB, err := ioutil.ReadFile(shareDir + "/words.txt")
	util.DieIf(err, "NeoType: error: cannot read file '%s/words.txt': %s\n", shareDir, err)
	dictionary := strings.Split(string(dictionaryB), "\n")
	rand.Seed(time.Now().Unix())
	var chosen []string
	for i := 0; i < config.Words; i++ {
		chosen = append(chosen, dictionary[rand.Intn(len(dictionary))])
	}

	// Declare Game
	g := game.Game{
		WordString: strings.Join(chosen, " "),
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

			fmt.Printf("wpm: %d\n", g.WPM())
			fmt.Printf("acc: %d\n", g.Accuracy())
			fmt.Printf("raw: %d\n", g.Raw())

			err = exec.Command("stty", "-F", "/dev/tty", "echo").Run()
			util.DieIf(err, "NeoType: error: cannot run command 'stty -F /dev/tty echo': %s\n", err)

			os.Exit(0)
		}
	}()

	// Hide input and remove CR buffer
	err = exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	util.DieIf(err, "NeoType: error: cannot run command 'stty -F /dev/tty cbreak min 1': %s\n", err)
	err = exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	util.DieIf(err, "NeoType: error: cannot run command 'stty -F /dev/tty -echo': %s\n", err)

	// Start game
	err = g.Start()
	util.DieIf(err, "NeoType: error: %s\n", err)

	fmt.Println("")
	fmt.Print("\033[?1049l")
	fmt.Printf("wpm: %d\n", g.WPM())
	fmt.Printf("acc: %d\n", g.Accuracy())
	fmt.Printf("raw: %d\n", g.Raw())

	err = exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	util.DieIf(err, "NeoType: error: cannot run command 'stty -F /dev/tty echo': %s\n", err)

	os.Exit(0)
}
