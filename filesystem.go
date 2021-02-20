package filesystem

import (
	"encoding/json"
	"errors"
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

// Dir interface for traversing & creating child Nodes
type Dir interface {
	Node
	CreateDir(name string) (Dir, error)
	CreateFile(name string, data string) error
	CreatePath(path []string) (Dir, error)
	Find(path []string) (Node, error)
	GetChild(name string) (Node, error)
	HasChild(name string) bool
	PrettyPrint() map[string]interface{}
}

// File interface for reading/writing a File Node
type File interface {
	Node
	Append(data string)
	Data() string
}

// Node is the interface grouping common functionality between File & Dir types
type Node interface {
	Name() string
}

// NewFileSystem returns a FileSystem
func NewFileSystem(opts ...FileSystemOption) FileSystem {
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

type FileSystemOption func(o fileSystemOptions) fileSystemOptions

// DisablePFlag will not auto-generate non-existent
// sub-directories leading up to the given directory
// in Mkdir or WriteFile calls
func DisablePFlag() FileSystemOption {
	return func(opt fileSystemOptions) fileSystemOptions {
		opt.pFlagDisabled = true
		return opt
	}
}

type fileSystemOptions struct {
	// see DisablePFlag
	pFlagDisabled bool
}

// fileSystem implements FileSystemer
type fileSystem struct {
	*dir
}

// Mkdir creates a directory at path or returns
// an err if the path is incomplete
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

// Readfile returns data in a file or error
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

func (fs fileSystem) PrettyPrint() (string, error) {
	m := fs.dir.PrettyPrint()

	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// parsePathSegments helper func to return info from path string
func (fs fileSystem) parsePathSegments(path string) (segments []string, segmentsLen int) {
	segments = strings.Split(path, "/")
	segmentsLen = len(segments)
	return
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
	opt   *fileSystemOptions
	name  string
	nodes *sync.Map
}

// Name of the Dir
func (d *dir) Name() string {
	return d.name
}

// CreateDir creates a child directory,
// returns an error if the dirname already exists
func (d *dir) CreateDir(name string) (Dir, error) {
	if d.HasChild(name) {
		return nil, fmt.Errorf("%s dir already exists", name)
	}

	dir := &dir{
		name:  name,
		nodes: &sync.Map{},
		opt:   d.opt,
	}

	d.nodes.Store(name, dir)

	return dir, nil
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

// CreatePath walks the Trie creating any missing directories
// & returns the final Dir or error
func (d *dir) CreatePath(path []string) (Dir, error) {
	name := path[0]

	child, err := d.GetChild(name)

	if err != nil {
		if !errors.Is(err, PathNotFoundError) {
			return nil, err
		}

		if child, err = d.CreateDir(name); err != nil {
			return nil, err
		}
	}

	// Node must be a Dir
	childDir, isDir := child.(Dir)
	if !isDir {
		return nil, fmt.Errorf("path is invalid: %s", path)
	}

	if len(path) == 1 {
		// return the final Dir
		return childDir, nil
	}

	// shift one from front of the path
	// and continue finding the next Node
	return childDir.CreatePath(path[1:])
}

// Find walks the Trie & returns a final Node or error
func (d *dir) Find(path []string) (Node, error) {
	name := path[0]

	child, err := d.GetChild(name)
	if err != nil {
		return nil, NewPathNotFoundError(path)
	}

	if len(path) == 1 {
		// return the final Node
		return child, nil
	}

	// to continue finding, Node must be a Dir
	childDir, isDir := child.(Dir)
	if !isDir {
		return nil, fmt.Errorf("path is invalid: %s", path)
	}

	// shift one from front of the path
	// and continue finding the next Node
	return childDir.Find(path[1:])
}

// GetChild returns child by name
// returns an error if the child does not exist
func (d *dir) GetChild(name string) (Node, error) {
	child, exists := d.nodes.Load(name)

	if !exists {
		return nil, NewPathNotFoundError([]string{name})
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

func (d *dir) PrettyPrint() map[string]interface{} {
	m := map[string]interface{}{}

	d.nodes.Range(func(key, value interface{}) bool {
		if d, ok := value.(Dir); ok {
			m[fmt.Sprint(key)] = d.PrettyPrint()
		} else if f, ok := value.(File); ok {
			m[fmt.Sprint(key)] = f.Data()
		} else {
			fmt.Println("wtf is this?", value)
		}
		return true
	})

	return m
}

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
