package metrics

type MetricLabelSuccess struct {
	Source             string
	StoreNumber        string
	Reason             string
	IsWalmartIdBlocked *bool
}

type MetricStoreNumber struct {
	StoreNumber string
}

type MetricSource struct {
	Source string
}

type MetricApiCall struct {
	Name string
	Type string
}

type MetricLabelCancelShippingGroup struct {
	ReasonId          int
	TotalCancellation bool
	StoreNumber       int
}

type MetricEvent struct {
	Event string
}
