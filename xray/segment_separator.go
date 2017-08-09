package xray

import (
	"fmt"

	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

// segmentDivider is used to divide segments slice into multiple slices.
// PutTraceSegments API has 64kb limitation for segments.
// see: http://docs.aws.amazon.com/xray/latest/devguide/xray-api-segmentdocuments.html
type segmentDivider struct {
	totalSize    int
	list         [][]*string
	currentIndex int
	currentByte  int
}

func createSegmentDivider(segments []*Segment) (*segmentDivider, error) {
	size := len(segments)
	sep := &segmentDivider{
		list:      [][]*string{make([]*string, 0, size)},
		totalSize: size,
	}

	errList := newErrors()
	for _, s := range segments {
		if !s.Trace {
			continue
		}

		byt, err := s.ToJSON()
		if err != nil {
			errList.Add(fmt.Errorf("error on segment.ToJSON(); segment=%+v; error=%s;", s, err.Error()))
			continue
		}
		sep.append(byt)
	}
	if errList.HasError() {
		return sep, errList
	}
	return sep, nil
}

func (s *segmentDivider) GetResult() [][]*string {
	return s.list
}

func (s *segmentDivider) isOverByte() bool {
	const maxByte = 61440 // 60kb
	return s.currentByte > maxByte
}

func (s *segmentDivider) append(byt []byte) {
	bytesize := len(byt)
	s.currentByte += bytesize
	if s.isOverByte() {
		s.currentIndex++
		s.currentByte = bytesize
		s.list[s.currentIndex] = make([]*string, 0, s.totalSize)

	}
	s.list[s.currentIndex] = append(s.list[s.currentIndex], pointers.String(string(byt)))
}
