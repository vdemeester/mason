# Mason

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

## Helpers

As previously said, `mason` provides some helpers to create
client-side builders, from the most low-level (almost `API` level) to
some higher level (with concept of Steps, commit/non-commit step,
etc…). Those *helpers* are designed to be composable.

### Base

### Step

## Binary

Mason currently holds a binary too, mostly to show off :P.
