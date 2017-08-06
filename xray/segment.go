package xray

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"
)

// replacer is used to remove `(`, `)`, `*` from segment's name.
// only `)` is not allowed on X-Ray though.
var replacer = strings.NewReplacer("(", "",
	")", "",
	"*", "")

// Segment is span data for AWS X-Ray.
type Segment struct {
	service *XRay
	Trace   bool

	// Required
	TraceID   string
	ID        string
	Name      string
	StartTime time.Time
	EndTime   time.Time

	// Optional
	User           string
	ParentID       string
	Error          string
	Annotations    map[string]interface{}
	Request        *http.Request
	ResponseStatus int
	ContentLength  int

	// subsegments
	SQLQuery string
}

// NewSegment creates new segment with given name.
func NewSegment(name string) *Segment {
	s := &Segment{
		Name:        name,
		Annotations: make(map[string]interface{}),
	}
	s.Init()
	return s
}

// NewSegmentFromRequest creates new segment from *http.Request.
func NewSegmentFromRequest(r *http.Request) *Segment {
	s := NewSegment(r.URL.Path)
	s.Request = r
	return s
}

// NewEmptySegment creates dummy segment.
// This data is created by sampling policy and does not send to AWS API.
func NewEmptySegment() *Segment {
	return &Segment{
		Annotations: make(map[string]interface{}),
	}
}

// Init initializes segment data.
func (s *Segment) Init() {
	s.TraceID = nextTraceID()
	s.ID = nextID()
	s.StartTime = time.Now()
	s.Trace = true
}

// NewChild creates child segment data from the Segment.
func (s *Segment) NewChild(name string) *Segment {
	return &Segment{
		service:     s.service,
		Trace:       s.Trace,
		Name:        name,
		TraceID:     s.TraceID,
		ParentID:    s.ID,
		ID:          nextID(),
		StartTime:   time.Now(),
		Annotations: make(map[string]interface{}),
	}
}

// ToJSON converts to json byte data for AWS X-Ray API.
func (s *Segment) ToJSON() ([]byte, error) {
	r := &SegmentRoot{
		TraceID:   s.TraceID,
		ID:        s.ID,
		Name:      sanitizeName(s.Name),
		StartTime: toEpochTime(s.StartTime),
		EndTime:   toEpochTime(s.EndTime),
		User:      s.User,
		ParentID:  s.ParentID,
		Error:     s.Error,
	}

	if len(s.Annotations) != 0 {
		r.Annotations = s.Annotations
	}

	req := NewSegmentRequest(s.Request)
	resp := NewSegmentResponse(s.ResponseStatus, s.ContentLength)
	if req != nil || resp != nil {
		r.SegmentHTTP = &SegmentHTTP{
			SegmentRequest:  req,
			SegmentResponse: resp,
		}
	}

	if s.SQLQuery != "" {
		sub := &SegmentRoot{
			TraceID:   r.TraceID,
			ID:        nextID(),
			Name:      r.Name,
			StartTime: r.StartTime,
			EndTime:   r.EndTime,
			SegmentSQL: &SegmentSQL{
				SanitizedQuery: s.SQLQuery,
			},
		}
		r.SubSegments = append(r.SubSegments, sub)
	}
	return json.Marshal(r)
}

// sanitizeName removes some special characters to avoid X-Ray's limitation.
// X-Ray is not allowd some special characters.
func sanitizeName(s string) string {
	return replacer.Replace(s)
}

// toEpochTime returns epoch time with float point.
func toEpochTime(dt time.Time) float64 {
	if dt.IsZero() {
		return 0.0
	}
	return float64(dt.UnixNano()) / float64(time.Second)
}

// Finish ends segment timer and add segment into daemon's spool.
func (s *Segment) Finish() {
	if !s.Trace {
		return
	}

	s.EndTime = time.Now()
	if s.service != nil {
		s.service.AddSegment(s)
	}
}

// SegmentRoot is root segment data for converting JSON as a X-Ray's specification.
type SegmentRoot struct {
	// Required
	TraceID   string  `json:"trace_id,omitempty"`
	ID        string  `json:"id,omitempty"`
	Name      string  `json:"name,omitempty"`
	StartTime float64 `json:"start_time,omitempty"`
	EndTime   float64 `json:"end_time,omitempty"`

	// Optional
	User         string                 `json:"user,omitempty"`
	ParentID     string                 `json:"parent_id,omitempty"`
	Error        string                 `json:"error,omitempty"`
	Annotations  map[string]interface{} `json:"annotations,omitempty"`
	*SegmentHTTP `json:"http,omitempty"`

	// subsegments
	SubSegments []*SegmentRoot `json:"subsegments,omitempty"`
	*SegmentSQL `json:"sql,omitempty"`
}

// SegmentHTTP is segment data for http.
type SegmentHTTP struct {
	*SegmentRequest  `json:"request,omitempty"`
	*SegmentResponse `json:"response,omitempty"`
}

// SegmentRequest is segment data for http request.
type SegmentRequest struct {
	URL           string `json:"url,omitempty"`
	Method        string `json:"method,omitempty"`
	UserAgent     string `json:"user_agent,omitempty"`
	ClientIP      string `json:"client_ip,omitempty"`
	XForwardedFor bool   `json:"x_forwarded_for,omitempty"`
}

// NewSegmentRequest creates SegmentRequest data from *http.Request.
func NewSegmentRequest(r *http.Request) *SegmentRequest {
	if r == nil {
		return nil
	}

	s := &SegmentRequest{
		URL:       r.URL.RequestURI(),
		Method:    r.Method,
		UserAgent: r.UserAgent(),
	}
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		s.ClientIP = ip
		s.XForwardedFor = true
	} else {
		s.ClientIP, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return s
}

// SegmentResponse is segment data for http response.
type SegmentResponse struct {
	Status        int `json:"status,omitempty"`
	ContentLength int `json:"content_length,omitempty"`
}

// NewSegmentResponse creates SegmentResponse data from http status code and content length.
func NewSegmentResponse(status, length int) *SegmentResponse {
	if status == 0 && length == 0 {
		return nil
	}
	return &SegmentResponse{
		Status:        status,
		ContentLength: length,
	}
}

// SegmentSQL is segment data for sql.
type SegmentSQL struct {
	SanitizedQuery string `json:"sanitized_query,omitempty"`
}
