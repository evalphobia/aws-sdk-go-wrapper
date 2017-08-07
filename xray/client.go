package xray

import (
	"net/http"
	"time"

	SDK "github.com/aws/aws-sdk-go/service/xray"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
	"github.com/evalphobia/aws-sdk-go-wrapper/private/errors"
	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

const (
	serviceName = "X-Ray"
)

// sampling all request and set sampling limit to 1000 req/s.
var defaultSamplingPolicy, _ = NewLimitedSampler(1, 1000)

// XRay has XRay client.
type XRay struct {
	client   *SDK.XRay
	daemon   *Daemon
	sampling SamplingPolicy

	logger log.Logger
	prefix string
}

// New returns initialized *Kinesis.
func New(conf config.Config) (*XRay, error) {
	sess, err := conf.Session()
	if err != nil {
		return nil, err
	}

	svc := &XRay{
		client:   SDK.New(sess),
		logger:   log.DefaultLogger,
		prefix:   conf.DefaultPrefix,
		sampling: defaultSamplingPolicy,
	}
	return svc, nil
}

// SetLogger sets logger.
func (svc *XRay) SetLogger(logger log.Logger) {
	svc.logger = logger
}

// SetSamplingPolicy sets sampling policy.
func (svc *XRay) SetSamplingPolicy(fraction, qps float64) error {
	s, err := NewLimitedSampler(fraction, qps)
	if err != nil {
		svc.Errorf("error on SetSamplingPolicy; fraction=%f; qps=%f; error=%s;", fraction, qps, err.Error())
		return err
	}
	svc.sampling = s
	return nil
}

// AddSegment adds the segment dat into background daemon.
func (svc *XRay) AddSegment(segments ...*Segment) {
	svc.daemon.Add(segments...)
}

// RunDaemon creates and runs background daemon.
func (svc *XRay) RunDaemon(size int, interval time.Duration) {
	svc.daemon = NewDaemon(size, interval, svc.PutTraceSegments)
	svc.daemon.Run()
}

// PutTraceSegments executes PutTraceSegments operation.
func (svc *XRay) PutTraceSegments(segments []*Segment) error {
	if len(segments) == 0 {
		return nil
	}

	list := make([]*string, len(segments))
	for i, s := range segments {
		if !s.Trace {
			continue
		}

		byt, err := s.ToJSON()
		if err != nil {
			svc.Errorf("error on segment.ToJSON(); segment=%+v; error=%s;", s, err.Error())
			continue
		}
		list[i] = pointers.String(string(byt))
	}

	notProcessed, err := svc.client.PutTraceSegments(&SDK.PutTraceSegmentsInput{
		TraceSegmentDocuments: list,
	})
	if err != nil {
		_list := make([]string, len(list))
		for i, s := range list {
			_list[i] = *s
		}
		svc.Errorf("error on `PutTraceSegments` operation; segments=%v; error=%s;", _list, err.Error())
	}
	_ = notProcessed // TODO
	return err
}

// NewSegment creates new Segment data with given name.
func (svc *XRay) NewSegment(name string) *Segment {
	s := NewSegment(name)
	s.service = svc
	return s
}

// NewSegmentFromRequest creates new Segment data from *http.Request.
func (svc *XRay) NewSegmentFromRequest(r *http.Request) *Segment {
	if !svc.sampling.CanSample() {
		return NewEmptySegment()
	}

	s := NewSegmentFromRequest(r)
	s.service = svc
	return s
}

// Infof logging information.
func (svc *XRay) Infof(format string, v ...interface{}) {
	svc.logger.Infof(serviceName, format, v...)
}

// Errorf logging error information.
func (svc *XRay) Errorf(format string, v ...interface{}) {
	svc.logger.Errorf(serviceName, format, v...)
}

func newErrors() *errors.Errors {
	return errors.NewErrors(serviceName)
}
