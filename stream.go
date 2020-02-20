<<<<<<< HEAD
package stream

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
)

//Stream - general intarface
type Stream struct {
	s      Streamer
	buf    [8]byte
	offset int64
	order  binary.ByteOrder
}

//Streamer - general intarface
type Streamer interface {
	io.ReaderAt
	io.Closer
	Size() int64
}

//New - create stream
func New(data interface{}) (*Stream, error) {
	var s Streamer
	var err error

	switch td := data.(type) {
	case *os.File:
		s, err = NewFileStream(td)
	case []byte:
		s, err = NewByteStream(td)
	default:
		err = errors.New("Unknown type")
	}
	if err != nil {
		return nil, err
	}
	return &Stream{s: s, order: binary.LittleEndian}, nil
}

//ChangeByteOrder - Set ByteOrder
func (s *Stream) ChangeByteOrder(order binary.ByteOrder) {
	s.order = order
}

//Close - close file stream
func (s *Stream) Close() error {
	if s == nil {
		return errors.New("Stream is nil")
	}
	if s.s == nil {
		return errors.New("Base stream is nil")
	}
	err := s.s.Close()
	s.offset = 0
	s.s = nil
	return err
}

//Size - bytes size
func (s *Stream) Size() int64 {
	if s == nil || s.s == nil {
		return 0
	}
	return s.s.Size()
}

//Seek - change position
func (s *Stream) Seek(offset int64, whence int) (int64, error) {
	position := s.offset
	switch whence {
	case 0:
		position = offset
	case 1:
		position += offset
	case 2:
		position = s.s.Size() - offset
	default:
		return -1, errors.New("Unknown whence")
	}

	if position < 0 || position > s.s.Size() {
		return -1, io.EOF
	}
	s.offset = position

	return s.offset, nil
}

//ReadAt - read bytes from file
func (s *Stream) ReadAt(b []byte, offset int64) (int, error) {
	if s.s == nil {
		return -1, errors.New("Satream is nil")
	}

	i, err := s.s.ReadAt(b, offset)
	if err == nil && i != len(b) {
		err = io.EOF
	}
	return i, err
}

//Read - read bytes from file
func (s *Stream) Read(b []byte) (int, error) {
	i, err := s.ReadAt(b, s.offset)
	if err != nil {
		return i, err
	}
	s.offset += int64(i)
	return i, err
}

//ReadByte - read byte from file
func (s *Stream) ReadByte() (byte, error) {
	_, err := s.Read(s.buf[:1])
	if err != nil {
		return s.buf[0], err
	}
	return s.buf[0], nil
}

//ReadBytes - read bytes from file
func (s *Stream) ReadBytes(count int) ([]byte, error) {
	result := make([]byte, count)
	if _, err := s.Read(result); err != nil {
		return nil, err
	}
	return result, nil
}

//ReadUInt16 - read uint16 from file
func (s *Stream) ReadUInt16() (v uint16, err error) {
	if _, err = s.Read(s.buf[:2]); err != nil {
		return 0, err
	}
	return s.order.Uint16(s.buf[:2]), nil
}

//ReadUInt32 - read uint32 from file
func (s *Stream) ReadUInt32() (uint32, error) {
	if _, err := s.Read(s.buf[:4]); err != nil {
		return 0, err
	}
	return s.order.Uint32(s.buf[:4]), nil
}

//ReadUInt64 - read uint64 from file
func (s *Stream) ReadUInt64() (uint64, error) {
	if _, err := s.Read(s.buf[:8]); err != nil {
		return 0, err
	}
	return s.order.Uint64(s.buf[:8]), nil
}

//ReadInt32 - read int32 from file
func (s *Stream) ReadInt32() (int32, error) {
	i, err := s.ReadUInt32()
	return int32(i), err
}

//ReadInt64 - read int64 from file
func (s *Stream) ReadInt64() (int64, error) {
	i, err := s.ReadUInt64()
	return int64(i), err
}

//SeekStart - change position relative to the start of the file,
func (s *Stream) SeekStart(offset int64) (int64, error) {
	return s.Seek(offset, 0)
}
=======
package stream

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
)

