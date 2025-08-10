package paginate

type Page struct {
	Number int `json:"number"`
	Skip   int `json:"skip"`
	Rows   int `json:"rows"`
}

type PaginatedResponse[T any] struct {
	Hits    []T  `json:"hits"`
	HasMore bool `json:"hasMore"`
	Total   int  `json:"total,omitempty"`
}

type Range struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type PaginatedRequest struct {
	Filter       *Filter  `json:"filter"`
	GroupingFunc Agg      `json:"groupingFunc"`
	Sorts        []Sort   `json:"sorts"`
	Groups       []string `json:"groups"`
	Page         Page     `json:"page"`
}

type Agg struct {
	Function  string `json:"func"`
	Field     string `json:"field"`
	ResolveAs string `json:"resolveAs"`
}

type Filter struct {
	Field     string     `json:"field"`
	Type      FilterType `json:"type"`
	Values    []string   `json:"values"`
	ValueType ValueType  `json:"valueType"`
	Filters   []Filter   `json:"filters"`

	Range Range `json:"range"`
}

type FilterType string

const (
	FilterTypeEquals                 FilterType = "EQUALS"
	FilterTypeGreaterThan            FilterType = "GT"
	FilterTypeLessThan               FilterType = "LT"
	FilterTypeGreaterThanEquals      FilterType = "GTE"
	FilterTypeLessThanEquals         FilterType = "LTE"
	FilterTypeNotEquals              FilterType = "NOT_EQUALS"
	FilterTypeContains               FilterType = "CONTAINS"
	FilterTypeLike                   FilterType = "LIKE"
	FilterTypeStartsWith             FilterType = "STARTS_WITH"
	FilterTypeEndsWith               FilterType = "ENDS_WITH"
	FilterTypeIn                     FilterType = "IN"
	FilterTypeNotIn                  FilterType = "NOT_IN"
	FilterTypeNot                    FilterType = "NOT"
	FilterTypeAnd                    FilterType = "AND"
	FilterTypeOr                     FilterType = "OR"
)

type ValueType string

const (
	ValueTypeString  ValueType = "string"
	ValueTypeNumeric ValueType = "numeric"
	ValueTypeBoolean ValueType = "boolean"
)

type Sort struct {
	Field string `json:"field"`
	Order string `json:"order"`
}
