// Package dockerfile is a implementation client-side of the Dockerfile builds
// supported by docker build (daemon-side)
package dockerfile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	// "github.com/docker/engine-api/types/strslice"
	"github.com/vdemeester/libmason"
	"github.com/vdemeester/libmason/builder"
	"github.com/vdemeester/mason/dockerfile/parser"
)

// DefaultDockerfile holds the default name for a Dockerfile
const DefaultDockerfile = "Dockerfile"

type handlerFunc func(args []string, heredoc string) error

// Builder holds attributes and defines method to build Dockerfile
type Builder struct {
	helper           libmason.Helper
	contextDirectory string
	dockerfilePath   string
	references       []string
	out              io.Writer
}

// NewBuilder creates a new Dockerfile Build with the specified arguments.
func NewBuilder(c client.APIClient, contextDirectory, dockerfilePath string, tags []string) (*Builder, error) {
	// Validate that the context is a directory
	if err := validateContextDirectory(contextDirectory); err != nil {
		return nil, fmt.Errorf("unable to access build context directory: %s", err)
	}
	// Validate that dockerfilePath exists and is valid
	if dockerfilePath == "" {
		dockerfilePath = filepath.Join(contextDirectory, DefaultDockerfile)
	}
	if err := validateDockerfilePath(dockerfilePath); err != nil {
		return nil, fmt.Errorf("unable to access build file: %s", err)
	}
	if err := validateReferences(tags); err != nil {
		return nil, fmt.Errorf("invalid specified references : %v", tags)
	}
	builder := &Builder{
		helper:           libmason.NewHelper(c),
		contextDirectory: contextDirectory,
		dockerfilePath:   dockerfilePath,
		references:       tags,
		out:              os.Stdout,
	}

	return builder, nil
}

// Run executes the build process.
func (b *Builder) Run() error {
	dockerfile, err := os.Open(b.dockerfilePath)
	if err != nil {
		return fmt.Errorf("unable to open Dockerfile: %s", err)
	}
	commands, err := parser.Parse(dockerfile)
	if err != nil {
		return fmt.Errorf("unable to parse Dockerfile: %s", err)
	}
	if len(commands) == 0 {
		return fmt.Errorf("no commands found in Dockerfile")
	}

	build := builder.WithSteps(builder.WithLogFunc(builder.NewBuilder(b.helper), log.Infof), b.toSteps(commands))

	image, err := build.Run(context.Background())
	if err != nil {
		return err
	}

	for _, ref := range b.references {
		log.Infof("Tag image %s with reference %s", image, ref)
		if err := b.helper.TagImage(context.Background(), image, ref); err != nil {
			return err
		}
	}

	return nil
}

func (b *Builder) toSteps(commands []*parser.Command) []builder.Step {
	steps := make([]builder.Step, len(commands))
	for i, command := range commands {
		var step builder.Step
		cmd, args := strings.ToUpper(command.Args[0]), command.Args[1:]
		switch cmd {
		case "FROM":
			step = &builder.FromStep{
				Reference: args[0],
			}
		case "COPY":
			step = builder.WithRemove(
				builder.WithCommit(
					builder.WithCreate(&CopyStep{
						srcPath:     args[0],
						destPath:    args[1],
						contextPath: b.contextDirectory,
					}, []string{}, []string{}, false),
				),
			)
		case "RUN":
			step = builder.WithRemove(
				builder.WithCommit(
					builder.WithCreate(&RunStep{
						heredoc: command.Heredoc,
					}, args[:1], args[1:], true),
				),
			)
		case "LABEL":
			step = builder.WithRemove(
				builder.WithCommit(&LabelStep{
					labels: map[string]string{
						args[0]: args[1],
					},
				}),
			)
		}
		steps[i] = step
	}
	return steps
}

func validateReferences(references []string) error {
	// FIXME(vdemeester) handle that using reference package from distribution
	return nil
}

func validateDockerfilePath(dockerfilePath string) error {
	_, err := os.Stat(dockerfilePath)
	return err
}

func validateContextDirectory(contextDirectory string) error {
	stat, err := os.Stat(contextDirectory)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("context must be a directory")
	}
	return nil
}
