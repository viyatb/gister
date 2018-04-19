a commandline gister in golang
---
[![GoDoc](https://godoc.org/github.com/viyatb/gister?status.svg)](https://godoc.org/github.com/viyatb/gister)
![](https://img.shields.io/github/issues/viyatb/gister.svg)


> This is a port of [gist](https://github.com/defunkt/gist) in Go

## Settings

1. [Create a personal access token] https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/)

2. Set the `GITHUB_TOKEN` environment variable to the value `username:token`,
   or write `username:token` to `~/.gist` file.

## Usage

1. Get the pre-built Linux (x86_64) built or download and build it yourself.

`gister` provides 3 optional CLI arguments.
  - `-public`: If `true`, the gist created will be public. Defaults to `true`.
  - `-d`: Provide a description. Defaults to `This is a gist`.
  - `-anonymous`: If `true`, the gist created will be anonymous. Set `false` to create a gist for a user. Defaults to `true`.

2. Run `gister --h` for the available options and usage.


## LICENSE

[MIT](LICENSE.md)
