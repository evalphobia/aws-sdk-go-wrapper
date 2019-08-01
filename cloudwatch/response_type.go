package cloudwatch

import (
	"time"

	SDK "github.com/aws/aws-sdk-go/service/cloudwatch"
)

type MetricStatisticsResponse struct {
	Datapoints []Datapoint
	Label      string
}

func NewMetricStatisticsResponse(out *SDK.GetMetricStatisticsOutput) *MetricStatisticsResponse {
	r := &MetricStatisticsResponse{}
	if out.Label != nil {
		r.Label = *out.Label
	}
	if len(out.Datapoints) == 0 {
		return r
	}

	r.Datapoints = make([]Datapoint, len(out.Datapoints))
	for i, d := range out.Datapoints {
		r.Datapoints[i] = NewDatapoint(d)
	}
	return r
}

type Datapoint struct {
	Average            float64
	ExtendedStatistics map[string]float64
	Maximum            float64
	Minimum            float64
	SampleCount        float64
	Sum                float64
	Unit               string
	Timestamp          time.Time
}

func NewDatapoint(d *SDK.Datapoint) Datapoint {
	dd := Datapoint{}
	dd.Average = getFloat64FromPointer(d.Average)
	dd.Maximum = getFloat64FromPointer(d.Maximum)
	dd.Minimum = getFloat64FromPointer(d.Minimum)
	dd.SampleCount = getFloat64FromPointer(d.SampleCount)
	if d.Sum != nil {
		dd.Sum = *d.Sum
	}
	if d.Unit != nil {
		dd.Unit = *d.Unit
	}
	if d.Timestamp != nil {
		dd.Timestamp = *d.Timestamp
	}

	s := make(map[string]float64, len(d.ExtendedStatistics))
	for key, val := range d.ExtendedStatistics {
		s[key] = getFloat64FromPointer(val)
	}
	dd.ExtendedStatistics = s
	return dd
}

func getFloat64FromPointer(v *float64) float64 {
	if v != nil {
		return *v
	}
	return 0.0
}
