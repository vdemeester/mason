package helper

import (
	"io"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
)

// APIBuilder abstracts calls to a Docker Daemon.
// FIXME(vdemeester) this is an extract a tiny bit modified from builder/builder.go
// from docker engine project.
type APIBuilder interface {
	// TODO: use digest reference instead of name

	// GetParentImage looks up a Docker image referenced by `name` and pull it if needed.
	GetParentImage(ctx context.Context, name string, options types.ImagePullOptions) (types.ImageInspect, error)

	// TagImage tags an image with newTag
	TagImage(ctx context.Context, image string, newReference string) error

	// ContainerAttachRaw attaches to container.
	ContainerAttachRaw(ctx context.Context, cID string, stdin io.Reader, stdout, stderr io.Writer, stream bool) error

	// ContainerCreate creates a new Docker container and returns potential warnings
	ContainerCreate(ctx context.Context, config types.ContainerCreateConfig) (types.ContainerCreateResponse, error)

	// ContainerRm removes a container specified by `id`.
	ContainerRm(ctx context.Context, name string, options types.ContainerRemoveOptions) error

	// Commit creates a new Docker image from an existing Docker container.
	Commit(ctx context.Context, name string, options types.ContainerCommitOptions) (string, error)
	//Commit(string, *backend.ContainerCommitConfig) (string, error)

	// ContainerKill stops the container execution abruptly.
	ContainerKill(ctx context.Context, containerID string, sig uint64) error

	// ContainerStart starts a new container
	ContainerStart(ctx context.Context, containerID string) error

	// ContainerWait stops processing until the given container is stopped.
	ContainerWait(ctx context.Context, containerID string, timeout time.Duration) (int, error)

	// ContainerUpdateCmdOnBuild updates container.Path and container.Args
	ContainerUpdateCmdOnBuild(ctx context.Context, containerID string, cmd []string) error

	// ContainerCopy copies/extracts a source FileInfo to a destination path inside a container
	// specified by a container object.
	// TODO: make an Extract method instead of passing `decompress`
	// TODO: do not pass a FileInfo, instead refactor the archive package to export a Walk function that can be used
	// with Context.Walk
	//ContainerCopy(name string, res string) (io.ReadCloser, error)
	// TODO: use copyBackend api
	CopyOnBuild(ctx context.Context, container string, destPath, srcPath string, decompress bool) error
}
