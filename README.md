# NeoType

Tired of needing to go to a website powered by bloated JavaScript just to do a typing test?

Fear no more, for NeoType is here! It runs in your UNIX terminal and is powered by classic ANSI escape codes!

## Installation

```
$ git clone https://github.com/tteeoo/neotype
$ cd neotype
$ install -d share ~/.local/share/neotype
$ go build
# cp neotype /usr/local/bin
```

Note: The above commands will compile NeoType, but a pre-compiled binary (Linux x86-64) is also provided on the latest GitHub release.

## Usage

Just run `neotype` to start, optionally provide a number of words: `neotype --words x`.

NeoType looks for the data directory at the following paths:

```
$NEOTYPE_DATA
$XDG_DATA_HOME/neotype
$HOME/.local/share/neotype
```

## License

The glorious Unlicense!

In other words, this software is dedicated to the public domain.
