package stream

import (
	"errors"
	"os"
)

//FileStream - file stream
type FileStream struct {
	f    *os.File
	size int64
}

//NewFileStream - create file stream
func NewFileStream(f *os.File) (*FileStream, error) {
	if f == nil {
		return nil, errors.New("file is nil")
	}

	fs := &FileStream{f: f}

	srcFile, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	fs.size = srcFile.Size()

	return fs, nil
}

//Close - close file stream
func (fs *FileStream) Close() error {
	fs.f = nil
	fs.size = 0
	return nil
}

//ReadAt - read bytes from file
func (fs *FileStream) ReadAt(b []byte, offset int64) (int, error) {
	return fs.f.ReadAt(b, offset)
}

//Size - file size
func (fs *FileStream) Size() int64 {
	return fs.size
}
