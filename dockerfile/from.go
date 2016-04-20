package dockerfile

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
)

func (b *Builder) handleFrom(args []string, heredoc string) error {
	if len(args) != 1 {
		return fmt.Errorf("FROM support only one arguments, got %#v", args)
	}

	image, err := b.helper.GetImage(context.Background(), args[0], types.ImagePullOptions{})
	if err != nil {
		return err
	}
	b.currentImage = image.ID
	b.currentCmd = image.Config.Cmd
	b.currentEntrypoint = image.Config.Entrypoint

	return nil
}
