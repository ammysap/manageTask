package paginate

import (
	"context"
	"fmt"
	"strings"
)

type ValueFormatter func(value string, valueType ValueType) string

func SearchWithCount[T any](
	ctx context.Context,
	search func(ctx context.Context, query string, skip int, limit int, sort string) ([]T, error),
	count func(ctx context.Context, query string) (int, error),
	request *PaginatedRequest,
	valueFormatter ...ValueFormatter,
) (*PaginatedResponse[T], error) {
	page := request.Page

	skip := page.Skip
	if skip <= 0 {
		skip = page.Number * page.Rows // assuming page to be 0 indexed here
	}

	limit := page.Rows
	resolvedFilter := resolveFilter(request.Filter, valueFormatter...)

	resp, err := search(
		ctx,
		resolvedFilter,
		skip,
		limit,
		resolveSorts(request.Sorts),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get paginated data, page: %+v : %w",
			page,
			err,
		)
	}

	total, err := count(ctx, resolvedFilter)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get count for paginated data, page: %+v : %w",
			page,
			err,
		)
	}

	return &PaginatedResponse[T]{
		Hits:    resp,
		HasMore: skip+limit < total,
		Total:   total,
	}, nil
}

func resolveFilter(
	filter *Filter,
	valueFormatter ...ValueFormatter,
) string {
	if filter == nil {
		return ""
	}

	return GenerateFilterClause(filter, valueFormatter...)
}

func resolveSorts(sorts []Sort) string {
	resolvedSorts := make([]string, len(sorts))

	for i, s := range sorts {
		resolvedSorts[i] = resolveSort(s)
	}

	return strings.Join(resolvedSorts, ", ")
}

func resolveSort(s Sort) string {
	return fmt.Sprintf("%s %s", s.Field, s.Order)
}

func resolveFields(groups []string, agg Agg) string {
	for idx := range groups {
		groups[idx] = fmt.Sprintf("%q", groups[idx])
	}

	var resolveAs string

	if agg.ResolveAs != "" {
		resolveAs = agg.ResolveAs
	} else {
		// the default value to maintain backward compatibility
		resolveAs = "func"
	}

	return fmt.Sprintf("%s, %s(%s) as %s", strings.Join(groups, ","),
		agg.Function,
		agg.Field,
		resolveAs,
	)
}
