package main

import (
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

func main() {

	// Default config
	config := util.Config{
		Words: 12,
	}

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
	shareDir, err := util.ResolveShare()
	util.DieIf(err, "neotype: error: %s\n", err)

	// Get config file
	configFile, err := util.ResolveConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "neotype: warning: could not locate config file:", err)
	} else {
		// Read config
		err = config.Read(configFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "neotype: warning: could not read config file:", err)
		}
	}

	// Get terminal width
	w, _, err := terminal.GetSize(0)
	util.DieIf(err, "neotype: error: cannot get terminal width: %s\n", err)

	// Generate words
	dictionaryB, err := ioutil.ReadFile(shareDir + "/words.txt")
	util.DieIf(err, "neotype: error: cannot read file '%s/words.txt': %s\n", shareDir, err)
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
		Line:       1,
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
			util.DieIf(err, "neotype: error: cannot run command 'stty -F /dev/tty echo': %s\n", err)

			os.Exit(0)
		}
	}()

	// Hide input and remove CR buffer
	err = exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	util.DieIf(err, "neotype: error: cannot run command 'stty -F /dev/tty cbreak min 1': %s\n", err)
	err = exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	util.DieIf(err, "neotype: error: cannot run command 'stty -F /dev/tty -echo': %s\n", err)

	// Start game
	err = g.Start()
	util.DieIf(err, "neotype: error: %s\n", err)

	fmt.Println("")
	fmt.Print("\033[?1049l")
	fmt.Printf("wpm: %d\n", g.WPM())
	fmt.Printf("acc: %d\n", g.Accuracy())
	fmt.Printf("raw: %d\n", g.Raw())

	err = exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	util.DieIf(err, "neotype: error: cannot run command 'stty -F /dev/tty echo': %s\n", err)
}
