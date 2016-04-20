package dockerfile

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
)

func (b *Builder) handleRun(args []string, heredoc string) error {
	if len(args) == 0 {
		return fmt.Errorf("RUN require at least one argument")
	}
	ctx := context.Background()
	// Create a container
	c, err := b.helper.ContainerCreate(ctx, types.ContainerCreateConfig{
		Config: &container.Config{
			Image:      b.currentImage,
			Entrypoint: args[:1],
			Cmd:        args[1:],
			OpenStdin:  true,
			StdinOnce:  true,
		},
	})
	if err != nil {
		return err
	}
	errChan := make(chan error)
	go func() {
		// FIXME(vdemeester) handle errors
		errChan <- b.helper.ContainerAttach(ctx, c.ID, strings.NewReader(heredoc), os.Stdout, os.Stderr)
	}()

	// Start the container
	if err := b.helper.ContainerStart(ctx, c.ID); err != nil {
		return err
	}

	if err := <-errChan; err != nil {
		return err
	}

	// Commit the container and remove it
	imageID, err := b.helper.ContainerCommit(ctx, c.ID, types.ContainerCommitOptions{
		Changes: []string{
			fmt.Sprintf("CMD %v", b.currentCmd),
			fmt.Sprintf("ENTRYPOINT %v", b.currentEntrypoint),
		},
	})
	if err != nil {
		return err
	}
	b.currentImage = imageID

	if err := b.helper.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{
		Force: true,
	}); err != nil {
		return err
	}
	return nil
}
