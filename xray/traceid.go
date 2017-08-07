package xray

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"sync/atomic"
	"time"
)

var (
	uniqueIDCounter   uint64
	uniqueIDIncrement uint64
)

func init() {
	// copy from https://github.com/GoogleCloudPlatform/google-cloud-go/blob/master/trace/trace.go
	binary.Read(rand.Reader, binary.LittleEndian, &uniqueIDCounter)
	binary.Read(rand.Reader, binary.LittleEndian, &uniqueIDIncrement)
	uniqueIDIncrement |= 1
}

func nextID() string {
	id1 := nextUniqueID()
	return fmt.Sprintf("%016x", id1)
}

func nextTraceID() string {
	id := nextUniqueID()
	return fmt.Sprintf("1-%x-%024x", time.Now().Unix(), id)
}

func nextUniqueID() uint64 {
	var id uint64
	for id == 0 {
		id = atomic.AddUint64(&uniqueIDCounter, uniqueIDIncrement)
	}
	return id
}
