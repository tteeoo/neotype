# NeoType

Tired of needing to go to a website powered by bloated JavaScript just to do a typing test?

Fear no more, for NeoType is here! It runs in your UNIX terminal and is powered by classic ANSI escape codes!

![preview](https://raw.githubusercontent.com/tteeoo/neotype/master/preview.gif)

It's currently quite minimal, partially by design, but I am also not sure what I should add. Open an issue (or better yet, a PR) if you have any feature requests!

## Installation

```
$ git clone https://github.com/tteeoo/neotype
$ cd neotype
$ mkdir ~/.local/share/neotype && cp share/words.txt ~/.local/share/neotype
$ go build
# cp neotype /usr/local/bin
```

Note: The above commands will compile NeoType, but a pre-compiled binary (Linux x86-64) is also provided on the latest GitHub release.

## Usage

Just run `neotype` to start.

Options:
* `-wordfile <string>`: The path to the file of newline-separated words to use. This can be an absolute or relative path, if it is invalid it will be treated as a relative path from the data directory (default "words.txt").
* `-words <int>`: The number of words to test with (default 20).

NeoType looks for the data directory at following paths, using the first valid path:

```
$NEOTYPE_DATA
$XDG_DATA_HOME/neotype
$HOME/.local/share/neotype
```

## License

The glorious Unlicense!

In other words, this software is dedicated to the public domain.
