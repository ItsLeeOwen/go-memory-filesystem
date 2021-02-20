package filesystem

import (
	"fmt"
	"strings"
)

// PathNotFoundError for errors.Is
var PathNotFoundError = &pathNotFoundError{}

func NewPathNotFoundError(path []string) *pathNotFoundError {
	return &pathNotFoundError{
		path: path,
	}
}

// pathNotFoundError for restrictive (non-autogenerating path) mode
type pathNotFoundError struct {
	path []string
}

func (e *pathNotFoundError) Error() string {
	return fmt.Sprintf("path '%s' does not exist", strings.Join(e.path, "/"))
}

func (e *pathNotFoundError) Is(err error) bool {
	_, ok := err.(*pathNotFoundError)
	if !ok {
		return false
	}
	return true
}
