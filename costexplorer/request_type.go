package costexplorer

import (
	"time"

	SDK "github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

const (
	GranularityMonthly = SDK.GranularityMonthly
	GranularityDaily   = SDK.GranularityDaily
	GranularityHourly  = SDK.GranularityHourly

	GroupByDimension = SDK.GroupDefinitionTypeDimension
	GroupByTag       = SDK.GroupDefinitionTypeTag

	MetricBlendedCost           = SDK.MetricBlendedCost
	MetricUnblendedCost         = SDK.MetricUnblendedCost
	MetricAmortizedCost         = SDK.MetricAmortizedCost
	MetricNetUnblendedCost      = SDK.MetricNetUnblendedCost
	MetricNetAmortizedCost      = SDK.MetricNetAmortizedCost
	MetricUsageQuantity         = SDK.MetricUsageQuantity
	MetricNormalizedUsageAmount = SDK.MetricNormalizedUsageAmount
)

// GetCostAndUsageInput is optional parameters for `GetCostAndUsage`.
type GetCostAndUsageInput struct {
	NextPageToken   string
	TimePeriodStart time.Time
	TimePeriodEnd   time.Time

	GranularityMonthly bool
	GranularityDaily   bool
	GranularityHourly  bool

	GroupByDimensionAZ              bool
	GroupByDimensionInstanceType    bool
	GroupByDimensionLegalEntityName bool
	GroupByDimensionLinkedAccount   bool
	GroupByDimensionOperation       bool
	GroupByDimensionPlatform        bool
	GroupByDimensionPurchaseType    bool
	GroupByDimensionService         bool
	GroupByDimensionTenancy         bool
	GroupByDimensionRecordType      bool
	GroupByDimensionUsageType       bool
	GroupByTagKeys                  []string

	MetricAmortizedCost         bool
	MetricBlendedCost           bool
	MetricNetAmortizedCost      bool
	MetricNetUnblendedCost      bool
	MetricNormalizedUsageAmount bool
	MetricUnblendedCost         bool
	MetricUsageQuantity         bool

	Filter *SDK.Expression
}

// ToInput converts to *SDK.GetCostAndUsageInput.
func (u GetCostAndUsageInput) ToInput() *SDK.GetCostAndUsageInput {
	in := &SDK.GetCostAndUsageInput{}

	// set NextPageToken
	if u.NextPageToken != "" {
		in.NextPageToken = pointers.String(u.NextPageToken)
	}

	// set TimePeriod
	if u.TimePeriodEnd.IsZero() {
		u.TimePeriodEnd = time.Now().AddDate(0, 0, -1)
	}
	if u.TimePeriodStart.IsZero() {
		u.TimePeriodStart = u.TimePeriodEnd.AddDate(0, 0, -1)
	}

	in.TimePeriod = &SDK.DateInterval{
		Start: pointers.String(u.TimePeriodStart.Format("2006-01-02")),
		End:   pointers.String(u.TimePeriodEnd.Format("2006-01-02")),
	}

	// set Granularity
	switch {
	case u.GranularityDaily:
		in.Granularity = pointers.String(GranularityDaily)
	case u.GranularityMonthly:
		in.Granularity = pointers.String(GranularityMonthly)
	case u.GranularityHourly:
		in.Granularity = pointers.String(GranularityHourly)
	default:
		in.Granularity = pointers.String(GranularityDaily)
	}

	// set GroupBy
	if u.GroupByDimensionAZ {
		in.GroupBy = append(in.GroupBy, newGroupDefinitionDimension(SDK.DimensionAz))
	}
	if u.GroupByDimensionInstanceType {
		in.GroupBy = append(in.GroupBy, newGroupDefinitionDimension(SDK.DimensionInstanceType))
	}
	if u.GroupByDimensionLinkedAccount {
		in.GroupBy = append(in.GroupBy, newGroupDefinitionDimension(SDK.DimensionLinkedAccount))
	}
	if u.GroupByDimensionOperation {
		in.GroupBy = append(in.GroupBy, newGroupDefinitionDimension(SDK.DimensionOperation))
	}
	if u.GroupByDimensionPurchaseType {
		in.GroupBy = append(in.GroupBy, newGroupDefinitionDimension(SDK.DimensionPurchaseType))
	}
	if u.GroupByDimensionService {
		in.GroupBy = append(in.GroupBy, newGroupDefinitionDimension(SDK.DimensionService))
	}
	if u.GroupByDimensionTenancy {
		in.GroupBy = append(in.GroupBy, newGroupDefinitionDimension(SDK.DimensionTenancy))
	}
	if u.GroupByDimensionRecordType {
		in.GroupBy = append(in.GroupBy, newGroupDefinitionDimension(SDK.DimensionRecordType))
	}
	if u.GroupByDimensionUsageType {
		in.GroupBy = append(in.GroupBy, newGroupDefinitionDimension(SDK.DimensionUsageType))
	}
	for _, v := range u.GroupByTagKeys {
		in.GroupBy = append(in.GroupBy, newGroupDefinitionTag(v))
	}

	// set Metrics
	if u.MetricAmortizedCost {
		in.Metrics = append(in.Metrics, pointers.String(SDK.MetricAmortizedCost))
	}
	if u.MetricBlendedCost {
		in.Metrics = append(in.Metrics, pointers.String(SDK.MetricBlendedCost))
	}
	if u.MetricNetAmortizedCost {
		in.Metrics = append(in.Metrics, pointers.String(SDK.MetricNetAmortizedCost))
	}
	if u.MetricNetUnblendedCost {
		in.Metrics = append(in.Metrics, pointers.String(SDK.MetricNetUnblendedCost))
	}
	if u.MetricNormalizedUsageAmount {
		in.Metrics = append(in.Metrics, pointers.String(SDK.MetricNormalizedUsageAmount))
	}
	if u.MetricUnblendedCost {
		in.Metrics = append(in.Metrics, pointers.String(SDK.MetricUnblendedCost))
	}
	if u.MetricUsageQuantity {
		in.Metrics = append(in.Metrics, pointers.String(SDK.MetricUsageQuantity))
	}
	if len(in.Metrics) == 0 {
		in.Metrics = append(in.Metrics, pointers.String(SDK.MetricUnblendedCost))
	}

	in.Filter = u.Filter
	return in
}

func newGroupDefinitionDimension(key string) *SDK.GroupDefinition {
	return &SDK.GroupDefinition{
		Type: pointers.String(GroupByDimension),
		Key:  pointers.String(key),
	}
}

func newGroupDefinitionTag(key string) *SDK.GroupDefinition {
	return &SDK.GroupDefinition{
		Type: pointers.String(GroupByTag),
		Key:  pointers.String(key),
	}
}
