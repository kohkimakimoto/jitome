# Jitome

![logo.png](logo.png)

Jitome is a simple file watcher.

## Installation

```
$ go get github.com/kohkimakimoto/jitome
```

## Usage

Run `jitome -i` to create `jitome.yml` file that is a main configuration file for jitome.
The following is an example of configuration.

```yaml
build:
  watch:
    - base: ""
      ignore_dir: [".git"]
      pattern: '.+\.go$'
  script: |
    go build .
```
run `jitome`. It is watching file changing.

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
* [mattn/goemon](https://github.com/mattn/goemon)
