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
	"github.com/docker/engine-api/types/strslice"
	"github.com/vdemeester/mason/base"
	"github.com/vdemeester/mason/dockerfile/parser"
)

// DefaultDockerfile holds the default name for a Dockerfile
const DefaultDockerfile = "Dockerfile"

type handlerFunc func(args []string, heredoc string) error

// Builder holds attributes and defines method to build Dockerfile
type Builder struct {
	helper           base.Helper
	contextDirectory string
	dockerfilePath   string
	references       []string
	out              io.Writer

	currentImage      string
	currentEntrypoint strslice.StrSlice
	currentCmd        strslice.StrSlice

	handlers map[string]handlerFunc
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
		helper:           base.NewHelper(c),
		contextDirectory: contextDirectory,
		dockerfilePath:   dockerfilePath,
		references:       tags,
		out:              os.Stdout,
	}

	builder.handlers = map[string]handlerFunc{
		"FROM":  builder.handleFrom,
		"LABEL": builder.handleLabel,
		"RUN":   builder.handleRun,
		"COPY":  builder.handleCopy,
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

	for stepNum, command := range commands {
		if err := b.dispatch(stepNum, command); err != nil {
			return err
		}
	}

	for _, ref := range b.references {
		log.Infof("Tag image %s with reference %s", b.currentImage, ref)
		if err := b.helper.TagImage(context.Background(), b.currentImage, ref); err != nil {
			return err
		}
	}

	return nil
}

func (b *Builder) dispatch(stepNum int, command *parser.Command) error {
	cmd, args := strings.ToUpper(command.Args[0]), command.Args[1:]

	if (stepNum == 0) != (cmd == "FROM") {
		return fmt.Errorf("FROM must be the first Dockerfile command")
	}

	handler, exists := b.handlers[cmd]
	if !exists {
		return fmt.Errorf("unknown command: %q", cmd)
	}

	log.Infof("Step %d: %#v\n", stepNum, command)
	// FIXME(vdemeester) do way more..
	if err := handler(args, command.Heredoc); err != nil {
		return err
	}
	log.Infof("--> %s", b.currentImage)

	return nil
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
