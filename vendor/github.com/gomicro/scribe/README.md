# Scribe

[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/gomicro/scribe/build.yml?branch=main)](https://github.com/gomicro/scribe/actions?query=workflow%3ABuild)
[![Go Reportcard](https://goreportcard.com/badge/github.com/gomicro/scribe)](https://goreportcard.com/report/github.com/gomicro/scribe)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/gomicro/scribe)
[![License](https://img.shields.io/github/license/gomicro/scribe.svg)](https://github.com/gomicro/scribe/blob/master/LICENSE.md)
[![Release](https://img.shields.io/github/release/gomicro/scribe.svg)](https://github.com/gomicro/scribe/releases/latest)

Scribe is a Go library for structured, indented output in CLI tools. It formats output as a tree of named sections and leaf lines, with each nesting level indented by two spaces. All visual styling is delegated to a caller-supplied `Theme`, so the same structured output can be rendered plain or with ANSI color without changing the calling code.

# Usage

Add scribe to your module:

```sh
go get github.com/gomicro/scribe
```

## Basic usage

```go
s, err := scribe.NewScribe(os.Stdout, scribe.DefaultTheme())
if err != nil {
    log.Fatal(err)
}

s.BeginDescribe("Deployment")
{
    s.BeginDescribe("Services")
    s.EndDescribe()
    s.Print("Starting api")
    s.Print("Starting worker")
}
s.EndDescribe()
```

Output:

```
Deployment

  Services
    Starting api
    Starting worker
```

## Themed usage

Pass a custom `Theme` to apply ANSI color or any other decoration to describe headers and error lines:

```go
import "github.com/gomicro/scribe/color"

theme := &scribe.Theme{
    Describe: color.CyanFg,
    Print:    scribe.NoopDecorator,
    Error: func(err error) string {
        return color.RedFg("Error: " + err.Error())
    },
}

s, err := scribe.NewScribe(os.Stdout, theme)
if err != nil {
    log.Fatal(err)
}

s.BeginDescribe("Deployment")
{
    s.Print("Pushing image")
    s.Error(errors.New("registry unreachable"))
}
s.EndDescribe()
```

The `color` subpackage provides ANSI helpers for all standard and high-intensity foreground colors (`RedFg`, `CyanFg`, `HiBlueFg`, etc.).

# Versioning

The library will be versioned in accordance with [Semver 2.0.0](http://semver.org). See the [releases](https://github.com/gomicro/scribe/releases) section for the latest version. Until version 1.0.0 the library is considered to be unstable.

# License

See [LICENSE.md](./LICENSE.md) for more information.
