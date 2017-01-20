package sej

import (
	"fmt"
	"io"
	"os"
	"time"
)

var (
	// NotifyTimeout is the timeout value in rare cases that the OS notification fails
	// to capture the file/directory change events
	NotifyTimeout = time.Hour
)

// Scanner implements reading of messages from segmented journal files
type Scanner struct {
	offset      uint64
	journalDir  *watchedJournalDir
	journalFile *JournalFile
	file        watchedReadSeekCloser
	message     Message
	err         error

	Timeout time.Duration // read timeout when no data arrived, default 0
}
type watchedReadSeekCloser interface {
	readSeekCloser
	Watch() chan bool
}

// NewScanner creates a scanner for reading dir/jnl starting from offset
func NewScanner(dir string, offset uint64) (*Scanner, error) {
	dir = JournalDirPath(dir)
	r := Scanner{}
	journalDir, err := openWatchedJournalDir(dir)
	if err != nil {
		return nil, err
	}
	journalFile, err := journalDir.Find(offset)
	if err != nil {
		return nil, err
	}
	if journalDir.IsLast(journalFile) {
		r.file, err = openWatchedFile(journalFile.FileName)
	} else {
		r.file, err = openDummyWatchedFile(journalFile.FileName)
	}
	if err != nil {
		return nil, err
	}
	r.offset = journalFile.FirstOffset
	r.journalFile = journalFile
	r.journalDir = journalDir
	for r.offset < offset && r.Scan() {
	}
	if r.Err() != nil {
		return nil, r.Err()
	}
	if r.offset != offset {
		return nil, fmt.Errorf("fail to find offset %d", offset)
	}
	return &r, nil
}

// Scan scans the next message and increment the offset
func (r *Scanner) Scan() bool {
	if r.err != nil {
		return false
	}
	for {
		fileChanged, dirChanged := r.file.Watch(), r.journalDir.Watch()
		var n int64
		n, r.err = r.message.ReadFrom(r.file)
		if r.err != nil {
			// rollback the reader
			if _, seekErr := r.file.Seek(-n, io.SeekCurrent); seekErr != nil {
				return false
			}

			// unexpected io error
			switch r.err {
			case io.EOF, io.ErrUnexpectedEOF:
			default:
				return false
			}

			// not the last one, open the next journal file
			if !r.journalDir.IsLast(r.journalFile) {
				if r.err = r.reopenFile(); r.err != nil {
					return false
				}
				continue
			}

			// the last one, wait for any changes
			var timeoutChan <-chan time.Time
			if r.Timeout != 0 {
				timeoutChan = time.After(r.Timeout)
			}
			select {
			case <-dirChanged:
				if r.err = r.reopenFile(); r.err != nil {
					return false
				}
			case <-fileChanged:
			case <-timeoutChan:
				r.err = ErrTimeout
				return false
			case <-time.After(NotifyTimeout):
			}
			continue
		}
		break
	}

	// check offset
	if r.message.Offset != r.offset {
		r.err = fmt.Errorf("offset is out of order, expect %d but got %d", r.offset, r.message.Offset)
		return false
	}
	r.offset++

	return true
}

func (r *Scanner) Message() *Message {
	return &r.message
}

func (r *Scanner) Err() error {
	return r.err
}

func (r *Scanner) reopenFile() error {
	journalFile, err := r.journalDir.Find(r.offset)
	if err != nil {
		return err
	}
	var newFile watchedReadSeekCloser
	if r.journalDir.IsLast(journalFile) {
		newFile, err = openWatchedFile(journalFile.FileName)
	} else {
		newFile, err = openDummyWatchedFile(journalFile.FileName)
	}
	if err != nil {
		return err
	}
	if err := r.file.Close(); err != nil {
		return err
	}
	r.file = newFile
	r.journalFile = journalFile
	return nil
}

// Offset returns the current offset of the reader
func (r *Scanner) Offset() uint64 {
	return r.offset
}

// Close closes the reader
func (r *Scanner) Close() error {
	err1 := r.journalDir.Close()
	err2 := r.file.Close()
	if err1 != nil {
		return err1
	}
	return err2
}

type dummyWatchedFile struct {
	*os.File
}

func openDummyWatchedFile(file string) (*dummyWatchedFile, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	return &dummyWatchedFile{File: f}, nil
}

func (f *dummyWatchedFile) Watch() chan bool { return nil }
