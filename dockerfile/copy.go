package dockerfile

import (
	"fmt"
	"path/filepath"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
)

func (b *Builder) handleCopy(args []string, heredoc string) error {
	if len(args) != 2 {
		return fmt.Errorf("COPY require 2 arguments")
	}
	ctx := context.Background()

	// Create a container
	c, err := b.helper.ContainerCreate(ctx, types.ContainerCreateConfig{
		Config: &container.Config{
			Image: b.currentImage,
		},
	})
	if err != nil {
		return err
	}

	// Copy \o/
	// FIXME(vdemeester) handle this way better, and compress it
	destPath := args[1]
	srcPath := filepath.Join(b.contextDirectory, args[0])
	if err := b.helper.CopyToContainer(ctx, c.ID, destPath, srcPath, false); err != nil {
		return err
	}

	// Commit the container and remove it
	imageID, err := b.helper.ContainerCommit(ctx, c.ID, types.ContainerCommitOptions{})
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
