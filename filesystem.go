package filesystem

import (
	"fmt"
	"strings"
	"sync"
)

// FileSystemer is the interface for reading/writing data to an in-memory filesystem
type FileSystemer interface {
	Mkdir(string) error
	WriteFile(string, string) error
	ReadFile(string) (string, error)
}

// NewFileSystem returns a FileSystemer
func NewFileSystem() FileSystemer {
	nodes := &sync.Map{}
	nodes.Store("", &dir{
		name:  "",
		nodes: &sync.Map{},
	})

	return &fileSystem{
		root: &dir{
			name:  "CoinbaseHD",
			nodes: nodes,
		},
	}
}

// fileSystem implements FileSystemer
type fileSystem struct {
	root Dir
}

// Mkdir creates a directory at path or returns
// an err if the path is incomplete
func (fs fileSystem) Mkdir(path string) error {
	segments, segmentsLen := fs.parsePathSegments(path)
	parentPath := segments[:segmentsLen-1]

	parent, err := fs.root.Find(parentPath)
	if err != nil {
		return err
	}

	parentDir, isDir := parent.(Dir)
	if !isDir {
		return fmt.Errorf("path is invalid: %s", path)
	}

	dirName := segments[segmentsLen-1]

	// add child or return err
	if err := parentDir.CreateDir(dirName); err != nil {
		return fmt.Errorf("error adding child dir %s: %w", path, err)
	}

	return nil
}

// WriteFile creates a file if it doesn't exist
// and appends data to the file
func (fs fileSystem) WriteFile(path string, data string) error {
	segments, segmentsLen := fs.parsePathSegments(path)
	parentPath := segments[:segmentsLen-1]

	parent, err := fs.root.Find(parentPath)
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

// Readfile returns data in a file or error
func (fs fileSystem) ReadFile(path string) (string, error) {
	segments := strings.Split(path, "/")

	node, err := fs.root.Find(segments)
	if err != nil {
		return "", err
	}

	file, isFile := node.(File)
	if !isFile {
		return "", fmt.Errorf("error reading non-file: %s", path)
	}

	return file.Data(), nil
}

// parsePathSegments helper func to return info from path string
func (fs fileSystem) parsePathSegments(path string) (segments []string, segmentsLen int) {
	segments = strings.Split(path, "/")
	segmentsLen = len(segments)
	return
}

// Node is the interface grouping common functionality between File & Dir types
type Node interface {
	Name() string
}

// Dir interface for traversing & creating child Nodes
type Dir interface {
	Node
	CreateDir(name string) error
	CreateFile(name string, data string) error
	Find(path []string) (Node, error)
	GetChild(name string) (Node, error)
	HasChild(name string) bool
}

// File interface for reading/writing a File Node
type File interface {
	Node
	Append(data string)
	Data() string
}

// file implements File
type file struct {
	data string
	name string
}

// Append data to the File
func (f *file) Append(data string) {
	f.data += data
}

// Data return data in the File
func (f *file) Data() string {
	return f.data
}

// Name of the File
func (f *file) Name() string {
	return f.name
}

// dir implements Dir
type dir struct {
	name  string
	nodes *sync.Map
}

// Name of the Dir
func (d *dir) Name() string {
	return d.name
}

// CreateDir creates a child directory,
// returns an error if the dirname already exists
func (d *dir) CreateDir(name string) error {
	if d.HasChild(name) {
		return fmt.Errorf("%s dir already exists", name)
	}

	d.nodes.Store(name, &dir{
		name:  name,
		nodes: &sync.Map{},
	})

	return nil
}

// CreateFile creates a child file,
// returns an error if the filename already exists
func (d *dir) CreateFile(name string, data string) error {
	if d.HasChild(name) {
		return fmt.Errorf("file already exists: %s", name)
	}

	d.nodes.Store(name, &file{
		name: name,
		data: data,
	})

	return nil
}

// GetChild returns child by name
// returns an error if the child does not exist
func (d *dir) GetChild(name string) (Node, error) {
	child, exists := d.nodes.Load(name)

	if !exists {
		return nil, fmt.Errorf("child does not exist: %s", name)
	}

	return child.(Node), nil
}

// HasChild returns bool if dir has child of name
func (d *dir) HasChild(name string) bool {
	if _, exists := d.nodes.Load(name); !exists {
		return false
	}

	return true
}

// Find walks the list & returns a final Node or error
func (d *dir) Find(path []string) (Node, error) {
	name := path[0]

	child, err := d.GetChild(name)
	if err != nil {
		return nil, fmt.Errorf("%w; path: %s", err, path)
	}

	if len(path) == 1 {
		// we have found the final Node
		return child, nil
	}

	// to continue finding Node must be a Dir
	childDir, isDir := child.(Dir)
	if !isDir {
		return nil, fmt.Errorf("path is invalid: %s", path)
	}

	// shift one from front of the path
	// and continue finding the next Node
	return childDir.Find(path[1:])
}
