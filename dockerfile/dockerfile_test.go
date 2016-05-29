package dockerfile

import (
	"testing"

	"github.com/vdemeester/libmason/test"
)

func TestNewBuilder(t *testing.T) {
	c := test.NewNopClient()
	contextDirectory := "."
	dockerfilePath := ""
	ref := []string{"something:fun"}
	_, err := NewBuilder(c, contextDirectory, dockerfilePath, ref)
	if err == nil {
		t.Fatal(err)
	}
}
