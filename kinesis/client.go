package kinesis

import (
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	SDK "github.com/aws/aws-sdk-go/service/kinesis"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

const (
	serviceName = "Kinesis"
)

// Kinesis has Kinesis client.
type Kinesis struct {
	client *SDK.Kinesis

	logger log.Logger
	prefix string

	streamsMu sync.RWMutex
	streams   map[string]*Stream
}

// New returns initialized *Kinesis.
func New(conf config.Config) (*Kinesis, error) {
	sess, err := conf.Session()
	if err != nil {
		return nil, err
	}

	svc := NewFromSession(sess)
	svc.prefix = conf.DefaultPrefix
	return svc, nil
}

// NewFromSession returns initialized *Kinesis from aws.Session.
func NewFromSession(sess *session.Session) *Kinesis {
	return &Kinesis{
		client:  SDK.New(sess),
		logger:  log.DefaultLogger,
		streams: make(map[string]*Stream),
	}
}

// GetClient gets aws client.
func (svc *Kinesis) GetClient() *SDK.Kinesis {
	return svc.client
}

// SetLogger sets logger.
func (svc *Kinesis) SetLogger(logger log.Logger) {
	svc.logger = logger
}

// SetPrefix sets prefix.
func (svc *Kinesis) SetPrefix(prefix string) {
	svc.prefix = prefix
}

// GetStream gets Kinesis Stream.
func (svc *Kinesis) GetStream(name string) (*Stream, error) {
	streamName := svc.prefix + name

	// get the stream from cache
	svc.streamsMu.RLock()
	s, ok := svc.streams[streamName]
	svc.streamsMu.RUnlock()
	if ok {
		return s, nil
	}

	s, err := NewStream(svc, name)
	if err != nil {
		return nil, err
	}

	svc.streamsMu.Lock()
	svc.streams[streamName] = s
	svc.streamsMu.Unlock()
	return s, nil
}

// CreateStream creates new Kinesis Stream.
func (svc *Kinesis) CreateStream(in *SDK.CreateStreamInput) error {
	_, err := svc.client.CreateStream(in)
	if err != nil {
		svc.Errorf("error on `CreateStream` operation; stream=%s; error=%s;", *in.StreamName, err.Error())
		return err
	}

	svc.Infof("success on `CreateStream` operation; stream=%s;", *in.StreamName)
	return nil
}

// CreateStreamWithName creates new Kinesis Stream by given name with prefix.
func (svc *Kinesis) CreateStreamWithName(name string) error {
	streamName := svc.prefix + name
	return svc.CreateStream(&SDK.CreateStreamInput{
		StreamName: pointers.String(streamName),
		ShardCount: pointers.Long(1),
	})
}

// IsExistStream checks if the Stream already exists or not.
func (svc *Kinesis) IsExistStream(name string) (bool, error) {
	streamName := svc.prefix + name
	desc, err := svc.client.DescribeStream(&SDK.DescribeStreamInput{
		StreamName: pointers.String(streamName),
	})

	switch {
	case isNonExistentStreamError(err):
		return false, nil
	case err != nil:
		svc.Errorf("error on `DescribeStream` operation; stream=%s; error=%s", streamName, err.Error())
		return false, err
	case desc == nil:
		return false, nil
	case desc.StreamDescription != nil: // stream exists
		return true, nil
	default:
		return false, nil
	}
}

func isNonExistentStreamError(err error) bool {
	const errNonExistentStream = "ResourceNotFoundException: "
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), errNonExistentStream)
}

// ForceDeleteStream deletes Kinesis stream by given name with prefix.
func (svc *Kinesis) ForceDeleteStream(name string) error {
	streamName := svc.prefix + name
	_, err := svc.client.DeleteStream(&SDK.DeleteStreamInput{
		StreamName: pointers.String(streamName),
	})
	if err != nil {
		svc.Errorf("error on `DeleteStream` operation; stream=%s; error=%s;", streamName, err.Error())
		return err
	}

	svc.Infof("success on `DeleteStream` operation; stream=%s;", streamName)
	return nil
}

// Infof logging information.
func (svc *Kinesis) Infof(format string, v ...interface{}) {
	svc.logger.Infof(serviceName, format, v...)
}

// Errorf logging error information.
func (svc *Kinesis) Errorf(format string, v ...interface{}) {
	svc.logger.Errorf(serviceName, format, v...)
}
