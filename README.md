a commandline gister in golang
---
[![GoDoc](https://godoc.org/github.com/viyatb/gister?status.svg)](https://godoc.org/github.com/viyatb/gister)
![](https://img.shields.io/github/issues/viyatb/gister.svg)


> This is a port of [gist](https://github.com/defunkt/gist) in Go

## Usage

1. Get the pre-built Linux (x86_64) built or download and build it yourself.

`gister` provides 3 optional CLI arguments.
  - `-p`: If `true`, the gist created will be public. Defaults to `true`.
  - `-d`: Provide a description. Defaults to `This is a gist`.
  - `-a`: If `true`, the gist created will be anonymous. Set `false` to create a gist for a user. Defaults to `true`.

2. Running `gister --h` for the available options and usage.


## LICENSE

[MIT](LICENSE.md)
