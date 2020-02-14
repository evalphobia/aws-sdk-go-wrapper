package costexplorer

import (
	SDK "github.com/aws/aws-sdk-go/service/costexplorer"
)

// UsageResult represents results from `GetCostAndUsage` operation.
type UsageResult struct {
	GroupDefinitions []GroupDefinition
	ResultsByTime    []ResultByTime
	NextPageToken    string
}

func NewUsageResult(output *SDK.GetCostAndUsageOutput) UsageResult {
	r := UsageResult{}
	if output == nil {
		return r
	}

	if output.NextPageToken != nil {
		r.NextPageToken = *output.NextPageToken
	}

	r.GroupDefinitions = newGroupDefinitions(output.GroupDefinitions)
	r.ResultsByTime = newResultsByTime(output.ResultsByTime)
	return r
}

type GroupDefinition struct {
	Key  string
	Type string
}

func newGroupDefinitions(list []*SDK.GroupDefinition) []GroupDefinition {
	if len(list) == 0 {
		return nil
	}

	results := make([]GroupDefinition, len(list))
	for i, v := range list {
		results[i] = newGroupDefinition(v)
	}
	return results
}

func newGroupDefinition(d *SDK.GroupDefinition) GroupDefinition {
	result := GroupDefinition{}
	if d == nil {
		return result
	}

	if d.Key != nil {
		result.Key = *d.Key
	}
	if d.Type != nil {
		result.Type = *d.Type
	}
	return result
}

type ResultByTime struct {
	Estimated       bool
	TimePeriodStart string
	TimePeriodEnd   string

	Groups []Group
	Total  MetricValues
}

func newResultsByTime(list []*SDK.ResultByTime) []ResultByTime {
	if len(list) == 0 {
		return nil
	}

	results := make([]ResultByTime, len(list))
	for i, v := range list {
		results[i] = newResultByTime(v)
	}
	return results
}

func newResultByTime(d *SDK.ResultByTime) ResultByTime {
	result := ResultByTime{}
	if d == nil {
		return result
	}

	if d.Estimated != nil {
		result.Estimated = *d.Estimated
	}

	if d.TimePeriod != nil {
		dt := d.TimePeriod
		if dt.Start != nil {
			result.TimePeriodStart = *dt.Start
		}
		if dt.End != nil {
			result.TimePeriodEnd = *dt.End
		}
	}

	result.Groups = newGroups(d.Groups)
	result.Total = newMetricValues(d.Total)
	return result
}

type Group struct {
	Keys []string
	MetricValues
}

func newGroups(list []*SDK.Group) []Group {
	if len(list) == 0 {
		return nil
	}

	results := make([]Group, len(list))
	for i, v := range list {
		results[i] = newGroup(v)
	}
	return results
}

func newGroup(g *SDK.Group) Group {
	result := Group{}
	if g == nil {
		return result
	}

	if len(g.Keys) != 0 {
		keys := make([]string, len(g.Keys))
		for i, v := range g.Keys {
			keys[i] = *v
		}
		result.Keys = keys
	}

	result.MetricValues = newMetricValues(g.Metrics)
	return result
}

type MetricValues struct {
	AmortizedCost         MetricValue
	NetAmortizedCost      MetricValue
	BlendedCost           MetricValue
	UnblendedCost         MetricValue
	NetUnblendedCost      MetricValue
	NormalizedUsageAmount MetricValue
	UsageQuantity         MetricValue
}

func newMetricValues(m map[string]*SDK.MetricValue) MetricValues {
	result := MetricValues{}
	if m == nil {
		return result
	}

	if v, ok := m["AmortizedCost"]; ok {
		result.AmortizedCost = newMetricValue(v)
	}
	if v, ok := m["NetAmortizedCost"]; ok {
		result.NetAmortizedCost = newMetricValue(v)
	}
	if v, ok := m["BlendedCost"]; ok {
		result.BlendedCost = newMetricValue(v)
	}
	if v, ok := m["UnblendedCost"]; ok {
		result.UnblendedCost = newMetricValue(v)
	}
	if v, ok := m["NetUnblendedCost"]; ok {
		result.NetUnblendedCost = newMetricValue(v)
	}
	if v, ok := m["NormalizedUsageAmount"]; ok {
		result.NormalizedUsageAmount = newMetricValue(v)
	}
	if v, ok := m["UsageQuantity"]; ok {
		result.UsageQuantity = newMetricValue(v)
	}
	return result
}

func (v MetricValues) GetOne() (amount, unit string) {
	switch {
	case v.UnblendedCost.Unit != "":
		return v.UnblendedCost.Amount, v.UnblendedCost.Unit
	case v.BlendedCost.Unit != "":
		return v.BlendedCost.Amount, v.BlendedCost.Unit
	case v.UsageQuantity.Unit != "":
		return v.UsageQuantity.Amount, v.UsageQuantity.Unit
	case v.NetUnblendedCost.Unit != "":
		return v.NetUnblendedCost.Amount, v.NetUnblendedCost.Unit
	case v.AmortizedCost.Unit != "":
		return v.AmortizedCost.Amount, v.AmortizedCost.Unit
	case v.NetAmortizedCost.Unit != "":
		return v.NetAmortizedCost.Amount, v.NetAmortizedCost.Unit
	case v.NormalizedUsageAmount.Unit != "":
		return v.NormalizedUsageAmount.Amount, v.NormalizedUsageAmount.Unit
	}
	return "", ""
}

type MetricValue struct {
	Amount string
	Unit   string
}

func newMetricValue(v *SDK.MetricValue) MetricValue {
	result := MetricValue{}
	if v == nil {
		return result
	}

	if v.Amount != nil {
		result.Amount = *v.Amount
	}
	if v.Unit != nil {
		result.Unit = *v.Unit
	}
	return result
}
