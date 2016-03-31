# Jitome

[![Build Status](https://travis-ci.org/kohkimakimoto/jitome.svg?branch=master)](https://travis-ci.org/kohkimakimoto/jitome)

![logo.png](logo.png)

Jitome is a simple file watcher.

* [Requirement](#requirement)
* [Installation](#installation)
* [Usage](#usage)
    * [Configuration](#configuration)
    * [Target](#target)
    * [Init target](#init-target)
* [Author](#author)
* [License](#license)
* [Inspired by](#inspired-by)

## Requirement

Go1.6 or later.

## Installation

```
$ go get github.com/kohkimakimoto/jitome
```

## Usage

Run `jitome -i` to create `jitome.yml` file that is a main configuration file for Jitome.

```yaml
build:
  watch:
    - base: ""
      ignore: [".git"]
      pattern: "*.go"
  script: |
    go build .
```

Run `jitome`. It is watching file changing.

```
$ jitome
```

When you change a `.go` file, Jitome detects it and runs a target.

```
[jitome] starting jitome...
[jitome] loading config 'jitome.yml'
[jitome] evaluating target 'build'.
[jitome] watching files...
[jitome] 'build' target detected 'color.go' changing [write]. running script.
[jitome] 'build' target finished script.
[jitome] 'build' target detected 'watcher.go' changing [write]. running script.
[jitome] 'build' target finished script.
[jitome] 'build' target detected 'watcher.go' changing [write]. running script.
[jitome] 'build' target finished script.
```

## Configuration

Default configuration file that Jitome uses is `jitome.yml` at the current directory. you can change it by using `-c` option.

The following is an example of the configuration.

```yaml
build:
  watch:
    - base: ""
      ignore: [".git"]
      pattern: "*.go"
  script: |
    go build .

test:
  watch:
    - base: ""
      ignore: [".git"]
      pattern: "*.go"
  script: |
    go test .
```

### Target

The top level property as the above `build` and `test` is a ***target***.

***target*** is a unit of config that defines watching patterns and a script that runs when it detects file changing.

### Init target

The `init` target is a special purpose target.
It has only `script` property and runs when Jitome starts.

```yaml
init:
  script: |
    echo "booted!"

```

## Author

Kohki Makimoto <kohki.makimoto@gmail.com>

## License

The MIT License (MIT)

## Inspired by

* [cespare/reflex](https://github.com/cespare/reflex)
* [romanoff/gow](https://github.com/romanoff/gow)
* [nathany/looper](https://github.com/nathany/looper)
* [mattn/goemon](https://github.com/mattn/goemon)
