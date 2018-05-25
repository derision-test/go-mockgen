package paths

import (
	"go/build"
	"os"
)

func Gopath() string {
	if gopath := os.Getenv("GOPATH"); gopath != "" {
		return gopath
	}

	return build.Default.GOPATH
}
