# Jitome

![logo.png](logo.png)

Jitome is a watcher for file changing.

## Installation

```
$ go get github.com/kohkimakimoto/jitome
```

## Usage

Run `jitome -i` to create `jitome.yml` file and you should edit it.

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

run `jitome`

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
