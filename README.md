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

## Helpers & Builders

As previously said, `mason` provides some helpers to create
client-side builders, from the most low-level (almost `API` level) to
some higher level (with concept of Steps, commit/non-commit step,
etc…). Those *helpers* are designed to be composable.

### Base

The base Helper is located in the `base` package. It's a low level
interface (and implementation) of commands that might be needed for a
builder (get the image, create a container, commit a container to an
image, etc.).

```go
import (
    "github.com/vdemeester/mason/base"
    "github.com/docker/engine-api/types"
    "github.com/docker/engine-api/types/container"
)
// […]

helper := base.Newhelper(client)
// […]

image, err := helper.GetImage(context.Background(), "busybox", types.ImagePullOptions{})
// […]

resp, err := helper.ContainerCreate(context.Background(), types.ContainerCreateConfig{
    Config: &container.Config{
        Image: image.ID,
    }
}
// […]

imageID, err := helper.ContainerCommit(context.Background, resp.ID, types.ContainerCommitOptions{})
// […]
```

### Step

The `builder` package currently holds a `StepBuilder` which consists
of a composition of Step executed in order.

```go
import (
    "github.com/vdemeester/mason/builder"
)
// […]

steps := []Step{
    &MyStep{},
    // A step with that needs to create a container
    builder.WithDefaultCreate(&AnotherStep{}),
    // A step that will commit the container
    builder.WithCommit(&AThirdStep{}),
    // Or remove the container
    builder.WithRemove(&AFourthStep{}),
    // Or all of them ?
    builder.WithCreate(build.WithCommitAndRemove(&MyStep{})),
}

builder := builder.WithSteps(builder.DefaultBuilder(client))
image, err := builder.Run()
// […]
```

See the [godoc](https://godoc.org/github.com/vdemeester/mason) on how to create steps.

## Binary

Mason currently holds a binary too, mostly to show off :P.
