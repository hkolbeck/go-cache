package relativistic

import (
	"sync"
	"sync/atomic"
)

type Reader struct {
	epoch int64
}

type writer struct {
	epoch int64
	sync.Mutex
}

type Relativistic struct {
	readers chan<- *Reader
	writers chan<- *writer
	epoch   int64
	sync.Mutex
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
				if r.epoch == 0 {
					delete(currentReaders, r)
				} else {
					currentReaders[r] = true
				}
			case w := <-writers:
				for reader := range currentReaders {
					if reader.epoch <= w.epoch {
						writers <- w
						continue ReaderLoop
					}
				}
				w.Unlock()
			}
		}
	}()

	return &Relativistic{
		readers: readers,
		writers: writers,
		epoch:   1,
	}
}

func (r *Relativistic) StartRead() *Reader {
	reader := &Reader{r.epoch}
	r.readers <- reader
	return reader
}

func (r *Relativistic) EndRead(reader *Reader) {
	reader.epoch = 0
	r.readers <- reader
}

func (r *Relativistic) WaitForReaders() {
	w := &writer{epoch: atomic.AddInt64(&(r.epoch), 1)}
	w.Lock()
	r.writers <- w
	w.Lock()
}
