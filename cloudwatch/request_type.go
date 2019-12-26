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
	in.Statistics = sliceStringToPointer(o.Statistics)
	in.ExtendedStatistics = sliceStringToPointer(o.ExtendedStatistics)

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

func sliceStringToPointer(list []string) []*string {
	if len(list) == 0 {
		return nil
	}

	result := make([]*string, len(list))
	for i, v := range list {
		result[i] = pointers.String(v)
	}
	return result
}
