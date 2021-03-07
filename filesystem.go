package filesystem

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

// FileSystem is the interface for reading/writing data to an in-memory filesystem
type FileSystem interface {
	Mkdir(string) error
	PrettyPrint() (string, error)
	ReadFile(string) (string, error)
	WriteFile(string, string) error
}

// Node is the interface grouping common functionality between File & Dir types
type Node interface {
	Name() string
}

type fileSystemOption func(o fileSystemOptions) fileSystemOptions

// DisablePFlag will not auto-generate non-existent
// sub-directories leading up to the given directory
// in Mkdir or WriteFile calls
func DisablePFlag() fileSystemOption {
	return func(opt fileSystemOptions) fileSystemOptions {
		opt.pFlagDisabled = true
		return opt
	}
}

type fileSystemOptions struct {
	// see DisablePFlag
	pFlagDisabled bool
}

// NewFileSystem returns a FileSystem
func NewFileSystem(opts ...fileSystemOption) FileSystem {
	opt := fileSystemOptions{}
	for _, o := range opts {
		opt = o(opt)
	}

	return &fileSystem{
		dir: &dir{
			name:  "",
			nodes: &sync.Map{},
			opt:   &opt,
		},
	}
}

// fileSystem implements the FileSystem interface
type fileSystem struct {
	*dir
}

// Mkdir creates a directory at path or returns an err
func (fs fileSystem) Mkdir(path string) error {
	segments, segmentsLen := fs.parsePathSegments(path)

	var err error
	var parent Node

	if !fs.opt.pFlagDisabled {
		_, err = fs.CreatePath(segments)
		return err
	}

	parentPath := segments[:segmentsLen-1]
	parent, err = fs.Find(parentPath)

	if err != nil {
		return err
	}

	parentDir, isDir := parent.(Dir)
	if !isDir {
		return fmt.Errorf("path is invalid: %s", path)
	}

	dirName := segments[segmentsLen-1]

	// add child or return err
	if _, err := parentDir.CreateDir(dirName); err != nil {
		return fmt.Errorf("error adding child dir %s: %w", path, err)
	}

	return nil
}

// PrettyPrint json string representation of the filesystem
func (fs fileSystem) PrettyPrint() (string, error) {
	m := fs.dir.PrettyPrint()

	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// Readfile returns contents of a file or error
func (fs fileSystem) ReadFile(path string) (string, error) {
	segments := strings.Split(path, "/")

	node, err := fs.Find(segments)
	if err != nil {
		return "", err
	}

	file, isFile := node.(File)
	if !isFile {
		return "", fmt.Errorf("error reading non-file: %s", path)
	}

	return file.Data(), nil
}

// WriteFile creates a file if it doesn't exist
// and appends data to the file
func (fs fileSystem) WriteFile(path string, data string) error {
	segments, segmentsLen := fs.parsePathSegments(path)
	parentPath := segments[:segmentsLen-1]

	var err error
	var parent Node

	if !fs.opt.pFlagDisabled {
		parent, err = fs.CreatePath(parentPath)
	} else {
		parent, err = fs.Find(parentPath)
	}

	if err != nil {
		return err
	}

	dir, isDir := parent.(Dir)
	if !isDir {
		return fmt.Errorf("path is invalid: %s", path)
	}

	fileName := segments[segmentsLen-1]

	child, err := dir.GetChild(fileName)
	if err != nil {
		// use static err type & check for 404 here
		if err := dir.CreateFile(fileName, data); err != nil {
			return fmt.Errorf("error writing file %s: %w", path, err)
		}

		return nil
	}

	file, isFile := child.(File)
	if !isFile {
		return fmt.Errorf("error writing to non-file: %s", path)
	}

	file.Append(data)

	return nil
}

// parsePathSegments helper func to return info from path string
func (fs fileSystem) parsePathSegments(path string) (segments []string, segmentsLen int) {
	segments = strings.Split(path, "/")
	segmentsLen = len(segments)
	return
}
