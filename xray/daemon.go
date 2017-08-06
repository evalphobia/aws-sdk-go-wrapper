package xray

import (
	"sync"
	"time"
)

// Daemon is background daemon for sending segments.
// This struct stores segments and sends segment to AWS X-Ray in each checkpoint timing.
type Daemon struct {
	flushSegments func([]*Segment) error

	spoolMu            sync.Mutex
	spool              []*Segment
	checkpointSize     int
	checkpointInterval time.Duration
	stopSignal         chan struct{}
}

// NewDaemon creates new Daemon.
// size is number of segments to send AWS API in single checkpoint.
// interval is the time of checkpoint interval.
// fn is function called in each checkpoint, to sends segments to AWS API.
func NewDaemon(size int, interval time.Duration, fn func([]*Segment) error) *Daemon {
	if size < 1 {
		size = 10
	}
	if interval == 0 {
		interval = 1 * time.Second
	}

	return &Daemon{
		spool:              make([]*Segment, 0, 4096),
		checkpointSize:     size,
		checkpointInterval: interval,
		stopSignal:         make(chan struct{}),
		flushSegments:      fn,
	}
}

// Add adds segment data into daemon.
func (d *Daemon) Add(segments ...*Segment) {
	d.spoolMu.Lock()
	d.spool = append(d.spool, segments...)
	d.spoolMu.Unlock()
}

// Flush gets segments from the internal spool and execute flushSegments function.
func (d *Daemon) Flush() {
	d.spoolMu.Lock()
	var segments []*Segment
	segments, d.spool = shiftSegment(d.spool, d.checkpointSize)
	d.spoolMu.Unlock()
	d.flushSegments(segments)
}

// shiftSegment retrieves segments.
func shiftSegment(slice []*Segment, size int) (part []*Segment, all []*Segment) {
	l := len(slice)
	if l <= size {
		return slice, slice[:0]
	}
	return slice[:size], slice[size:]
}

// Run sets timer to flush data in each checkpoint as a background daemon.
func (d *Daemon) Run() {
	ticker := time.NewTicker(d.checkpointInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				d.Flush()
			case <-d.stopSignal:
				ticker.Stop()
				return
			}
		}
	}()
}

// Stop stops daemon.
func (d *Daemon) Stop() {
	d.stopSignal <- struct{}{}
}
