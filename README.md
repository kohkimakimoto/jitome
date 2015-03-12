# Jitome

![logo.png](logo.png)

Jitome is a watcher for file changing

![jitome.gif](jitome.gif)

## Installation

```
$ go get github.com/kohkimakimoto/jijtome
```

## Usage

Run `jitome init` to create `.jitome.yml` file and you should edit it.

The following is an example of configuration.

```yaml
# .jitome.yml
build:
    watch:   "*.go"
    command: "go build"

test:
    watch:   "*.go"
    command: "your test command"
```

You can also use [TOML](https://github.com/toml-lang/toml) format to write configuration.

```toml
# .jitome.toml
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

To use `FILE` environment variable, you can get a file path that is detected changing

```yaml
display_modified_file:
    watch: "*"
    command: "cat $FILE"
```

## Author

Kohki Makimoto <kohki.makimoto@gmail.com>

## License

The MIT License (MIT)

## Inspired by

* [cespare/reflex](https://github.com/cespare/reflex)
* [romanoff/gow](https://github.com/romanoff/gow)
* [nathany/looper](https://github.com/nathany/looper)
