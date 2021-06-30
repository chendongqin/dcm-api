package hbase

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	CompareOperatorLess           = "<"
	CompareOperatorGreater        = ">"
	CompareOperatorEqual          = "="
	CompareOperatorLessOrEqual    = "<="
	CompareOperatorGreaterOrEqual = ">="
	CompareOperatorNotEqual       = "!="
)

type FilterInterface interface {
	ToString() string
}

type Filter struct {
	filters []FilterInterface
}

func NewFilters() *Filter {
	return &Filter{filters: make([]FilterInterface, 0)}
}

func (receiver *Filter) Add(filterInterface FilterInterface) *Filter {
	return receiver.And(filterInterface)
}

func (receiver *Filter) IsEmpty() bool {
	return len(receiver.filters) <= 0
}

func (receiver *Filter) And(filterInterface FilterInterface) *Filter {
	receiver.filters = append(receiver.filters, filterInterface)
	return receiver
}

func (receiver *Filter) ToByteArray() []byte {
	var arr []string
	for _, filter := range receiver.filters {
		filterString := filter.ToString()
		arr = append(arr, filterString)
	}
	finalFilterString := strings.Join(arr, " AND ")
	return []byte(finalFilterString)
}

func b2s(val bool) string {
	if val {
		return "true"
	}
	return "false"
}

func BinaryComparator(val interface{}) string {
	return "binary:" + toString(val)
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch result := v.(type) {
	case string:
		return result
	case []byte:
		return string(result)
	default:
		return fmt.Sprint(result)
	}
}

// SingleColumnValueFilter

type SingleColumnValueFilter struct {
	Family                string
	Qualifier             string
	CompareOperator       string
	Comparator            string
	FilterIfColumnMissing bool
	LatestVersion         bool
}

func (f *SingleColumnValueFilter) ToString() string {
	return "SingleColumnValueFilter('" + f.Family + "', '" + f.Qualifier + "', " + f.CompareOperator + ", '" + f.Comparator + "', " + b2s(f.FilterIfColumnMissing) + ", " + b2s(f.LatestVersion) + ")"
}

func NewSingleColumnValueFilter(family string, qualifier string, compareOperator string, comparator string, filterIfColumnMissing bool, latestVersion bool) *SingleColumnValueFilter {
	return &SingleColumnValueFilter{
		Family:                family,
		Qualifier:             qualifier,
		CompareOperator:       compareOperator,
		Comparator:            comparator,
		FilterIfColumnMissing: filterIfColumnMissing,
		LatestVersion:         latestVersion,
	}
}

type PrefixFilter struct {
	Prefix string
}

func (f *PrefixFilter) ToString() string {
	return "PrefixFilter('" + f.Prefix + "')"
}

// ColumnPaginationFilter
type ColumnPaginationFilter struct {
	Limit  int
	Offset int
}

func (f *ColumnPaginationFilter) ToString() string {
	return "ColumnPaginationFilter(" + strconv.Itoa(f.Limit) + ", " + strconv.Itoa(f.Offset) + ")"
}
