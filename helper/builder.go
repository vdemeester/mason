package helper

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http/httputil"
	"os"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"

	// FIXME(vdemeester) Remove dependency to docker/docker
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/docker/pkg/term"
)

var _ APIBuilder = &Builder{}

// Builder is a client-side builder implementation. It satisfy APIBuilder interface.
type Builder struct {
	client client.APIClient
}

// NewBuilder creates a new helper builder from a docker client
func NewBuilder(cli client.APIClient) *Builder {
	return &Builder{
		client: cli,
	}
}

// GetParentImage looks up a Docker image referenced by `ref`.
func (b *Builder) GetParentImage(ctx context.Context, ref string, options types.ImagePullOptions) (types.ImageInspect, error) {
	if imageInspect, _, err := b.client.ImageInspectWithRaw(ctx, ref, false); err == nil {
		return imageInspect, nil
	}
	// FIXME(vdemeester) still ways to come
	// 1. parse name/reference of the image
	// 2. from that get the registry and get the AuthConfig
	// 3. Create a privilegedFunc if needed
	// 4. Call ImagePull
	// Try to pull it
	responseBody, err := b.client.ImagePull(ctx, ref, options)
	if err != nil {
		return types.ImageInspect{}, err
	}

	// FIXME(vdemeester) do something better
	var writeBuff io.Writer = os.Stdout
	outFd, isTerminalOut := term.GetFdInfo(os.Stdout)

	defer responseBody.Close()
	if err := jsonmessage.DisplayJSONMessagesStream(responseBody, writeBuff, outFd, isTerminalOut, nil); err != nil {
		return types.ImageInspect{}, err
	}
	imageInspect, _, err := b.client.ImageInspectWithRaw(context.Background(), ref, false)
	return imageInspect, err
}

// TagImage tags an image with newTag
func (b *Builder) TagImage(ctx context.Context, image string, newReference string) error {
	// FIXME(vdemeester) Use reference for ImageTag \o/
	return b.client.ImageTag(ctx, image, newReference, types.ImageTagOptions{
		Force: true,
	})
}

// ContainerAttachRaw attaches to container.
func (b *Builder) ContainerAttachRaw(ctx context.Context, container string, stdin io.Reader, stdout, stderr io.Writer, stream bool) error {
	// pipe stdin, stderr and stdout (and stream) in containerAttachOptions
	resp, errAttach := b.client.ContainerAttach(ctx, container, types.ContainerAttachOptions{
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Stream: stream,
	})
	if errAttach != nil && errAttach != httputil.ErrPersistEOF {
		return errAttach
	}
	defer resp.Close()

	// FIXME(vdemeester) TTY stuff
	if err := holdHijackedConnection(ctx, true, stdin, stdout, stderr, resp); err != nil {
		return err
	}

	if errAttach != nil {
		return errAttach
	}
	return nil
}

// ContainerCreate creates a new Docker container and returns potential warnings
func (b *Builder) ContainerCreate(ctx context.Context, createConfig types.ContainerCreateConfig) (types.ContainerCreateResponse, error) {
	return b.client.ContainerCreate(context.Background(), createConfig.Config, createConfig.HostConfig, createConfig.NetworkingConfig, createConfig.Name)
}

// ContainerRm removes a container specified by `id`.
func (b *Builder) ContainerRm(ctx context.Context, container string, options types.ContainerRemoveOptions) error {
	return b.client.ContainerRemove(context.Background(), container, options)
}

// Commit creates a new Docker image from an existing Docker container.
func (b *Builder) Commit(ctx context.Context, container string, options types.ContainerCommitOptions) (string, error) {
	commitResponse, err := b.client.ContainerCommit(context.Background(), container, options)
	return commitResponse.ID, err
}

// ContainerKill stops the container execution abruptly.
func (b *Builder) ContainerKill(ctx context.Context, containerID string, sig uint64) error {
	return b.client.ContainerKill(context.Background(), containerID, string(sig))
}

// ContainerStart starts a new container
func (b *Builder) ContainerStart(ctx context.Context, container string) error {
	return b.client.ContainerStart(context.Background(), container)
}

// ContainerWait stops processing until the given container is stopped.
func (b *Builder) ContainerWait(ctx context.Context, container string, timeout time.Duration) (int, error) {
	// FIXME(vdemeester) TODO
	return 0, nil
}

// ContainerUpdateCmdOnBuild updates container.Path and container.Args
func (b *Builder) ContainerUpdateCmdOnBuild(ctx context.Context, container string, cmd []string) error {
	return nil
}

// CopyOnBuild copies/extracts a source FileInfo to a destination path inside a container
// specified by a container object.
func (b *Builder) CopyOnBuild(ctx context.Context, container string, destPath, srcPath string, decompress bool) error {
	dstInfo := archive.CopyInfo{Path: destPath}
	// FIXME(vdemeester) handle link follow here ?
	dstStat, err := b.client.ContainerStatPath(ctx, container, destPath)
	if err == nil {
		dstInfo.Exists, dstInfo.IsDir = true, dstStat.Mode.IsDir()
	}
	srcInfo, err := archive.CopyInfoSourcePath(srcPath, false)
	if err != nil {
		return err
	}
	srcArchive, err := archive.TarResource(srcInfo)
	if err != nil {
		return err
	}
	defer srcArchive.Close()
	destDir, preparedArchive, err := archive.PrepareArchiveCopy(srcArchive, srcInfo, dstInfo)
	if err != nil {
		return err
	}
	defer preparedArchive.Close()
	// FIXME(vdemeester) update signature
	return b.client.CopyToContainer(ctx, container, destDir, preparedArchive, types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	})
}

// FIXME(vdemeester) Handle context :)
func holdHijackedConnection(ctx context.Context, tty bool, inputStream io.Reader, outputStream, errorStream io.Writer, resp types.HijackedResponse) error {
	var err error
	receiveStdout := make(chan error, 1)
	if outputStream != nil || errorStream != nil {
		go func() {
			// When TTY is ON, use regular copy
			if tty && outputStream != nil {
				_, err = io.Copy(outputStream, resp.Reader)
			} else {
				_, err = stdcopy.StdCopy(outputStream, errorStream, resp.Reader)
			}
			receiveStdout <- err
		}()
	}

	stdinDone := make(chan struct{})
	go func() {
		if inputStream != nil {
			io.Copy(resp.Conn, inputStream)
		}

		if err := resp.CloseWrite(); err != nil {
		}
		close(stdinDone)
	}()

	select {
	case err := <-receiveStdout:
		if err != nil {
			return err
		}
	case <-stdinDone:
		if outputStream != nil || errorStream != nil {
			if err := <-receiveStdout; err != nil {
				return err
			}
		}
	}

	return nil
}

// encodeAuthToBase64 serializes the auth configuration as JSON base64 payload
func encodeAuthToBase64(authConfig types.AuthConfig) (string, error) {
	buf, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(buf), nil
}