//Stream - general intarface
type Stream struct {
	s      Streamer
	buf    [8]byte
	offset int64
	order  binary.ByteOrder
}

//Streamer - general intarface
type Streamer interface {
	io.ReaderAt
	io.Closer
	Size() int64
}

//New - create stream
func New(data interface{}) (*Stream, error) {
	var s Streamer
	var err error

	switch td := data.(type) {
	case *os.File:
		s, err = NewFileStream(td)
	case []byte:
		s, err = NewByteStream(td)
	default:
		err = errors.New("Unknown type")
	}
	if err != nil {
		return nil, err
	}
	return &Stream{s: s, order: binary.LittleEndian}, nil
}

//ChangeByteOrder - Set ByteOrder
func (s *Stream) ChangeByteOrder(order binary.ByteOrder) {
	s.order = order
}

//Close - close file stream
func (s *Stream) Close() error {
	if s == nil {
		return errors.New("Stream is nil")
	}
	if s.s == nil {
		return errors.New("Base stream is nil")
	}
	err := s.s.Close()
	s.offset = 0
	s.s = nil
	return err
}

//Size - bytes size
func (s *Stream) Size() int64 {
	if s == nil || s.s == nil {
		return 0
	}
	return s.s.Size()
}

//Seek - change position
func (s *Stream) Seek(offset int64, whence int) (int64, error) {
	position := s.offset
	switch whence {
	case 0:
		position = offset
	case 1:
		position += offset
	case 2:
		position = s.s.Size() - offset
	default:
		return -1, errors.New("Unknown whence")
	}

	if position < 0 || position > s.s.Size() {
		return -1, io.EOF
	}
	s.offset = position

	return s.offset, nil
}

//ReadAt - read bytes from file
func (s *Stream) ReadAt(b []byte, offset int64) (int, error) {
	if s.s == nil {
		return -1, errors.New("Satream is nil")
	}

	i, err := s.s.ReadAt(b, offset)
	if err == nil && i != len(b) {
		err = io.EOF
	}
	return i, err
}

//Read - read bytes from file
func (s *Stream) Read(b []byte) (int, error) {
	i, err := s.ReadAt(b, s.offset)
	if err != nil {
		return i, err
	}
	s.offset += int64(i)
	return i, err
}

//ReadByte - read byte from file
func (s *Stream) ReadByte() (byte, error) {
	_, err := s.Read(s.buf[:1])
	if err != nil {
		return s.buf[0], err
	}
	return s.buf[0], nil
}

//ReadBytes - read bytes from file
func (s *Stream) ReadBytes(count int) ([]byte, error) {
	result := make([]byte, count)
	if _, err := s.Read(result); err != nil {
		return nil, err
	}
	return result, nil
}

//ReadUInt16 - read uint16 from file
func (s *Stream) ReadUInt16() (v uint16, err error) {
	if _, err = s.Read(s.buf[:2]); err != nil {
		return 0, err
	}
	return s.order.Uint16(s.buf[:2]), nil
}

//ReadUInt32 - read uint32 from file
func (s *Stream) ReadUInt32() (uint32, error) {
	if _, err := s.Read(s.buf[:4]); err != nil {
		return 0, err
	}
	return s.order.Uint32(s.buf[:4]), nil
}

//ReadUInt64 - read uint64 from file
func (s *Stream) ReadUInt64() (uint64, error) {
	if _, err := s.Read(s.buf[:8]); err != nil {
		return 0, err
	}
	return s.order.Uint64(s.buf[:8]), nil
}

//ReadInt32 - read int32 from file
func (s *Stream) ReadInt32() (int32, error) {
	i, err := s.ReadUInt32()
	return int32(i), err
}

//ReadInt64 - read int64 from file
func (s *Stream) ReadInt64() (int64, error) {
	i, err := s.ReadUInt64()
	return int64(i), err
}

//SeekStart - change position relative to the start of the file,
func (s *Stream) SeekStart(offset int64) (int64, error) {
	return s.Seek(offset, 0)
}
>>>>>>> 697c24c94b71ba05e8594b9f777cff7f44aa8412
