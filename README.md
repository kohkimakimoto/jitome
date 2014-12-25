# Jitome

![logo.png](logo.png)

Jitome watches files and runs a command when files change.

![jitome.gif](jitome.gif)

## Installation

```
$ go get github.com/kohkimakimoto/jijtome
```

## Usage

Create `.jitome` configuration file. The following is an example. The file format is [TOML](https://github.com/toml-lang/toml).

```toml
[build]
watch=["*.go"]
command="go build"

[test]
watch=["*.go"]
command="your test command"
```

Run `jitome`

```
$ jitome
```

## Author

Kohki Makimoto <kohki.makimoto@gmail.com>

## License

The MIT License (MIT)

## Inspired by

* [cespare/reflex](https://github.com/cespare/reflex)
* [romanoff/gow](https://github.com/romanoff/gow)
* [nathany/looper](https://github.com/nathany/looper)
