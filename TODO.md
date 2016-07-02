TODO
====

### Reader

* Search offset (done)
* Read segment files (done)
* monitoring
    - dir (done)
    - append (done)
    - append or dir (done)
* handle truncation of the last message
* Offset persistence

### Writer

* Write from the last offset (done)
* Segmentation (done)
* Lock to prevent other writer (done)
* startup corruption detection (done)
* startup corruption correction

### Command line tool

* range
* dump
* tail
* count
* rollback
* clean
	- delete files according to cleaning rules

### Benchmark

* Writer
    - Sync
    - Async
* Reader
* Single writer and multiple readers
