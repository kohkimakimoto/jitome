# Jitome

![logo.png](logo.png)

Jitome is a watcher for file changing

## Installation

```
$ go get github.com/kohkimakimoto/jitome
```

## Usage

Run `jitome init` to create `.jitome.yml` file and you should edit it.

The following is an example of configuration.

```yaml
# .jitome.yml
build:
    watch:   '.+\.go$'
    exclude: 'test\.go$'
    command: 'go build'

test:
    watch:   '.+\.go$'
    command: 'your test command'
```

The top level directives `build` and `test` are tasks that must be unique name in all of tasks.
`watch` is a regular expression string to define watching files.
`exclude` is a regular expression string to define excluding watching files.
`command` is a executed command  when it detects file changing.

## Author

Kohki Makimoto <kohki.makimoto@gmail.com>

## License

The MIT License (MIT)

## Inspired by

* [cespare/reflex](https://github.com/cespare/reflex)
* [romanoff/gow](https://github.com/romanoff/gow)
* [nathany/looper](https://github.com/nathany/looper)
* [mattn/goemon](https://github.com/mattn/goemon)
