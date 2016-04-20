package dockerfile

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
)

func (b *Builder) handleLabel(args []string, heredoc string) error {
	if len(args) != 2 {
		return fmt.Errorf("LABEL requires exactly two arguments")
	}

	ctx := context.Background()

	// Create a container
	c, err := b.helper.ContainerCreate(ctx, types.ContainerCreateConfig{
		Config: &container.Config{
			Image: b.currentImage,
			Labels: map[string]string{
				args[0]: args[1],
			},
		},
	})
	if err != nil {
		return err
	}

	// Commit the container and remove it
	imageID, err := b.helper.Commit(ctx, c.ID, types.ContainerCommitOptions{})
	if err != nil {
		return err
	}
	b.currentImage = imageID

	if err := b.helper.ContainerRm(ctx, c.ID, types.ContainerRemoveOptions{
		Force: true,
	}); err != nil {
		return err
	}
	return nil
}