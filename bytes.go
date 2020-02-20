<<<<<<< HEAD
package stream

import (
	"errors"
	"io"
)

//ByteStream - bytes stream
type ByteStream struct {
	b []byte
}

//NewByteStream - create file stream
func NewByteStream(b []byte) (*ByteStream, error) {
	if b == nil {
		return nil, errors.New("bytes is nil")
	}
	return &ByteStream{b: b}, nil
}

//Close - close file stream
func (bs *ByteStream) Close() error {
	bs.b = nil
	return nil
}

//Size - bytes size
func (bs *ByteStream) Size() int64 {
	return int64(len(bs.b))
}

//ReadAt - read bytes from file
func (bs *ByteStream) ReadAt(b []byte, offset int64) (int, error) {
	var err error
	end := offset + int64(len(b))

	if offset < 0 || offset > bs.Size() {
		return 0, io.EOF
	} else if end > bs.Size() {
		end = bs.Size()
	}
	return copy(b, bs.b[offset:end]), err
}
=======
package stream

import (
	"errors"
	"io"
)

//ByteStream - bytes stream
type ByteStream struct {
	b []byte
}

//NewByteStream - create file stream
func NewByteStream(b []byte) (*ByteStream, error) {
	if b == nil {
		return nil, errors.New("bytes is nil")
	}
	return &ByteStream{b: b}, nil
}

//Close - close file stream
func (bs *ByteStream) Close() error {
	bs.b = nil
	return nil
}

//Size - bytes size
func (bs *ByteStream) Size() int64 {
	return int64(len(bs.b))
}

//ReadAt - read bytes from file
func (bs *ByteStream) ReadAt(b []byte, offset int64) (int, error) {
	var err error
	end := offset + int64(len(b))

	if offset < 0 || offset > bs.Size() {
		return 0, io.EOF
	} else if end > bs.Size() {
		end = bs.Size()
	}
	return copy(b, bs.b[offset:end]), err
}
>>>>>>> 697c24c94b71ba05e8594b9f777cff7f44aa8412
