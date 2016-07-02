# Mason
[![GoDoc](https://godoc.org/github.com/vdemeester/mason?status.png)](https://godoc.org/github.com/vdemeester/mason)
[![Build Status](https://travis-ci.org/vdemeester/mason.svg?branch=master)](https://travis-ci.org/vdemeester/mason)
[![Go Report Card](https://goreportcard.com/badge/github.com/vdemeester/mason)](https://goreportcard.com/report/github.com/vdemeester/mason)
[![License](https://img.shields.io/github/license/vdemeester/mason.svg)]()

Mason an helper to build client-driven docker container image
builder. *It is still very experimental*.

The goal of `mason` is to provide few helpers to ease the pain of
creating client-side docker image builder for those who find the
`Dockerfile` and `docker build` a little bit too limited.

It also holds a command (`mason`) with sub-command and simple, example
client-side builder. It's probably temporally as it's more examples
that actual ready-to-use binaries.

It uses [engine-api](https://github.com/docker/engine-api) and is
pretty tied to it (some structs of `engine-api` are popping up for now).

It is based on [libmason](https://github.com/vdemeester/libmason)
which provides helpers to create your own client-side buildler.

## TODO

- Features to support
    - [ ] Build cache mechanism (probably in `libmason` as it's mostly
      related to the `Step` builder)
- Builders
    - [ ] `dockerfile`: `Dockerfile` reference implementation
    - [ ] `dockramp`: `Dockerfile` *dockramp* divergence 👼
    - [ ] `to-be-named`: Enhanced `Dockerfile` that support more
      commands (like auto-commit, squashing, …)
