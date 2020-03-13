package cloudwatch

import (
	"time"

	SDK "github.com/aws/aws-sdk-go/service/cloudwatch"

	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

type MetricStatisticsInput struct {
	StartTime          time.Time
	EndTime            time.Time
	Period             int64
	MetricName         string
	Namespace          string
	Unit               string
	Statistics         []string
	ExtendedStatistics []string

	Dimensions []Dimension
	// Key: Dimension.Name, Value: Dimension.Value.
	// If you use same key and differenct values, then use Dimensions.
	DimensionsMap map[string]string
}

func (o MetricStatisticsInput) ToInput() *SDK.GetMetricStatisticsInput {
	in := &SDK.GetMetricStatisticsInput{}
	if !o.StartTime.IsZero() {
		in.StartTime = &o.StartTime
	}
	if !o.EndTime.IsZero() {
		in.EndTime = &o.EndTime
	}
	if o.Period != 0 {
		in.Period = &o.Period
	}
	if o.MetricName != "" {
		in.MetricName = &o.MetricName
	}
	if o.Namespace != "" {
		in.Namespace = &o.Namespace
	}
	if o.Unit != "" {
		in.Unit = &o.Unit
	}
	in.Statistics = pointers.SliceString(o.Statistics)
	in.ExtendedStatistics = pointers.SliceString(o.ExtendedStatistics)

	in.Dimensions = make([]*SDK.Dimension, 0, len(o.DimensionsMap)+len(o.Dimensions))
	for key, val := range o.DimensionsMap {
		in.Dimensions = append(in.Dimensions, &SDK.Dimension{
			Name:  pointers.String(key),
			Value: pointers.String(val),
		})
	}
	for _, v := range o.Dimensions {
		in.Dimensions = append(in.Dimensions, &SDK.Dimension{
			Name:  pointers.String(v.Name),
			Value: pointers.String(v.Value),
		})
	}
	return in
}

type Dimension struct {
	Name  string
	Value string
}

type PutMetricDataInput struct {
	MetricData []MetricDatum
	Namespace  string
}

func (o *PutMetricDataInput) AddMetric(d MetricDatum) {
	o.MetricData = append(o.MetricData, d)
}

func (o PutMetricDataInput) ToInput() *SDK.PutMetricDataInput {
	in := &SDK.PutMetricDataInput{}
	if o.Namespace != "" {
		in.Namespace = &o.Namespace
	}

	in.MetricData = make([]*SDK.MetricDatum, 0, len(o.MetricData))
	for _, v := range o.MetricData {
		in.MetricData = append(in.MetricData, v.ToSDKValue())
	}
	return in
}

type MetricDatum struct {
	MetricName        string
	Unit              string
	StorageResolution int64
	Value             float64
	HasValue          bool // use as true when value == 0
	Values            []float64
	Counts            []float64
	Timestamp         time.Time

	StatisticValues StatisticSet
	Dimensions      []Dimension
}

func (d MetricDatum) ToSDKValue() *SDK.MetricDatum {
	in := &SDK.MetricDatum{}
	if d.MetricName != "" {
		in.MetricName = &d.MetricName
	}
	if d.Unit != "" {
		in.Unit = &d.Unit
	}
	if d.StorageResolution != 0 {
		in.StorageResolution = &d.StorageResolution
	}
	if d.HasValue || d.Value != 0 {
		in.Value = &d.Value
	}
	in.Values = pointers.SliceFloat64(d.Values)
	in.Counts = pointers.SliceFloat64(d.Counts)
	if !d.Timestamp.IsZero() {
		in.Timestamp = &d.Timestamp
	}

	in.StatisticValues = d.StatisticValues.ToSDKValue()
	for _, v := range d.Dimensions {
		in.Dimensions = append(in.Dimensions, &SDK.Dimension{
			Name:  pointers.String(v.Name),
			Value: pointers.String(v.Value),
		})
	}
	return in
}

type StatisticSet struct {
	Maximum     float64
	Minimum     float64
	SampleCount float64
	Sum         float64

	// use as true when value == 0
	HasMaximum     bool
	HasMinimum     bool
	HasSampleCount bool
	HasSum         bool
}

func (d StatisticSet) ToSDKValue() *SDK.StatisticSet {
	hasValue := false

	in := &SDK.StatisticSet{}
	if d.HasMaximum || d.Maximum != 0 {
		in.Maximum = &d.Maximum
		hasValue = true
	}
	if d.HasMinimum || d.Minimum != 0 {
		in.Minimum = &d.Minimum
		hasValue = true
	}
	if d.HasSampleCount || d.SampleCount != 0 {
		in.SampleCount = &d.SampleCount
		hasValue = true
	}
	if d.HasSum || d.Sum != 0 {
		in.Sum = &d.Sum
		hasValue = true
	}

	if !hasValue {
		return nil
	}
	return in
}
