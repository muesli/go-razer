# go-razer

Control your Razer (Chroma) devices from Go

## Installation

Make sure you have a working Go environment (Go 1.7 or higher is required).
See the [install instructions](http://golang.org/doc/install.html).

To install go-razer, simply run:

    go get github.com/muesli/go-razer

## Lantern CLI

Lantern is a CLI tool (using go-razer) to control your Razer devices. Run the
following command to install it:

    go get github.com/muesli/go-razer/cmd/lantern

### Hardware Effects

Available effects are `wave`, `reactive`, `spectrum`, `breath`, `breathdual`, `breathrandom`, `starlight`, `starlightdual`, `starlightrandom`, `ripple` and `ripplerandom`.

```
$ lantern -brightness 90 -effect wave
$ lantern -brightness 90 -effect starlightdual -color "#00ff00" -secondary "#aa00aa"
```

### Plain Background Color

```
$ lantern -color "#6b6b00"
```

### top Mode

Monitor your system's CPU usage by turning your keyboard into a gauge:

```
$ lantern -brightness 100 -top
```

### Themes

There are currently only two available themes (feel free to submit more!) named
`happy` and `rainbow`. To try them out, run:

```
$ lantern -brightness 100 -theme happy
```

## Development

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/muesli/go-razer)
[![Build Status](https://travis-ci.org/muesli/go-razer.svg?branch=master)](https://travis-ci.org/muesli/go-razer)
[![Go ReportCard](http://goreportcard.com/badge/muesli/go-razer)](http://goreportcard.com/report/muesli/go-razer)
