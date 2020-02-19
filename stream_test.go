package stream

import (
	"encoding/binary"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func NewFile(filename string) (f *os.File, err error) {
	return os.OpenFile(filename, os.O_RDWR, 0644)
}

func NewBytes(filename string) (b []byte, err error) {
	return ioutil.ReadFile(filename)
}

func NewArrayStream() ([]*Stream, error) {
	arr := make([]*Stream, 0, 2)
	f, err := NewFile("./test.txt")
	if err != nil {
		return arr, err
	}
	s1, err := New(f)
	if err != nil {
		return arr, err
	}
	arr = append(arr, s1)
	b, err := NewBytes("./test.txt")
	if err != nil {
		return arr, err
	}
	s2, err := New(b)
	if err != nil {
		return arr, err
	}
	if err != nil {
		return arr, err
	}
	arr = append(arr, s2)
	return arr, err
}

func TestNew(t *testing.T) {
	fs, err := NewFileStream(nil)
	assert.Error(t, err, "should be error")
	assert.Nil(t, fs, "should be nil")

	err = CopyFile("./test.txt", "./delete.txt")
	assert.NoError(t, err, "should not be error")

	f, err := NewFile("./delete.txt")
	assert.NoError(t, err, "should not be error")
	assert.NotNil(t, f, "should not be nil")
	f.Close()
	err = os.Remove("./delete.txt")
	assert.NoError(t, err, "should not be error")

	fs, err = NewFileStream(f)
	assert.Error(t, err, "should be error")
	assert.Nil(t, fs, "should be nil")

	bs, err := NewByteStream(nil)
	assert.Error(t, err, "should be error")
	assert.Nil(t, bs, "should be nil")

	s, err := New(nil)
	assert.Error(t, err, "should be error")
	assert.Nil(t, s, "should be nil")

}

func TestClose(t *testing.T) {
	arr, err := NewArrayStream()
	assert.NoError(t, err, "should not be error")

	for _, s := range arr {
		s.Close()

		assert.Equal(t, s.offset, int64(0), "they should be equal")
		assert.Equal(t, s.Size(), int64(0), "they should be equal")
	}
	var s *Stream
	err = s.Close()
	assert.Error(t, err, "should be error")

	s = &Stream{}
	err = s.Close()
	assert.Error(t, err, "should be error")

	i := s.Size()
	assert.Equal(t, i, int64(0), "they should be equal")
}

func TestNewStream(t *testing.T) {
	f, err := NewFile("./test.txt")
	assert.NoError(t, err, "should not be error")
	defer f.Close()

	fs, err := NewFileStream(f)
	assert.NoError(t, err, "should not be error")
	assert.NotNil(t, fs, "should not be nil")
	assert.Equal(t, fs.Size(), int64(9), "they should be equal")

	fs.Close()

}

func TestSeek(t *testing.T) {
	arr, err := NewArrayStream()
	assert.NoError(t, err, "should not be error")

	for _, s := range arr {
		assert.Equal(t, s.Size(), int64(9), "they should be equal")

		i, err := s.SeekStart(4)
		assert.Equal(t, i, int64(4), "they should be equal")
		assert.NoError(t, err, "should not be error")

		i, err = s.SeekStart(6)
		assert.Equal(t, i, int64(6), "they should be equal")
		assert.NoError(t, err, "should not be error")

		i, err = s.SeekStart(20)
		assert.NotEqual(t, i, (20), "they should not be equal")
		assert.Error(t, err, "should be error")

		i, err = s.Seek(-1, -1)
		assert.Error(t, err, "should be error")

		i, err = s.Seek(3, 2)
		assert.Equal(t, i, int64(6), "they should be equal")
		assert.NoError(t, err, "should not be error")

		i, err = s.Seek(1, 1)
		assert.Equal(t, i, int64(7), "they should be equal")
		assert.NoError(t, err, "should not be error")

		s.Close()
	}
}

func TestRead(t *testing.T) {
	arr, err := NewArrayStream()
	assert.NoError(t, err, "should not be error")

	for _, s := range arr {
		_, err := s.SeekStart(5)
		assert.NoError(t, err, "should not be error")

		b, err := s.ReadBytes(3)
		assert.NoError(t, err, "should not be error")
		assert.Equal(t, b, []byte{'f', 'i', 'l'}, "they should be equal")

		b, err = s.ReadBytes(2)
		assert.Error(t, err, "should be error")
		assert.Equal(t, err, io.EOF, "they should be equal")

		b = make([]byte, 2)
		var i int
		i, err = s.ReadAt(b, 8)
		assert.Error(t, err, "should be error")
		assert.Equal(t, err, io.EOF, "they should be equal")
		assert.Equal(t, i, 1, "they should be equal")

		_, err = s.SeekStart(1)
		assert.NoError(t, err, "should not be error")
		bt, err := s.ReadByte()
		assert.NoError(t, err, "should not be error")
		assert.Equal(t, bt, uint8('e'), "they should be equal")

		b, err = s.ReadBytes(2)
		assert.NoError(t, err, "should not be error")
		assert.Equal(t, b, []byte{'s', 't'}, "they should be equal")

		_, err = s.Seek(0, 2)
		assert.NoError(t, err, "should not be error")
		_, err = s.ReadByte()
		assert.Error(t, err, "should be error")

		i, err = s.ReadAt(b, 20)
		assert.Error(t, err, "should be error")
		assert.Equal(t, err, io.EOF, "they should be equal")

		s.Close()
	}
}

func TestReadInt(t *testing.T) {
	arr, err := NewArrayStream()
	assert.NoError(t, err, "should not be error")

	for _, s := range arr {
		//16
		_, err := s.SeekStart(2)
		assert.NoError(t, err, "should not be error")

		u16, err := s.ReadUInt16()
		assert.NoError(t, err, "should not be error")

		s.SeekStart(2)
		mySlice, _ := s.ReadBytes(2)
		assert.Equal(t, u16, binary.LittleEndian.Uint16(mySlice), "they should be equal")

		//32
		s.SeekStart(2)
		u32, err := s.ReadUInt32()
		assert.NoError(t, err, "should not be error")

		s.SeekStart(2)
		mySlice, _ = s.ReadBytes(4)
		assert.Equal(t, u32, binary.LittleEndian.Uint32(mySlice), "they should be equal")

		s.SeekStart(2)
		i32, err := s.ReadInt32()
		assert.NoError(t, err, "should not be error")
		assert.Equal(t, i32, int32(binary.LittleEndian.Uint32(mySlice)), "they should be equal")

		//64
		s.SeekStart(0)
		u64, err := s.ReadUInt64()
		assert.NoError(t, err, "should not be error")

		s.SeekStart(0)
		mySlice, _ = s.ReadBytes(8)
		assert.Equal(t, u64, binary.LittleEndian.Uint64(mySlice), "they should be equal")

		s.SeekStart(0)
		i64, err := s.ReadInt64()
		assert.NoError(t, err, "should not be error")
		assert.Equal(t, i64, int64(binary.LittleEndian.Uint64(mySlice)), "they should be equal")

		s.Seek(0, 2)
		_, err = s.ReadInt64()
		assert.Error(t, err, "should be error")

		_, err = s.ReadInt32()
		assert.Error(t, err, "should be error")

		_, err = s.ReadUInt16()
		assert.Error(t, err, "should be error")

		s.ChangeByteOrder(binary.BigEndian)
		s.SeekStart(0)
		i64, err = s.ReadInt64()
		assert.NoError(t, err, "should not be error")
		assert.Equal(t, i64, int64(binary.BigEndian.Uint64(mySlice)), "they should be equal")

		s.Close()

		_, err = s.ReadInt64()
		assert.Error(t, err, "should be error")
	}
}
