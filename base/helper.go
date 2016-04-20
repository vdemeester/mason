package base

import (
	"io"
	"os"

	"github.com/docker/engine-api/client"
)

var _ Helper = &DefaultHelper{}

// DefaultHelper is a client-side builder base helper implementation.
type DefaultHelper struct {
	client       client.APIClient
	outputWriter io.Writer
}

// NewHelper creates a new Helper from a docker client
func NewHelper(cli client.APIClient) *DefaultHelper {
	return &DefaultHelper{
		client:       cli,
		outputWriter: os.Stdout,
	}
}

// WithOutputWriter lets you specify a writer for the small amount of output this
// package will generate (Pull & such)
func (h *DefaultHelper) WithOutputWriter(w io.Writer) *DefaultHelper {
	h.outputWriter = w
	return h
}
