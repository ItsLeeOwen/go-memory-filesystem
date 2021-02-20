package filesystem

import (
	"errors"
	"fmt"
	"sync"
)

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
// returns an error if the file already exists
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
		}
		return true
	})

	return m
}
