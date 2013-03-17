package relativistic

import "sync"

type Reader struct {
	epoch int64
}

type Writer struct {
	lock  *sync.Mutex
	epoch int64
}

type Relativistic struct {
	readers chan<- *Reader
	writers chan<- *sync.Mutex
	epoch   int64
	lock    *sync.Mutex
}

func New(queueSize int) *Relativistic {
	if queueSize < 0 {
		panic("Negative queue size")
	}

	readers := make(chan *Reader, queueSize)
	writers := make(chan *writer, 1)

	go func() {
		currentReaders := make(map[*Reader]bool)

	ReaderLoop:
		for {
			select {
			case r := <-readers:
				if *r == 0 {
					delete(currentReaders, r)
				} else {
					currentReaders[r] = true
				}
			case w := <-writers:
				for reader := range currentReaders {
					if reader.epoch <= writer.epoch {
						writers <- w
						continue ReaderLoop
					}
				}
				w.lock.Unlock()
			}
		}
	}()

	return &Relativistic{
		readers: readers,
		writers: writers,
		epoch:   1,
		lock:    *sync.Mutex{},
	}
}

func (r *Relativistic) startRead() *Reader {
	reader := &reader{r.epoch}
	r.readers <- reader
	return reader
}

func (r *Relativistic) endRead(reader *Reader) {
	reader.epoch = 0
	r.readers <- reader
}

func (r *Relativistic) startWrite() *Writer {
	lock.Lock()
	return &Writer{&sync.Mutex, r.epoch}
}

func (r *Relativistic) waitForReaders(writer *Writer) {
	writer.lock.Lock()
	r.writers <- writer
	writer.lock.Lock()
}

func (r *Relativistic) endWrite() {
	lock.Unlock()
}
