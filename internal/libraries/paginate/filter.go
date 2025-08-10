package paginate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	LIKE       = "LIKE"
	IN         = "IN"
	NOTIN      = "NOT IN"
	ARRAYAND   = "&&"
	BETWEEN    = "BETWEEN"
	NOTBETWEEN = "NOT BETWEEN"
	timeFormat = "2006-01-02"
)

//nolint:gocyclo // this is filters function there will be different cases to handle
func GenerateFilterClause(
	filter *Filter,
	valueFormatter ...ValueFormatter,
) string {
	if filter == nil {
		return ""
	}

	if filter.Type == FilterTypeAnd || filter.Type == FilterTypeOr ||
		filter.Type == FilterTypeNot {
		if len(filter.Filters) == 0 {
			return ""
		}

		nestedClauses := make([]string, len(filter.Filters))

		for i, nestedFilter := range filter.Filters {
			nestedClauses[i] = GenerateFilterClause(&nestedFilter, valueFormatter...)
		}

		separator := " AND " // default to AND
		if filter.Type == FilterTypeOr {
			separator = " OR "
		}

		nestedClause := fmt.Sprintf(
			"(%s)",
			strings.Join(nestedClauses, separator),
		)

		if filter.Type == FilterTypeNot {
			nestedClause = "NOT " + nestedClause
		}

		return nestedClause
	}

	var (
		operator string
		values   string
		from     string
		to       string
	)

	switch filter.Type {
	case FilterTypeEquals:
		operator = "="
	case FilterTypeNotEquals:
		operator = "!="
	case FilterTypeGreaterThan:
		operator = ">"
	case FilterTypeLessThan:
		operator = "<"
	case FilterTypeGreaterThanEquals:
		operator = ">="
	case FilterTypeLessThanEquals:
		operator = "<="
	case FilterTypeContains:
		operator = LIKE
		filter.Values[0] = "%" + filter.Values[0] + "%"
	case FilterTypeLike:
		operator = LIKE
	case FilterTypeStartsWith:
		operator = LIKE
		filter.Values[0] += "%"
	case FilterTypeEndsWith:
		operator = LIKE
		filter.Values[0] = "%" + filter.Values[0]
	case FilterTypeIn:
		operator = IN
		inValues := make([]string, len(filter.Values))

		for i, value := range filter.Values {
			inValues[i] = formatValue(
				value,
				filter.ValueType,
				valueFormatter...)
		}

		values = "(" + strings.Join(inValues, ", ") + ")"
	case FilterTypeNotIn:
		operator = NOTIN
		inValues := formatValueForNotIn(filter, valueFormatter...)

		values = "(" + strings.Join(inValues, ", ") + ")"

	case BETWEEN:
		operator = BETWEEN
		from = strconv.Itoa(filter.Range.From)
		to = strconv.Itoa(filter.Range.To)

	case NOTBETWEEN:
		operator = NOTBETWEEN
		from = strconv.Itoa(filter.Range.From)
		to = strconv.Itoa(filter.Range.To)

	default:
		return ""
	}

	if operator == BETWEEN || operator == NOTBETWEEN {
		filterClause := fmt.Sprintf("%s %s '%s' AND '%s'",
			filter.Field, operator, from, to)

		return filterClause
	}

	if operator != IN && operator != NOTIN {
		values = formatValue(
			filter.Values[0],
			filter.ValueType,
			valueFormatter...,
		) // assuming only one value for now
	}

	field := filter.Field

	if operator == LIKE {
		field = "lower(" + field + ")" // match with lowercase
		values = strings.ToLower(values)
	}

	filterClause := fmt.Sprintf("%s %s %s", field, operator, values)

	return filterClause
}

func formatValueForNotIn(filter *Filter,
	valueFormatter ...ValueFormatter,
) []string {
	inValues := make([]string, len(filter.Values))

	for i, value := range filter.Values {
		inValues[i] = formatValue(
			value,
			filter.ValueType,
			valueFormatter...)
	}

	return inValues
}

func formatValue(
	value string,
	valueType ValueType,
	valueFormatter ...ValueFormatter,
) string {
	if valueType == ValueTypeString || valueType == "" {
		return getEscapedStringValue(value)
	}

	if len(valueFormatter) > 0 {
		return valueFormatter[0](value, valueType)
	}

	return value
}

func getEscapedStringValue(
	value string,
) string {
	return fmt.Sprintf("'%s'", strings.ReplaceAll(value, "'", "''"))
}

func ConvertUnixTimestamp(unixTimestamp int64, format string) string {
	t := time.Unix(unixTimestamp, 0)

	formattedDate := t.Format(format)

	return formattedDate
}
