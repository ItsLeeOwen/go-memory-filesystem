package filesystem

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
