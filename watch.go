package sej

import (
	"io"
	"os"
	"sync"

	"gopkg.in/fsnotify.v1"
	"h12.me/errors"
)

type watchedJournalDir struct {
	dir     *JournalDir
	watcher *changeWatcher
}

func openWatchedJournalDir(dir string) (*watchedJournalDir, error) {
	dirFile, err := openOrCreateDir(dir)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if err := dirFile.Close(); err != nil {
		return nil, errors.Wrap(err)
	}
	watcher, err := newChangeWatcher(dir, fsnotify.Create|fsnotify.Remove)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	journalDir, err := OpenJournalDir(dir)
	if err != nil {
		watcher.Close()
		return nil, errors.Wrap(err)
	}
	return &watchedJournalDir{
		dir:     journalDir,
		watcher: watcher,
	}, nil
}

func (d *watchedJournalDir) Watch() chan bool {
	return d.watcher.Watch()
}

func (d *watchedJournalDir) Find(offset uint64) (*JournalFile, error) {
	if err := d.watcher.Err(); err != nil {
		return nil, err
	}
	if err := d.reload(); err != nil {
		return nil, err
	}
	return d.dir.find(offset)
}

func (d *watchedJournalDir) IsLast(f *JournalFile) bool {
	if !d.dir.isLast(f) {
		return false
	}
	d.reload()
	return d.dir.isLast(f)
}
func (d *watchedJournalDir) reload() error {
	journalDir, err := OpenJournalDir(d.dir.path)
	if err != nil {
		return err
	}
	d.dir = journalDir
	return nil
}

func (d *watchedJournalDir) Close() error {
	return d.watcher.Close()
}

type watchedFile struct {
	file    *os.File
	watcher *changeWatcher
}

func openWatchedFile(name string) (*watchedFile, error) {
	watcher, err := newChangeWatcher(name, fsnotify.Write)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(name)
	if err != nil {
		watcher.Close()
		return nil, err
	}
	return &watchedFile{
		file:    file,
		watcher: watcher,
	}, nil
}

func (f *watchedFile) Watch() chan bool {
	return f.watcher.Watch()
}

func (f *watchedFile) Seek(offset int64, whence int) (int64, error) {
	return f.file.Seek(offset, whence)
}

func (f *watchedFile) Read(p []byte) (n int, err error) {
	if err := f.watcher.Err(); err != nil {
		return 0, err
	}
	n, err = f.file.Read(p)
	if err == io.EOF {
		if err := f.reopen(); err != nil {
			return n, err
		}
		return f.file.Read(p)
	}
	return n, err
}
func (f *watchedFile) reopen() error {
	oldOffset, err := f.file.Seek(0, os.SEEK_CUR)
	if err != nil {
		return err
	}
	newFile, err := os.Open(f.file.Name())
	if err != nil {
		return err
	}
	if _, err := newFile.Seek(oldOffset, os.SEEK_SET); err != nil {
		newFile.Close()
		return err
	}
	if err := f.file.Close(); err != nil {
		return err
	}
	f.file = newFile
	return nil
}

func (f *watchedFile) Close() error {
	err1 := f.file.Close()
	err2 := f.watcher.Close()
	if err1 != nil {
		return err1
	}
	return err2
}

// changeWatcher compresses multiple change messages into one
type changeWatcher struct {
	watcher   *fsnotify.Watcher
	watchedOp fsnotify.Op
	changedCh chan bool
	wg        sync.WaitGroup

	err error
	mu  sync.RWMutex
}

func newChangeWatcher(name string, op fsnotify.Op) (*changeWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if err := watcher.Add(name); err != nil {
		watcher.Close()
		return nil, errors.Wrap(err)
	}
	w := &changeWatcher{
		watcher:   watcher,
		watchedOp: op,
		changedCh: make(chan bool, 1), // make sure at least one message can be received when needed
	}
	w.wg.Add(2)
	go w.watchEvent()
	go w.watchError()
	return w, nil
}

// Watch returns an empty channel for receiving a single event after the method is called
func (w *changeWatcher) Watch() chan bool {
clearChan:
	for {
		select {
		case <-w.changedCh:
		default:
			break clearChan
		}
	}
	return w.changedCh
}

func (w *changeWatcher) watchEvent() {
	defer w.wg.Done()
	for event := range w.watcher.Events {
		if event.Op&w.watchedOp > 0 {
			select {
			case w.changedCh <- true: // send at least one
			default: // or skip the rest
			}
		}
	}
}

func (w *changeWatcher) watchError() {
	defer w.wg.Done()
	for err := range w.watcher.Errors {
		w.mu.Lock()
		w.err = err
		w.mu.Unlock()
	}
}

func (w *changeWatcher) Err() error {
	w.mu.RLock()
	err := w.err
	w.mu.RUnlock()
	return err
}

func (w *changeWatcher) Close() error {
	w.watcher.Close()
	w.wg.Wait()
	return nil
}
