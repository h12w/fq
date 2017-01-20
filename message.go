package sej

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	metaSize = 25
)

var (
	errMessageCorrupted = errors.New("last message of the journal file is courrupted")
)

// Message in a segmented journal file
type Message struct {
	Offset    uint64
	Timestamp time.Time
	Type      byte
	Value     []byte
}

type readSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}

// WriteTo writes the message to w
func (m *Message) WriteTo(w io.Writer) (int64, error) {
	cnt := int64(0) // total bytes written
	msgLen := int32(len(m.Value))

	n, err := writeUint64(w, m.Offset)
	cnt += int64(n)
	if err != nil {
		return cnt, err
	}

	ts := m.Timestamp
	if ts.IsZero() {
		ts = time.Now()
	}
	n, err = writeInt64(w, ts.UnixNano())
	cnt += int64(n)
	if err != nil {
		return cnt, err
	}

	n, err = writeByte(w, m.Type)
	cnt += int64(n)
	if err != nil {
		return cnt, err
	}

	n, err = writeInt32(w, msgLen)
	cnt += int64(n)
	if err != nil {
		return cnt, err
	}

	n, err = w.Write(m.Value)
	cnt += int64(n)
	if err != nil {
		return cnt, err
	}

	n, err = writeInt32(w, msgLen)
	cnt += int64(n)
	if err != nil {
		return cnt, err
	}

	return cnt, nil
}

// ReadFrom reads a message from a io.ReadSeeker.
// When an error occurs, it will rollback the seeker and then returns the original error.
func (m *Message) ReadFrom(r io.ReadSeeker) (n int64, err error) {
	cnt := int64(0) // total bytes read
	var msgLen, msgLen2 int32

	nn, err := readUint64(r, &m.Offset)
	cnt += int64(nn)
	if err != nil {
		return cnt, err
	}

	var unixNano int64
	nn, err = readInt64(r, &unixNano)
	cnt += int64(nn)
	if err != nil {
		return cnt, err
	}
	m.Timestamp = time.Unix(0, unixNano)

	nn, err = readByte(r, &m.Type)
	cnt += int64(nn)
	if err != nil {
		return cnt, err
	}

	nn, err = readInt32(r, &msgLen)
	cnt += int64(nn)
	if err != nil {
		return cnt, err
	}

	m.Value = make([]byte, int(msgLen))
	nn, err = io.ReadFull(r, m.Value)
	cnt += int64(nn)
	if err != nil {
		return cnt, err
	}
	if nn != int(msgLen) {
		return cnt, fmt.Errorf("message is truncated at %d", m.Offset)
	}

	nn, err = readInt32(r, &msgLen2)
	cnt += int64(nn)
	if err != nil {
		return cnt, err
	}
	if msgLen != msgLen2 {
		return cnt, fmt.Errorf("data corruption detected by size2 at %d", m.Offset)
	}

	return cnt, nil
}

func readMessageBackward(r io.ReadSeeker) (*Message, error) {
	var msgLen int32
	if _, err := r.Seek(-4, os.SEEK_CUR); err != nil {
		return nil, err
	}
	if _, err := readInt32(r, &msgLen); err != nil {
		return nil, err
	}
	if _, err := r.Seek(-metaSize-int64(msgLen), os.SEEK_CUR); err != nil {
		return nil, err
	}
	var msg Message
	_, err := msg.ReadFrom(r)
	return &msg, err
}

func writeByte(w io.Writer, i byte) (int, error) {
	return w.Write([]byte{i})
}

func readByte(r io.ReadSeeker, i *byte) (int, error) {
	var b [1]byte
	n, err := io.ReadFull(r, b[:])
	if err != nil {
		return n, err
	}
	*i = b[0]
	return n, nil
}

func writeInt64(w io.Writer, i int64) (int, error) {
	return w.Write([]byte{byte(i >> 56), byte(i >> 48), byte(i >> 40), byte(i >> 32), byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)})
}

func readInt64(r io.ReadSeeker, i *int64) (int, error) {
	var b [8]byte
	n, err := io.ReadFull(r, b[:])
	if err != nil {
		return n, err
	}
	*i = int64(b[0])<<56 | int64(b[1])<<48 | int64(b[2])<<40 | int64(b[3])<<32 |
		int64(b[4])<<24 | int64(b[5])<<16 | int64(b[6])<<8 | int64(b[7])
	return n, nil
}

func writeUint64(w io.Writer, i uint64) (int, error) {
	return w.Write([]byte{byte(i >> 56), byte(i >> 48), byte(i >> 40), byte(i >> 32), byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)})
}

func readUint64(r io.ReadSeeker, i *uint64) (int, error) {
	var b [8]byte
	n, err := io.ReadFull(r, b[:])
	if err != nil {
		return n, err
	}
	*i = uint64(b[0])<<56 | uint64(b[1])<<48 | uint64(b[2])<<40 | uint64(b[3])<<32 |
		uint64(b[4])<<24 | uint64(b[5])<<16 | uint64(b[6])<<8 | uint64(b[7])
	return n, nil
}

func writeInt32(w io.Writer, i int32) (int, error) {
	return w.Write([]byte{byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)})
}

func readInt32(r io.ReadSeeker, i *int32) (int, error) {
	var b [4]byte
	n, err := io.ReadFull(r, b[:])
	if err != nil {
		return n, err
	}
	*i = int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])
	return n, nil
}

func writeUint32(w io.Writer, i uint32) (int, error) {
	return w.Write([]byte{byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)})
}

func readUint32(r io.ReadSeeker, i *uint32) (int, error) {
	var b [4]byte
	n, err := io.ReadFull(r, b[:])
	if err != nil {
		return n, err
	}
	*i = uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
	return n, nil
}

// LatestOffset returns the offset after the last message in a journal file
func (journalFile *JournalFile) LastOffset() (uint64, error) {
	file, err := os.Open(journalFile.FileName)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	fileSize, err := file.Seek(0, os.SEEK_END)
	if err != nil {
		return 0, err
	}
	if fileSize == 0 {
		return journalFile.FirstOffset, nil
	}
	msg, err := readMessageBackward(file)
	if err != nil {
		return 0, errMessageCorrupted
	}
	return msg.Offset + 1, nil
}
