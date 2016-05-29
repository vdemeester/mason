package dockerfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/vdemeester/libmason"
	"github.com/vdemeester/libmason/builder"
)

// LabelStep is a step for LABEL.
type LabelStep struct {
	labels map[string]string
}

func (s *LabelStep) String() string {
	return fmt.Sprintf("LABEL %v", s.labels)
}

// Execute implements Step.Execute. It executes the step based on the specified config and helper.
func (s *LabelStep) Execute(ctx context.Context, helper libmason.Helper, config *builder.Config) (*builder.Config, error) {
	c, err := helper.ContainerCreate(ctx, types.ContainerCreateConfig{
		Config: &container.Config{
			Image:  config.ImageID,
			Labels: s.labels,
		},
	})
	if err != nil {
		return nil, err
	}

	config.Put(builder.ContainerID, c.ID)
	return config, nil
}

// RunStep is a step for RUN.
type RunStep struct {
	heredoc string
}

func (s *RunStep) String() string {
	return fmt.Sprintf("RUN %v", s.heredoc)
}

// Execute implements Step.Execute. It executes the step based on the specified config and helper.
func (s *RunStep) Execute(ctx context.Context, helper libmason.Helper, config *builder.Config) (*builder.Config, error) {
	containerID, ok := config.Get(builder.ContainerID)
	if !ok {
		return nil, fmt.Errorf("%s missing in config, cannot commit the container", builder.ContainerID)
	}
	errChan := make(chan error)
	go func() {
		// FIXME(vdemeester) handle errors
		errChan <- helper.ContainerAttach(ctx, containerID.(string), strings.NewReader(s.heredoc), os.Stdout, os.Stderr)
	}()

	// Start the container
	if err := helper.ContainerStart(ctx, containerID.(string)); err != nil {
		return nil, err
	}

	if err := <-errChan; err != nil {
		return nil, err
	}

	return config, nil
}

// CopyStep is a step for COPY.
type CopyStep struct {
	srcPath     string
	destPath    string
	contextPath string
}

func (s *CopyStep) String() string {
	return fmt.Sprintf("COPY %s %s", s.srcPath, s.destPath)
}

// Execute implements Step.Execute. It executes the step based on the specified config and helper.
func (s *CopyStep) Execute(ctx context.Context, helper libmason.Helper, config *builder.Config) (*builder.Config, error) {
	containerID, ok := config.Get(builder.ContainerID)
	if !ok {
		return nil, fmt.Errorf("%s missing in config, cannot commit the container", builder.ContainerID)
	}
	srcPath := filepath.Join(s.contextPath, s.srcPath)
	if err := helper.CopyToContainer(ctx, containerID.(string), s.destPath, srcPath, false); err != nil {
		return nil, err
	}
	return config, nil
}
